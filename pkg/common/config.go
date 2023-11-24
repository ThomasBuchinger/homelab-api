package common

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ServerModeDev         = "dev"
	ServerModeInternal    = "internal"
	ServerModePublic      = "public"
	ServerModeCopyGeoipDB = "copy_geoip_database"

	EnvMode                = "MODE"
	EnvGeoipDatabase       = "GEOIP_DATABASE"
	EnvCopyGeipDestination = "COPY_GEOIP_DATABASE"
	EnvAuthCountries       = "AUTH_COUNTRIES"
	EnvAuthUser            = "AUTH_COUNTRIES"
	EnvAuthPass            = "AUTH_COUNTRIES"
)

type Serverconfig struct {
	GinMode       string
	JsonLogs      bool
	RootLogger    *zap.SugaredLogger
	AuthCountries string

	EnableInternalApis bool
	EnableLegacyApi    bool

	TaskCopyGeoip bool

	EnableGeoip       bool
	GeoipDatabasePath string
}

func GetServerConfig() Serverconfig {
	var logger *zap.Logger
	if EnableFeatureInMode([]string{ServerModeDev}, true, false) {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	return Serverconfig{
		GinMode:       EnableFeatureInMode([]string{ServerModeDev}, gin.DebugMode, gin.ReleaseMode),
		JsonLogs:      EnableFeatureInMode([]string{ServerModeInternal, ServerModePublic}, true, false),
		RootLogger:    logger.Sugar(),
		AuthCountries: GetEnvWithDefault(EnvAuthCountries, ""),

		EnableInternalApis: EnableFeatureInMode([]string{ServerModeDev, ServerModeInternal}, true, false),
		EnableLegacyApi:    EnableFeatureInMode([]string{ServerModeDev, ServerModeInternal}, true, false),

		TaskCopyGeoip: GetEnvWithDefault(EnvCopyGeipDestination, "") != "",

		EnableGeoip:       featureGeoip.Enabled,
		GeoipDatabasePath: featureGeoip.DatapasePath,
	}
}

func EnableFeatureInMode[V any](mode []string, enabled V, disabled V) V {
	mode_envvar := GetEnvWithDefault(EnvMode, "dev")
	for _, v := range mode {
		if mode_envvar == v {
			return enabled
		}
	}
	return disabled
}

func GetEnvWithDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}
