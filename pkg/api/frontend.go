package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/reconciler"
)

var SyncthingDeviceIdMapping map[string]string = map[string]string{
	"JJIQV62-GSKZ5JB-7ORT5P2-K7OBLU6-2GCHWE5-3HNI4X7-2XVEBQD-P6D6JQT": "__skip__",
	"KA73NYG-7KJX5BO-TSWMNWC-J2WYBWV-B2WQK7X-G2B3TSX-M6EKEJS-BA2KLQQ": "BS13",
	"2C3RPBD-V4ZPEJW-5SDPWK6-H2FHE3Z-O4GF7AR-7QY7PDZ-3CXYXH6-OIV46Q6": "BUC Lenovo",
	"DWP4JSL-6SSJQW3-ZKG25PY-NWD7IZH-5OFOV2E-BVREQ63-BKMKOUO-44CKNAJ": "Nokia 3.4",
}
var SyncthingFolderStateMapping map[int]string = map[int]string{
	0: "Sync",
	1: "Sync",
	2: "Sync",
	3: "Progress",
	4: "Progress",
	5: "Progress",
	6: "Progress",
	7: "Progress",
	8: "Error",
}

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

func handleComponentSyncthing(c *gin.Context) {
	syncthing := reconciler.SyncthingMetric
	devices := []gin.H{}
	for devId, metric_value := range syncthing.Metrics["device_connections"].GroupedValues {
		name, ok := SyncthingDeviceIdMapping[devId]
		status := "disconnected"
		if metric_value > 0 {
			status = "OK"
		}
		if ok && name != "__skip__" {
			devices = append(devices, gin.H{"display_name": SyncthingDeviceIdMapping[devId], "id": devId, "status": status})
		}
	}
	folders := []gin.H{}
	for folder, metric_value := range syncthing.Metrics["folder_state"].GroupedValues {
		folders = append(folders, gin.H{"display_name": folder, "status": SyncthingFolderStateMapping[int(metric_value)], "raw_status": metric_value})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  syncthing.GetSatus(),
		"reason":  syncthing.GetReason(),
		"url":     "http://syncthing.buc.sh",
		"devices": devices,
		"folders": folders,
	})
}

func handleComponentCombinedKubernetes(c *gin.Context) {
	conf := common.GetServerConfig()
	evergreenData := reconciler.EvergreenMetric.Metrics
	prodData := reconciler.ProdMetric.Metrics

	c.JSON(http.StatusOK, gin.H{
		"status": reconciler.EvergreenMetric.GetSatus(),
		"reason": reconciler.EvergreenMetric.GetReason(),

		"evergreen_url": conf.Homelab.Evergreen.ConsoleUrl,
		"prod_url":      conf.Homelab.Prod.ConsoleUrl,
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
