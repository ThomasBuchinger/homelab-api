package reconciler

import (
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
)

type Reconciler interface {
	GetStatus() ReconcilerStatus
	GetReason() string
	Reconcile() (bool, int)
	// GetData() interface{}
}

type ReconcilerStatus string

const (
	ReconcilerStatusOK                   ReconcilerStatus = "OK"
	ReconcilerStatusInvalid              ReconcilerStatus = "INVALID"
	ReconcilerStatusDown                 ReconcilerStatus = "DOWN"
	ReconcilerStatusWarnReconcilerStatus                  = "WARN"
)

var KubernetesMetricProd *MetricsReconciler
var KubernetesMetricEvergreen *MetricsReconciler

func ReconcileLoop() {
	sleeptime, _ := time.ParseDuration("1m")
	conf := common.GetServerConfig()

	// kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Failed"} 0
	KubernetesMetricProd = NewMetricsReconciler(conf.HomelabEnv.ProdKubeStateMetricsUrl)
	KubernetesMetricProd.AddMetric("num_pod_total", MetricSelect{Name: "kube_pod_status_phase"})
	KubernetesMetricProd.AddMetric("num_pod_succeess", MetricSelect{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Succeeded"}})
	KubernetesMetricProd.AddMetric("num_pod_running", MetricSelect{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Running"}})
	KubernetesMetricProd.AddMetric("num_pvc_total", MetricSelect{Name: "kube_persistentvolume_status_phase"})
	KubernetesMetricProd.AddMetric("num_pvc_bound", MetricSelect{Name: "kube_persistentvolume_status_phase", Labels: map[string]string{"phase": "Bound"}})

	KubernetesMetricEvergreen = NewMetricsReconciler(conf.HomelabEnv.EvergreenKubeStateMetricsUrl)
	KubernetesMetricEvergreen.AddMetric("num_pod_total", MetricSelect{Name: "kube_pod_status_phase"})
	KubernetesMetricEvergreen.AddMetric("num_pod_succeess", MetricSelect{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Succeeded"}})
	KubernetesMetricEvergreen.AddMetric("num_pod_running", MetricSelect{Name: "kube_pod_status_phase", Labels: map[string]string{"phase": "Running"}})
	KubernetesMetricEvergreen.AddMetric("num_pvc_total", MetricSelect{Name: "kube_persistentvolume_status_phase"})
	KubernetesMetricEvergreen.AddMetric("num_pvc_bound", MetricSelect{Name: "kube_persistentvolume_status_phase", Labels: map[string]string{"phase": "Bound"}})

	for {
		conf.RootLogger.Logln(zap.DebugLevel, "Running Reconcilers...")

		KubernetesMetricEvergreen.Reconcile()
		KubernetesMetricProd.Reconcile()

		time.Sleep(sleeptime)
	}
}
