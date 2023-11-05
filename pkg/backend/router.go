package backend

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func SetupApi(r *gin.Engine) *gin.Engine {

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
	}
	return r
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
