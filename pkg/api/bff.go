package api

import (
	"net/http"
	"strings"

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

func SetupBffApiEndpoints(r *gin.Engine) *gin.Engine {
	r.GET("/api/public/bff/paperless", handleComponentPaperless)
	r.GET("/api/public/bff/syncthing", handleComponentSyncthing)
	r.GET("/api/public/bff/kubernetes", handleComponentCombinedKubernetes)
	r.GET("/api/public/bff/nasv3", handleComponentNasv3)

	r.DELETE("/api/private/bff/syncthing/restart", handleCommandRestartSyncthing)
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
	m := reconciler.NasNodeMetrics
	disks := []gin.H{}
	for disk, metric_value := range m.Metrics["btrfs_errors"].GroupedValues {
		cap_total, total_ok := m.Metrics["btrfs_total"].GroupedValues[disk]
		cap_free, free_ok := m.Metrics["btrfs_unused"].GroupedValues[disk]
		if strings.Contains(disk, "loop") && total_ok && free_ok {
			continue
		}

		disks = append(disks, gin.H{"display_name": disk, "btrfs_error": metric_value == 0, "capacity_free": cap_free, "capacity_total": cap_total})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           m.GetSatus(),
		"reason":           m.GetReason(),
		"url":              "http://10.0.0.19",
		"backup_url":       "http://10.0.0.19:9898",
		"parity_status":    "TODO",
		"parity_timestamp": "20240822T03:00:00Z",
		"parity_reason":    "Success",
		"backup_status":    "TODO",
		"backup_timestamp": "20240822T03:00:00Z",
		"backup_reason":    "Success",
		"disk_total":       m.Metrics["fs_total"].Value,
		"disk_free":        m.Metrics["fs_free"].Value,
		"disks":            disks,
	})
}

func handleComponentBastion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":                "OK",
		"reason":                "Success",
		"url":                   "http://10.0.0.22:9090",
		"service_dnf_automatic": "active",
		"last_update":           "20240822T06:00:00Z",
		"service_wireguard":     "active",
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

func handleCommandRestartSyncthing(c *gin.Context) {
	ApiLogger.Debugln("Restarting Syncthing...")
	config := common.GetServerConfig()
	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodDelete, config.Homelab.Syncthing.RestartUrl, http.NoBody)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ApiLogger.Error("Failed to call Syncthing-Helper Restart Hook: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": common.ReconcilerStatusError,
			"reason": err.Error(),
		})
	}

	c.DataFromReader(res.StatusCode, res.ContentLength, "application/json", res.Body, map[string]string{})
}
