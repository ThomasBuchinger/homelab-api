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

func ReconcileLoop() {
	sleeptime, _ := time.ParseDuration("1m")
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

	for {
		conf.RootLogger.Logln(zap.DebugLevel, "Running Reconcilers...")

		EvergreenMetric.Reconcile()
		ProdMetric.Reconcile()
		SyncthingMetric.Reconcile()

		time.Sleep(sleeptime)
	}
}
