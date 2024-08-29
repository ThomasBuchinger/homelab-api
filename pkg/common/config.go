package common

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	ConfigFileContent
}
type ConfigFileContent struct {
	Homelab struct {
		Evergreen struct {
			KubeStateMetricsUrl string
			ConsoleUrl          string
		}
		Prod struct {
			KubeStateMetricsUrl string
			ConsoleUrl          string
		}
		Syncthing struct {
			MetricsUrl string
			AuthUser   string
			AuthPass   string
		}
		Paperless struct {
			MetricsUrl string
			AuthUser   string
			AuthPass   string
		}
	}
}

func SetupViperConfig() {
	// viper.SetConfigName("*")
	// viper.AddConfigPath("./config")
	// viper.AddConfigPath("/config")

	viper.SetDefault("homelab.evergreen.consoleUrl", "")
	viper.SetDefault("homelab.evergreen.kubeStateMetricsUrl", "")
	viper.SetDefault("homelab.prod.consoleUrl", "")
	viper.SetDefault("homelab.prod.kubeStateMetricsUrl", "")
	viper.SetDefault("homelab.syncthing.metricsUrl", "d")
	viper.SetDefault("homelab.syncthing.authUser", "")
	viper.SetDefault("homelab.syncthing.authPass", "")
	viper.SetDefault("homelab.id", "default")

	viper.BindEnv("homelab.syncthing.authUser", "SYNCTHING_USER")
	viper.BindEnv("homelab.syncthing.authPass", "SYNCTHING_PASSWORD", "SYNCTHING_PASS")

	viper.SetConfigFile("./config/local.yaml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.SetConfigFile("./config/secret.yaml")
	viper.MergeInConfig()
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

	var labconf ConfigFileContent
	err := viper.Unmarshal(&labconf)
	if err != nil {
		fmt.Println("PANIC: got an error reading config")
		panic(err)
	}
	conf.ConfigFileContent = labconf
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
