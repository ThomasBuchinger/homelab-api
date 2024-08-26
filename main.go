package main

import (
	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/reconciler"
)

func main() {
	serverConfig := common.GetServerConfig()
	router := api.SetupRouter()

	serverConfig.RootLogger.Info("Starting HomeLAB API Server on :8080...")
	defer serverConfig.RootLogger.Sync()
	go reconciler.ReconcileLoop()
	router.Run(common.GetEnvWithDefault("BIND_ADDR", ":8080"))
}
