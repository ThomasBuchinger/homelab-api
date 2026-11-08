package reconciler

import (
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/metricscraper"
	"go.uber.org/zap"
)

var ProdMetric *metricscraper.MetricsReconciler
var EvergreenMetric *metricscraper.MetricsReconciler
var SyncthingMetric *metricscraper.MetricsReconciler
var NasNodeMetrics *metricscraper.MetricsReconciler
var NasSmartMetrics *metricscraper.MetricsReconciler
var PaperlessMetrics *metricscraper.MetricsReconciler
var FrigateMetrics *metricscraper.MetricsReconciler
var OpenwrtMetrics *metricscraper.MetricsReconciler
var HomeassistantMetrics *metricscraper.MetricsReconciler

func ReconcileLoop() {
	sleeptime, _ := time.ParseDuration("5m")
	conf := common.GetServerConfig()

	// kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Failed"} 0
	ProdMetric = metricscraper.NewMetricsReconciler(conf.Homelab.Prod.KubeStateMetricsUrl)
	ProdMetric.AddMetric("num_pod_total", metricscraper.Metric{Name: "kube_pod_status_phase"})
	ProdMetric.AddMetric("num_pod_succeess", metricscraper.Metric{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Succeeded"}})
	ProdMetric.AddMetric("num_pod_running", metricscraper.Metric{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Running"}})
	ProdMetric.AddMetric("num_pvc_total", metricscraper.Metric{Name: "kube_persistentvolume_status_phase"})
	ProdMetric.AddMetric("num_pvc_bound", metricscraper.Metric{Name: "kube_persistentvolume_status_phase", Labels: map[string]string{"phase": "Bound"}})

	EvergreenMetric = metricscraper.NewMetricsReconciler(conf.Homelab.Evergreen.KubeStateMetricsUrl)
	EvergreenMetric.AddMetric("num_pod_total", metricscraper.Metric{Name: "kube_pod_status_phase"})
	EvergreenMetric.AddMetric("num_pod_succeess", metricscraper.Metric{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Succeeded"}})
	EvergreenMetric.AddMetric("num_pod_running", metricscraper.Metric{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Running"}})
	EvergreenMetric.AddMetric("num_pvc_total", metricscraper.Metric{Name: "kube_persistentvolume_status_phase"})
	EvergreenMetric.AddMetric("num_pvc_bound", metricscraper.Metric{Name: "kube_persistentvolume_status_phase", Labels: map[string]string{"phase": "Bound"}})

	SyncthingMetric = metricscraper.NewMetricsReconciler(conf.Homelab.Syncthing.MetricsUrl)
	SyncthingMetric.AuthUser = conf.Homelab.Syncthing.AuthUser
	SyncthingMetric.AuthPass = conf.Homelab.Syncthing.AuthPass
	SyncthingMetric.AddMetric("folder_state", metricscraper.Metric{Name: "syncthing_model_folder_state", GroupBy: "folder"})
	SyncthingMetric.AddMetric("device_connections", metricscraper.Metric{Name: "syncthing_connections_active", GroupBy: "device"})

	NasNodeMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Nas.NodeExporter)
	NasNodeMetrics.AddMetric("btrfs_errors", metricscraper.Metric{Name: "node_btrfs_device_errors_total", GroupBy: "device"})
	NasNodeMetrics.AddMetric("btrfs_total", metricscraper.Metric{Name: "node_btrfs_device_size_bytes", GroupBy: "device"})
	NasNodeMetrics.AddMetric("btrfs_unused", metricscraper.Metric{Name: "node_btrfs_device_unused_bytes", GroupBy: "device"})
	NasNodeMetrics.AddMetric("fs_free", metricscraper.Metric{Name: "node_filesystem_avail_bytes", Labels: map[string]string{"mountpoint": "/mnt/user"}})
	NasNodeMetrics.AddMetric("fs_total", metricscraper.Metric{Name: "node_filesystem_size_bytes", Labels: map[string]string{"mountpoint": "/mnt/user"}})

	NasSmartMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Nas.SmartctlExporter)

	PaperlessMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Paperless.MetricsUrl)
	PaperlessMetrics.AddMetric("tag_name", metricscraper.Metric{Name: "paperless_tag_info", Labels: map[string]string{"id": conf.Homelab.Paperless.NewDocTag}, GroupBy: "name"})
	PaperlessMetrics.AddMetric("docs_new", metricscraper.Metric{Name: "paperless_tag_document_count", Labels: map[string]string{"id": conf.Homelab.Paperless.NewDocTag}})
	PaperlessMetrics.AddMetric("docs_total", metricscraper.Metric{Name: "paperless_document_type_document_count"})

	FrigateMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Frigate.MetricsUrl)
	FrigateMetrics.AddMetric("camera_fps", metricscraper.Metric{Name: "frigate_camera_fps", GroupBy: "camera_name"})
	FrigateMetrics.AddMetric("service_info", metricscraper.Metric{Name: "frigate_service_info"})
	FrigateMetrics.AddMetric("uptime", metricscraper.Metric{Name: "frigate_service_uptime_seconds"})
	FrigateMetrics.AddMetric("storage_free", metricscraper.Metric{Name: "frigate_storage_free_bytes", Labels: map[string]string{"storage": "/media/frigate/recordings"}})
	FrigateMetrics.AddMetric("storage_total", metricscraper.Metric{Name: "frigate_storage_total_bytes", Labels: map[string]string{"storage": "/media/frigate/recordings"}})
	FrigateMetrics.AddMetric("storage_used", metricscraper.Metric{Name: "frigate_storage_used_bytes", Labels: map[string]string{"storage": "/media/frigate/recordings"}})
	FrigateMetrics.AddMetric("events_person", metricscraper.Metric{Name: "frigate_camera_events_total", Labels: map[string]string{"label": "person"}, GroupBy: "camera"})

	OpenwrtMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Iotrouter.MetricsUrl)
	OpenwrtMetrics.AddMetric("uptime", metricscraper.Metric{Name: "node_boot_time_seconds"})
	OpenwrtMetrics.AddMetric("dhcp_leases_active", metricscraper.Metric{Name: "dnsmasq_leases_allocated_4_total"})
	OpenwrtMetrics.AddMetric("accesspoint_up", metricscraper.Metric{Name: "node_network_carrier", Labels: map[string]string{"device": "phy0-ap0"}})
	OpenwrtMetrics.AddMetric("network_info_iot", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "br-iot"}})
	OpenwrtMetrics.AddMetric("network_info_lan", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "br-lan"}})
	OpenwrtMetrics.AddMetric("network_info_docker", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "docker0"}})
	OpenwrtMetrics.AddMetric("network_info_eth0", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "eth0"}})
	OpenwrtMetrics.AddMetric("network_info_eth1", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "eth1"}})
	OpenwrtMetrics.AddMetric("network_info_ap", metricscraper.Metric{Name: "node_network_info", Labels: map[string]string{"device": "phy0-ap0"}})
	OpenwrtMetrics.AddMetric("wifi_info", metricscraper.Metric{Name: "wifi_radio_info"})
	OpenwrtMetrics.AddMetric("wifi_network", metricscraper.Metric{Name: "wifi_network_info"})
	OpenwrtMetrics.AddMetric("wifi_stations_active", metricscraper.Metric{Name: "wifi_stations_total"})
	OpenwrtMetrics.AddMetric("wifi_stations_inactive", metricscraper.Metric{Name: "wifi_station_inactive_milliseconds", GroupBy: "mac"})
	OpenwrtMetrics.AddMetric("wifi_stations_signal", metricscraper.Metric{Name: "wifi_station_signal_dbm", GroupBy: "mac"})

	HomeassistantMetrics = metricscraper.NewMetricsReconciler(conf.Homelab.Homeassistant.MetricsUrl)
	// HomeassistantMetrics.AddMetric("uptime", metricscraper.Metric{Name: "node_boot_time_seconds"})
	HomeassistantMetrics.AddMetric("entity_available", metricscraper.Metric{Name: "homeassistant_entity_available", GroupBy: "entity"})

	for {
		conf.RootLogger.Logln(zap.DebugLevel, "Running Reconcilers...")

		EvergreenMetric.Reconcile()
		ProdMetric.Reconcile()
		SyncthingMetric.Reconcile()
		NasNodeMetrics.Reconcile()
		PaperlessMetrics.Reconcile()

		OpenwrtMetrics.Reconcile()
		FrigateMetrics.Reconcile()
		HomeassistantMetrics.Reconcile()

		time.Sleep(sleeptime)
	}
}
