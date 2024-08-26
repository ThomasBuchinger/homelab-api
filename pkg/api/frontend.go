package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/reconciler"
)

func setupFrontendApiEndpoints(r *gin.Engine) *gin.Engine {
	r.GET("/api/component/paperless", handleComponentPaperless)
	r.GET("/api/component/syncthing", handleComponentSyncthing)
	r.GET("/api/component/kubernetes", handleComponentCombinedKubernetes)
	return r
}

func handleComponentPaperless(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":          "OK",
		"reason":          "Success",
		"total_documents": "-1",
		"url":             "https://paperless.buc.sh",
		"alt_url":         "https://paperless.10.0.0.21.nip.io",
	})
}
func handleComponentNasv3(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":           "OK",
		"reason":           "Success",
		"url":              "http://10.0.0.19",
		"backup_status":    "OK",
		"backup_timestamp": "20240822T03:00:00Z",
		"backup_reason":    "Success",
		"disk_total":       "XX TB",
		"disk_free":        "YY TB",
		"disk_smart_ok":    "4",
		"disk_smart_fail":  "2",
	})
}

func handleComponentBastion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":            "OK",
		"reason":            "Success",
		"url":               "http://10.0.0.22:9090",
		"last_update":       "20240822T06:00:00Z",
		"service_wireguard": "active",
	})
}

type ApiSyncthingDeviceStatusV1 struct {
	Id          string
	DisplayName string
	Status      string
}
type ApiSyncthingFolderStatusV1 struct {
	Id          string
	DisplayName string
	Status      string
}

func handleComponentSyncthing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"reason":  "Success",
		"url":     "http://syncthing.buc.sh",
		"alt_url": "syncthing.10.0.0.21.nip.io",
		"devices": []ApiSyncthingDeviceStatusV1{
			ApiSyncthingDeviceStatusV1{Id: "laptop", DisplayName: "My Laptop", Status: "Connected"},
		},
		"folders": []ApiSyncthingFolderStatusV1{
			ApiSyncthingFolderStatusV1{Id: "dir1", DisplayName: "A Directory", Status: "Sync"},
		},
	})
}

func handleComponentCombinedKubernetes(c *gin.Context) {
	conf := common.GetServerConfig()
	evergreenData := reconciler.KubernetesMetricEvergreen.Metrics
	prodData := reconciler.KubernetesMetricProd.Metrics

	c.JSON(http.StatusOK, gin.H{
		"status": reconciler.KubernetesMetricEvergreen.GetSatus(),
		"reason": reconciler.KubernetesMetricEvergreen.GetReason(),

		"url_evergreen": conf.HomelabEnv.EvergreenConsoleUrl,
		"urk_prod":      conf.HomelabEnv.ProdConsoleUrl,
		"pod_healthy": evergreenData["num_pod_succeess"].Value +
			evergreenData["num_pod_running"].Value +
			prodData["num_pod_succeess"].Value +
			prodData["num_pod_running"].Value,
		"pod_total":   evergreenData["num_pod_total"].Value + prodData["num_pod_total"].Value,
		"pvc_healthy": evergreenData["num_pvc_bound"].Value + prodData["num_pvc_bound"].Value,
		"pvc_total":   evergreenData["num_pvc_total"].Value + prodData["num_pvc_total"].Value,
		"evergreen": gin.H{
			"pod_healthy": evergreenData["num_pod_succeess"].Value + evergreenData["num_pod_running"].Value,
			"pod_total":   evergreenData["num_pod_total"].Value,
			"pvc_healthy": evergreenData["num_pvc_bound"].Value,
			"pvc_total":   evergreenData["num_pvc_total"].Value,
		},
		"prod": gin.H{
			"pod_healthy": prodData["num_pod_succeess"].Value + prodData["num_pod_running"].Value,
			"pod_total":   prodData["num_pod_total"].Value,
			"pvc_healthy": prodData["num_pvc_bound"].Value,
			"pvc_total":   prodData["num_pvc_total"].Value,
		},
	})
}
