package backend

import (
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

type AuthorizationConfig struct {
	AllowedCountries []string
	TempAllow        map[string]TempAllow
}
type TempAllow struct {
	Port        int
	ExpireDate  time.Time
	RequireAuth bool
	PathRegex   []string
}

var authConfig AuthorizationConfig = AuthorizationConfig{
	AllowedCountries: []string{"AT", "HR", "IT", "CZ"},
	TempAllow:        map[string]TempAllow{},
}

func AuthByGeoip(ip string) bool {
	cfg := common.GetServerConfig()

	if !cfg.EnableGeoip {
		return false // Deny all traffic if GeoIP is not available
	}
	req_country, err := common.LookupIP(ip)
	if err != nil {
		return false
	}

	for _, allowed_country := range authConfig.AllowedCountries {
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
	conf.TempAllow[host] = TempAllow{
		Port:        0,
		ExpireDate:  time.Now().Add(expires),
		RequireAuth: auth,
		PathRegex:   paths,
	}
}
