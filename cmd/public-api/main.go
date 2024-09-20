package main

import (
	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func main() {
	common.SetupViperConfig()
	serverConfig := common.GetServerConfig()

	router := api.SetupDefaultRouter()
	router = api.SetupAuthApiEndpoints(router)

	serverConfig.RootLogger.Info("Starting PUBLIC HomeLAB API Server on :8080...")
	defer serverConfig.RootLogger.Sync()
	// go reconciler.ReconcileLoop()
	router.Run(common.GetEnvWithDefault("BIND_ADDR", serverConfig.BindAddr))
}
