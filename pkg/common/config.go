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

	EnvMode = "MODE"
)

type Serverconfig struct {
	GinMode    string
	JsonLogs   bool
	RootLogger *zap.SugaredLogger

	HomelabEnv struct {
		EvergreenKubeStateMetricsUrl string
		EvergreenConsoleUrl          string
		ProdKubeStateMetricsUrl      string
		ProdConsoleUrl               string
	}
}

func GetServerConfig() Serverconfig {
	conf := Serverconfig{
		GinMode:  EnableFeatureInMode([]string{ServerModeDev}, gin.DebugMode, gin.ReleaseMode),
		JsonLogs: EnableFeatureInMode([]string{ServerModeInternal, ServerModePublic}, true, false),
	}

	var logger *zap.Logger
	if EnableFeatureInMode([]string{ServerModeDev}, true, false) {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	conf.RootLogger = logger.Sugar()

	conf.HomelabEnv.EvergreenConsoleUrl = "http://evergreen-console.10.0.0.16.nip.io"
	conf.HomelabEnv.EvergreenKubeStateMetricsUrl = "http://kube-state-metrics.10.0.0.16.nip.io/metrics"
	conf.HomelabEnv.ProdConsoleUrl = "http://prod-console.10.0.0.21.nip.io"
	conf.HomelabEnv.ProdKubeStateMetricsUrl = "http://kube-state-metrics.10.0.0.21.nip.io/metrics"

	return conf
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
