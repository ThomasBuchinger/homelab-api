package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/backend"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func embeddReactUI(router *gin.Engine) *gin.Engine {
	router.Use(static.Serve("/", static.LocalFile("./ui/out", true)))
	return router
}

func main() {
	fmt.Println("Stating Homelab API...")
	router := gin.Default()
	gin.SetMode(common.GetServerConfig().GinMode)

	router = embeddReactUI(router)
	router = backend.SetupApi(router)

	router.Run(":8080")
}
