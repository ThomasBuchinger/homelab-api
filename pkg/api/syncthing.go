package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/kubernetes"
)

func SetupSyncthingApiEndpoints(r *gin.Engine) *gin.Engine {
	r.GET("/api/syncthinghelper/metricsproxy", handleSyncthingMetricsProxy)
	r.DELETE("/api/syncthinghelper/restart", handleSyncthingRestart)

	return r
}

func handleSyncthingMetricsProxy(c *gin.Context) {
	conf := common.GetServerConfig()
	ApiLogger.Debug("Calling Mstrics: ", conf.Homelab.Syncthing.InternalMetricsUrl)
	ProxyWithBasicAuth(
		conf.Homelab.Syncthing.InternalMetricsUrl,
		conf.Homelab.Syncthing.AuthUser,
		conf.Homelab.Syncthing.AuthPass,
		c,
	)
}

func handleSyncthingRestart(c *gin.Context) {
	config := common.GetServerConfig()
	if config.Homelab.Syncthing.MockRestart != "" {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"reason": "Success",
		})
		return
	}

	k8s, err := kubernetes.NewKubernetesClient()
	if err != nil {
		ApiLogger.Error("Failed to initialize Kubernetes Client")
		c.JSONP(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": "Failed to create Kubernetes Client",
			"reason":  err.Error(),
		})
		return
	}

	err = k8s.RestartDeployment("syncthing", "syncthing")
	if err != nil {
		ApiLogger.Errorf("Failed restart Deployment %s/%s", "syncthing", "syncthing")
		c.JSONP(http.StatusInternalServerError, gin.H{
			"status":  "FAILED",
			"message": "Unable to update Deployment",
			"reason":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"reason": "Success",
	})
}
