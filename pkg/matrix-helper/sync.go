package matrixhelper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/kubernetes"
	"go.uber.org/zap"
)

type MatrixHelperConfig struct {
	ClientsecretNamespace, ClientsecretName, ClientsecretKey string
	MasNamespace, MasSecretName, MasSecretKey, MasDeployment string
	PollInterval                                             time.Duration
	HomeserverURL, ServerName, BotUsername, BotPassword      string
}

var logger *zap.SugaredLogger
var k8s *kubernetes.KubernetesClient
var matrixClient *MatrixClient

func NewMatrixHelperConfig() MatrixHelperConfig {
	return MatrixHelperConfig{
		ClientsecretNamespace: common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_NAMESPACE", "keycloak"),
		ClientsecretName:      common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_NAME", "keycloak-client-secret-client-matrix"),
		ClientsecretKey:       common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_KEY", "CLIENT_SECRET"),
		MasNamespace:          "matrix",
		MasDeployment:         "matrix-matrix-authentication-service",
		MasSecretName:         "matrix-matrix-authentication-service",
		MasSecretKey:          "user-keycloak",
		PollInterval:          10 * time.Minute,
		HomeserverURL:         "https://matrix.cloud.buc.sh",
		ServerName:            "matrix.cloud.buc.sh",

		BotUsername: common.GetEnvWithDefault("MATRIX_HELPER_BOT_USERNAME", ""),
		BotPassword: common.GetEnvWithDefault("MATRIX_HELPER_BOT_PASSWORD", ""),
	}
}

func SyncLoop(ctx context.Context, l *zap.SugaredLogger) {
	var err error
	logger = l
	cfg := NewMatrixHelperConfig()
	matrixClient = NewMatrixClient(cfg.HomeserverURL, cfg.ServerName, cfg.BotUsername, cfg.BotPassword)
	k8s, err = kubernetes.NewKubernetesClient()
	if err != nil {
		panic(fmt.Errorf("creating Kubernetes client: %w", err))
	}
	for {
		syncMatrix(cfg)
		syncOidcClient(cfg)
		select {
		case <-ctx.Done():
			return
		case <-time.After(cfg.PollInterval):
		}
	}
}

func syncOidcClient(cfg MatrixHelperConfig) {
	sourceSecretLabel := fmt.Sprintf("%s/%s:%s", cfg.ClientsecretNamespace, cfg.ClientsecretName, cfg.ClientsecretKey)
	masConfigLabel := fmt.Sprintf("Deployment: %s/%s | Secret: %s:%s", cfg.MasNamespace, cfg.MasDeployment, cfg.MasSecretName, cfg.MasSecretKey)

	sourceSecret, err := k8s.GetSecret(cfg.ClientsecretNamespace, cfg.ClientsecretName)
	if err != nil {
		panic(fmt.Errorf("Unable to get Source Secret: "+sourceSecretLabel+": %w", err))
	}
	masSecret, err := k8s.GetSecret(cfg.MasNamespace, cfg.MasSecretName)
	if err != nil {
		panic(fmt.Errorf("Unable to get Mas Config Secret: "+masConfigLabel+": %w", err))
	}

	desiredMasConfig, changed, err := ensureOidcClientSecret(string(masSecret.Data[cfg.MasSecretKey]), string(sourceSecret.Data[cfg.ClientsecretKey]))
	if err != nil {
		panic(fmt.Errorf("ensuring OIDC client secret | Source: %s | Mas: %s: %w", sourceSecretLabel, masConfigLabel, err))
	}

	metricCheckTimestamp.SetToCurrentTime()
	if !changed {
		logger.Infof("OK: No action needed")
		return
	}

	// Update Secret...
	masSecret.Data[cfg.MasSecretKey] = []byte(desiredMasConfig)
	if err := k8s.UpdateSecret(cfg.MasNamespace, masSecret); err != nil {
		panic(err)
	}
	metricChangeTotal.Inc()
	metricChangeTimestamp.SetToCurrentTime()

	// ... and restart MAS
	if err := k8s.RestartDeployment(cfg.MasNamespace, cfg.MasDeployment); err != nil {
		panic(err)
	}
	metricMasRestartTotal.Inc()
	metricMasRestartTimestamp.SetToCurrentTime()
	logger.Infof("OK: applied new client_secret to %s/%s and requested restart of %s/%s", cfg.MasNamespace, cfg.MasSecretName, cfg.MasNamespace, cfg.MasDeployment)
}

func syncMatrix(cfg MatrixHelperConfig) {

	ensureUser("bot", cfg.BotPassword)
	ensureUser("glados", cfg.BotPassword)

	if err := matrixClient.Login(); err != nil {
		panic(fmt.Errorf("logging in as %s: %w", cfg.BotUsername, err))
	}
	metricMatrixConnected.Set(1)

	ensureRoom("home", "public", "public_chat")
	ensureRoom("alerts", "public", "public_chat")
	ensureRoom("bot", "public", "public_chat")
	ensureRoom("timerec", "public", "public_chat")
	ensureRoom("glados", "public", "public_chat")

	// join the default bot to rooms
	joinRoom("home")
	joinRoom("bot")

	setMessageRetention("home", retentionPeriod)
	setMessageRetention("bot", retentionPeriod)

	metricMatrixCheckTimestamp.SetToCurrentTime()
	logger.Infof("OK: Matrix stack in sync (users, rooms, retention)")
}

func ensureUser(user, password string) {
	podNamespace, podName, err := k8s.GetPodByLabel(masNamespace, masPodLabel)
	if err != nil {
		panic(fmt.Errorf("ensuring user %s: %w", user, err))
	}

	if _, _, err := k8s.ExecInPod(podNamespace, podName, masContainer, MasCliUnlockCommand(user)); err == nil {
		logger.Infof("OK: User %s exists", user)
		return
	}

	logger.Infof("Creating user %s", user)
	_, stderr, err := k8s.ExecInPod(podNamespace, podName, masContainer, MasCliRegisterCommand(user, password))
	if err != nil {
		panic(fmt.Errorf("ensuring user %s: %w (stderr: %s)", user, err, stderr))
	}
	metricMatrixUserCreatedTotal.Inc()
	logger.Infof("OK: Created user %s", user)
}

func ensureRoom(room, visibility, preset string) {
	alias := RoomAlias(room, matrixClient.ServerName())
	_, found, err := matrixClient.LookupRoom(alias)
	if err != nil {
		panic(fmt.Errorf("looking up room %s: %w", alias, err))
	}
	if found {
		return
	}

	logger.Infof("Creating room %s", alias)
	if _, err := matrixClient.CreateRoom(alias, visibility, preset); err != nil {
		panic(fmt.Errorf("creating room %s: %w", alias, err))
	}
	metricMatrixRoomsCreatedTotal.Inc()
	logger.Infof("OK: Created room %s", alias)
}

func joinRoom(room string) {
	alias := RoomAlias(room, matrixClient.ServerName())
	if _, err := matrixClient.JoinRoom(alias); err != nil {
		panic(fmt.Errorf("joining room %s: %w", alias, err))
	}
}

func setMessageRetention(room string, maxLifetime time.Duration) {
	alias := RoomAlias(room, matrixClient.ServerName())
	roomID, found, err := matrixClient.LookupRoom(alias)
	if err != nil {
		panic(fmt.Errorf("looking up room %s for retention: %w", alias, err))
	}
	if !found {
		panic(fmt.Errorf("room %s not found in directory before setting retention", alias))
	}

	var current *time.Duration
	currentValue, err := matrixClient.GetRetention(roomID)
	switch {
	case err == nil:
		current = &currentValue
	case errors.Is(err, errRetentionNotSet):
		// No policy yet: the write below applies it.
	default:
		panic(fmt.Errorf("reading retention for %s: %w", alias, err))
	}

	if !NeedsRetentionUpdate(current, maxLifetime) {
		return
	}

	if err := matrixClient.PutRetention(roomID, maxLifetime); err != nil {
		panic(fmt.Errorf("writing retention for %s: %w", alias, err))
	}
	logger.Infof("OK: Set retention on %s to %s", alias, maxLifetime)
}
