package main

import (
	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/reconciler"
)

func main() {
	common.SetupViperConfig()
	serverConfig := common.GetServerConfig()
	router := api.SetupDefaultRouter()
	router = api.SetupStaticFileServing(router)
	router = api.SetupFrontendApiEndpoints(router)
	if common.EnableFeatureInMode([]string{common.ServerModeDev}, true, false) {
		router = api.SetupSyncthingApiEndpoints(router)
	}

	serverConfig.RootLogger.Info("Starting HomeLAB API Server on :8080...")
	defer serverConfig.RootLogger.Sync()
	go reconciler.ReconcileLoop()
	router.Run(common.GetEnvWithDefault("BIND_ADDR", serverConfig.BindAddr))
}
