package backend

import (
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func SetupRouter() *gin.Engine {
	serverConfig := common.GetServerConfig()
	logger := serverConfig.RootLogger.Desugar().Named("access")

	gin.SetMode(serverConfig.GinMode)
	router := gin.New()
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{TimeFormat: time.RFC3339, UTC: true,
		SkipPaths: []string{
			"/api/livez",
			"/api/readyz",
		},
	}))
	router.Use(ginzap.RecoveryWithZap(logger, false))
	router.Use(requestid.New())

	router.Use(static.Serve("/", static.LocalFile("./ui/out", true)))
	router.Use(static.Serve("/geoip", static.LocalFile("/geoip", true)))

	router = setupApiEndpoints(router)
	return router
}

func setupApiEndpoints(r *gin.Engine) *gin.Engine {

	r.GET("/api/livez", handlePing)
	r.GET("/api/readyz", handlePing)

	// Publicly accessible API endpoints
	r.GET("/api/public/server-config", handleServerConfig)
	r.GET("/api/public/client-config", handleClientConfig)
	r.GET("/api/public/ping", handlePing)
	r.GET("/api/public/health", handlePublicHealth)

	// API endpoints only available in "internal" mode
	serverConfig := common.GetServerConfig()
	if serverConfig.EnableInternalApis {
		r.GET("/api/legacy/ping", handlePing)
		r.GET("/api/legacy/proxy", handleLegacyProxy)
		r.GET("/api/legacy/metrics", handleLegacyMetrics)
	}
	if serverConfig.EnableInternalApis {
		r.GET("/api/internal/ping", handlePing)
		r.GET("/api/auth/simple/*authpath", handleAuthSimple)
		r.GET("/api/auth/login/*authpath", handleAuthWithCred)
	}
	return r
}

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Authorized group (uses gin.BasicAuth() middleware)
// Same than:
// authorized := r.Group("/")
// authorized.Use(gin.BasicAuth(gin.Credentials{
//	  "foo":  "bar",
//	  "manu": "123",
//}))
// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
// 	"foo":  "bar", // user:foo password:bar
// 	"manu": "123", // user:manu password:123
// }))

// authorized.POST("admin", func(c *gin.Context) {
// 	user := c.MustGet(gin.AuthUserKey).(string)
// 	// Parse JSON
// 	var json struct {
// 		Value string `json:"value" binding:"required"`
// 	}

// 	if c.Bind(&json) == nil {
// 		db[user] = json.Value
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	}
// })
