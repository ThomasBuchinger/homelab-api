package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

func SetupSyncthingApiEndpoints(r *gin.Engine) *gin.Engine {
	r.GET("/api/syncthinghelper/metricsproxy", handleSyncthingMetricsProxy)
	r.GET("/api/syncthinghelper/restart", handleSyncthingRestart)

	return r
}

func handleSyncthingMetricsProxy(c *gin.Context) {
	conf := common.GetServerConfig()
	ProxyWithBasicAuth(
		conf.Homelab.Syncthing.InternalMetricsUrl,
		conf.Homelab.Syncthing.AuthUser,
		conf.Homelab.Syncthing.AuthPass,
		c,
	)
}

func handleSyncthingRestart(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":          "OK",
		"reason":          "Success",
		"total_documents": "-1",
		"url":             "https://paperless.buc.sh",
		"alt_url":         "https://paperless.10.0.0.21.nip.io",
	})
}
