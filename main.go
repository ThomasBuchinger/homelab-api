package main

import (
	"fmt"
	"os"

	"github.com/thomasbuchinger/homelab-api/pkg/backend"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func copyGeoipDatabase(sourceFile, destinationFile string) {
	input, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}
}

func main() {
	common.FeatureGeoipInit()
	defer common.FeatureGeoipClose()

	serverConfig := common.GetServerConfig()
	if serverConfig.TaskCopyGeoip {
		copyGeoipDatabase(common.GetServerConfig().GeoipDatabasePath, common.GetEnvWithDefault(common.EnvCopyGeipDestination, ""))
		fmt.Println("Successfully Copied Database")
		os.Exit(0)
	}
	router := backend.SetupRouter()

	serverConfig.RootLogger.Info("Starting HomeLAB API Server on :8080...")
	defer serverConfig.RootLogger.Sync()
	router.Run(common.GetEnvWithDefault("BIND_ADDR", ":8080"))
}
