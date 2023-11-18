package backend

import (
	"strings"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

type AuthorizationConfig struct {
	AllowedCountries []string
	Users            map[string]AuthUser
	Rules            map[string]AuthAllowRule
}
type AuthAllowRule struct {
	Port        int
	ExpireDate  time.Time
	RequireAuth bool
	PathRegex   []string
}
type AuthUser struct {
	Password string
}

var authConfig AuthorizationConfig = CreateAuthConfig(
	common.GetEnvWithDefault("AUTH_COUNTRIES", ""),
	common.GetEnvWithDefault("AUTH_USER", ""),
	common.GetEnvWithDefault("AUTH_PASS", ""),
)

func CreateAuthConfig(allowed_countries, auth_user, auth_pass string) AuthorizationConfig {
	countries := []string{}
	for _, c := range strings.Split(allowed_countries, ",") {
		countries = append(countries, strings.ToUpper(strings.Trim(c, " ")))
	}
	users := map[string]AuthUser{}
	if auth_user != "" && auth_pass != "" {
		users = map[string]AuthUser{auth_user: {Password: auth_pass}}
	}

	return AuthorizationConfig{
		AllowedCountries: countries,
		Users:            users,
		Rules:            map[string]AuthAllowRule{},
	}
}

func AuthByGeoip(ip string, conf AuthorizationConfig) bool {
	cfg := common.GetServerConfig()

	if !cfg.EnableGeoip {
		return false // Deny all traffic if GeoIP is not available
	}
	req_country, err := common.LookupIP(ip)
	if err != nil {
		return false
	}

	for _, allowed_country := range conf.AllowedCountries {
		if req_country == allowed_country {
			return true
		}
	}
	return false
}

func AllowInternalIp(ip string) bool {
	return common.IsIpAddressInternal(ip)
}

func AuthByCredentials(user string, password string) bool {
	auth_user := common.GetEnvWithDefault("AUTH_USER", "")
	auth_pass := common.GetEnvWithDefault("AUTH_PASS", "")

	return (auth_user != "" && auth_user == user) && (auth_pass != "" && auth_pass == password)
}

func AuthByTempAllow(host string, path string) bool {
	return host == path && false
}

func (conf *AuthorizationConfig) AddTempAllow(host string, paths []string, expires time.Duration, auth bool) {
	conf.Rules[host] = AuthAllowRule{
		Port:        0,
		ExpireDate:  time.Now().Add(expires),
		RequireAuth: auth,
		PathRegex:   paths,
	}
}
