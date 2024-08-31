package main

import (
	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func main() {
	common.SetupViperConfig()
	serverConfig := common.GetServerConfig()
	router := api.SetupDefaultRouter()
	router = api.SetupSyncthingApiEndpoints(router)

	router.Run(common.GetEnvWithDefault("BIND_ADDR", serverConfig.BindAddr))
}
