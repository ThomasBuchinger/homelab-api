package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/backend"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func embeddReactUI(router *gin.Engine) *gin.Engine {
	router.Use(static.Serve("/", static.LocalFile("./ui/out", true)))
	router.Use(static.Serve("/geoip", static.LocalFile("/geoip", true)))
	return router
}

func CopyGeoipDatabase(sourceFile, destinationFile string) {
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

	if common.GetEnvWithDefault("COPY_GEOIP_DATABASE", "") != "" {
		CopyGeoipDatabase(common.GetServerConfig().GeoipDatabasePath, common.GetEnvWithDefault("COPY_GEOIP_DATABASE", ""))
		fmt.Println("Successfully Copied Database")
		os.Exit(0)
	}

	fmt.Println("Stating Homelab API...")
	router := gin.Default()
	gin.SetMode(common.GetServerConfig().GinMode)

	router = embeddReactUI(router)
	router = backend.SetupApi(router)

	router.Run(common.GetEnvWithDefault("BIND_ADDR", ":8080"))
}
