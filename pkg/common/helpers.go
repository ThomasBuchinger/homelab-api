package common

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ServerModeDev      = "dev"
	ServerModeInternal = "internal"
	ServerModePublic   = "public"
)

type Serverconfig struct {
	GinMode            string
	EnableInternalApis bool
	EnableLegacyApi    bool
	EnableGeoip        bool

	GeoipDatabasePath string
}

func GetServerConfig() Serverconfig {
	return Serverconfig{
		GinMode:            EnableFeatureInMode([]string{"dev"}, gin.DebugMode, gin.ReleaseMode),
		EnableInternalApis: EnableFeatureInMode([]string{"dev", "internal"}, true, false),
		EnableLegacyApi:    EnableFeatureInMode([]string{"dev", "internal"}, true, false),
		EnableGeoip:        featureGeoip.Enabled,

		GeoipDatabasePath: featureGeoip.DatapasePath,
	}
}

func IsIpAddressInternal(ip string) bool {
	if ip == "127.0.0.1" || ip == "::ffff:127.0.0.1" {
		return true
	}
	if strings.HasPrefix(ip, "10.0.0") && ip != "10.0.0.21" {
		return true
	}
	return false
}

func GetEnvWithDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}

func EnableFeatureInMode[V any](mode []string, enabled V, disabled V) V {
	mode_envvar := GetEnvWithDefault("MODE", "dev")
	for _, v := range mode {
		if mode_envvar == v {
			return enabled
		}
	}
	return disabled
}
