package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	ServerModeDev         = "dev"
	ServerModeInternal    = "internal"
	ServerModePublic      = "public"
	ServerModeCopyGeoipDB = "copy_geoip_database"
)

type Serverconfig struct {
	GinMode    string
	JsonLogs   bool
	RootLogger *zap.SugaredLogger

	ConfigFileContent
}
type ConfigFileContent struct {
	BindAddr string
	Mode     string

	Homelab struct {
		Evergreen struct {
			KubeStateMetricsUrl string
			ConsoleUrl          string
		}
		Prod struct {
			KubeStateMetricsUrl string
			ConsoleUrl          string
		}
		Nas struct {
			MetricsUrl string
		}
		Syncthing struct {
			MetricsUrl         string
			InternalMetricsUrl string
			RestartUrl         string
			MockRestart        string
			AuthUser           string
			AuthPass           string
		}
		Paperless struct {
			MetricsUrl string
			AuthUser   string
			AuthPass   string
		}
	}
	EnvoyAuthConfig struct {
		AllowedCountries string
	}
}

func SetupViperConfig() {

	viper.SetDefault("bindAddr", ":8080")
	viper.SetDefault("mode", "dev")
	viper.SetDefault("homelab.evergreen.consoleUrl", "http://evergreen-console.10.0.0.16.nip.io")
	viper.SetDefault("homelab.evergreen.kubeStateMetricsUrl", "http://kube-state-metrics.10.0.0.16.nip.io/metrics")
	viper.SetDefault("homelab.prod.consoleUrl", "http://prod-console.10.0.0.21.nip.io")
	viper.SetDefault("homelab.prod.kubeStateMetricsUrl", "http://kube-state-metrics.10.0.0.21.nip.io/metrics")
	viper.SetDefault("homelab.nas.metricsUrl", "http://10.0.0.19:9100/metrics")
	viper.SetDefault("homelab.syncthing.metricsUrl", "http://syncthing.10.0.0.21.nip.io/metrics")
	viper.SetDefault("homelab.syncthing.internalMetricsUrl", "http://syncthing:8384/metrics")
	viper.SetDefault("homelab.syncthing.restartUrl", "http://syncthing.10.0.0.21.nip.io/api/syncthinghelper/restart")
	viper.SetDefault("homelab.syncthing.mockRestart", "")
	viper.SetDefault("homelab.syncthing.authUser", "")
	viper.SetDefault("homelab.syncthing.authPass", "")
	viper.SetDefault("envoyAuthConfig.allowedCountries", "AT,CH,DE,IT,HR")

	viper.BindEnv("bindAddr", "BIND_ADDR")
	viper.BindEnv("mode", "MODE")
	viper.BindEnv("envoyAuthConfig.allowedCountries", "AUTH_ALLOWED_COUNTRIES")
	viper.BindEnv("homelab.syncthing.internalMetricsUrl", "SYNCTHING_INTERNAL_METRICSURL")

	viper.BindEnv("homelab.syncthing.metricsUrl", "SYNCTHING_METRICSURL")
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

func GetConfigEnvoyAuthAllowedCountries() []string {
	var allowed_countries []string
	for _, country := range strings.Split(GetServerConfig().EnvoyAuthConfig.AllowedCountries, ",") {
		allowed_countries = append(allowed_countries, strings.Trim(strings.ToUpper(country), " "))
	}
	return allowed_countries
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

func GetEnvWithDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}
