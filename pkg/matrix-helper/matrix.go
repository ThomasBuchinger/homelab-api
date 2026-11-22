package matrixhelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	masNamespace    = "matrix"
	masPodLabel     = "app.kubernetes.io/name=matrix-authentication-service"
	masContainer    = "matrix-authentication-service"
	botDeviceID     = "matrixhelper"
	retentionPeriod = 7 * 24 * time.Hour
)

var errRetentionNotSet = errors.New("retention state not set")

type MatrixClient struct {
	httpClient  *http.Client
	apiBase     string
	serverName  string
	user        string
	password    string
	accessToken string
}

func NewMatrixClient(homeserverURL, serverName, user, password string) *MatrixClient {
	return &MatrixClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiBase:    strings.TrimRight(homeserverURL, "/") + "/_matrix/client/v3",
		serverName: serverName,
		user:       user,
		password:   password,
	}
}

func (c *MatrixClient) ServerName() string {
	return c.serverName
}

func (c *MatrixClient) Login() error {
	if c.accessToken != "" {
		return nil
	}
	body := map[string]string{
		"type":      "m.login.password",
		"user":      c.user,
		"password":  c.password,
		"device_id": botDeviceID,
	}
	var resp struct {
		AccessToken string `json:"access_token"`
	}
	if err := c.do(http.MethodPost, "/login", false, body, &resp); err != nil {
		return err
	}
	if resp.AccessToken == "" {
		return errors.New("login response did not contain an access token")
	}
	c.accessToken = resp.AccessToken
	return nil
}
func (c *MatrixClient) LookupRoom(alias string) (roomID string, found bool, err error) {
	var resp struct {
		RoomID string `json:"room_id"`
	}
	err = c.do(http.MethodGet, "/directory/room/"+url.PathEscape(alias), true, nil, &resp)
	if err != nil {
		var apiErr *apiError
		if errors.As(err, &apiErr) && apiErr.status == http.StatusNotFound {
			return "", false, nil
		}
		return "", false, err
	}
	return resp.RoomID, true, nil
}

func (c *MatrixClient) CreateRoom(alias, visibility, preset string) (roomID string, err error) {
	body := map[string]string{
		"room_alias_name": alias,
		"visibility":      visibility,
		"preset":          preset,
	}
	var resp struct {
		RoomID string `json:"room_id"`
	}
	if err = c.do(http.MethodPost, "/createRoom", true, body, &resp); err != nil {
		return "", err
	}
	return resp.RoomID, nil
}

func (c *MatrixClient) JoinRoom(alias string) (roomID string, err error) {
	var resp struct {
		RoomID string `json:"room_id"`
	}
	if err = c.do(http.MethodPost, "/join/"+url.PathEscape(alias), true, nil, &resp); err != nil {
		return "", err
	}
	return resp.RoomID, nil
}

func (c *MatrixClient) GetRetention(roomID string) (maxLifetime time.Duration, err error) {
	var resp struct {
		MaxLifetime int64 `json:"max_lifetime"`
	}
	err = c.do(http.MethodGet, "/rooms/"+url.PathEscape(roomID)+"/state/m.room.retention", true, nil, &resp)
	if err != nil {
		var apiErr *apiError
		if errors.As(err, &apiErr) && apiErr.status == http.StatusNotFound {
			return 0, errRetentionNotSet
		}
		return 0, err
	}
	return time.Duration(resp.MaxLifetime) * time.Millisecond, nil
}

func (c *MatrixClient) PutRetention(roomID string, maxLifetime time.Duration) error {
	body := map[string]int64{
		"max_lifetime": int64(maxLifetime / time.Millisecond),
	}
	return c.do(http.MethodPut, "/rooms/"+url.PathEscape(roomID)+"/state/m.room.retention", true, body, nil)
}

type apiError struct {
	status int
	body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("unexpected HTTP status %d: %s", e.status, e.body)
}

func (c *MatrixClient) do(method, path string, auth bool, reqBody, respBody any) error {
	var bodyReader io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("%s %s: marshaling request: %w", method, path, err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.apiBase+path, bodyReader)
	if err != nil {
		return fmt.Errorf("%s %s: creating request: %w", method, path, err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth {
		if c.accessToken == "" {
			return fmt.Errorf("%s %s: not logged in", method, path)
		}
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s %s: reading response: %w", method, path, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &apiError{status: resp.StatusCode, body: string(data)}
	}
	if respBody != nil {
		if err := json.Unmarshal(data, respBody); err != nil {
			return fmt.Errorf("%s %s: unmarshaling response: %w", method, path, err)
		}
	}
	return nil
}

func RoomAlias(room, serverName string) string {
	return "#" + room + ":" + serverName
}

func NeedsRetentionUpdate(current *time.Duration, desired time.Duration) bool {
	if current == nil {
		return true
	}
	return *current != desired
}

func MasCliUnlockCommand(user string) []string {
	return []string{"mas-cli", "manage", "unlock-user", user}
}

func MasCliRegisterCommand(user, password string) []string {
	return []string{
		"mas-cli", "manage", "register-user", user,
		"--password", password,
		"--display-name", user,
		"--email", user + "@buc.sh",
		"--no-admin",
		"--yes",
		"--ignore-password-complexity",
	}
}
