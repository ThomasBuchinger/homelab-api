package reconciler_test

import (
	"testing"

	"github.com/thomasbuchinger/homelab-api/pkg/reconciler"
)

const MOCK_METRICS string = `
# Fake Test Metrics
fake_up 1
fake_down 0
fake_pods_total{namespace="ns1",deployment="name"} 2
fake_pods_total{namespace="ns2",deployment="name"} 2

# Real Example Metrics
kube_configmap_info{namespace="cert-manager",configmap="kube-root-ca.crt"} 1
kube_configmap_info{namespace="gotify",configmap="kube-root-ca.crt"} 1
kube_configmap_info{namespace="kube-system",configmap="coredns"} 1
kube_configmap_info{namespace="monitoring",configmap="monitoring-kube-prometheus-persistentvolumesusage"} 1
kube_configmap_info{namespace="spdf",configmap="app-config"} 1
kube_configmap_info{namespace="spdf",configmap="kube-root-ca.crt"} 1
kube_configmap_info{namespace="argocd",configmap="kube-root-ca.crt"} 1
kube_configmap_info{namespace="kube-system",configmap="kube-apiserver-legacy-service-account-token-tracking"} 1
kube_configmap_info{namespace="monitoring",configmap="grafana"} 1
kube_configmap_info{namespace="monitoring",configmap="loki"} 1
kube_configmap_info{namespace="monitoring",configmap="monitoring-kube-prometheus-pod-total"} 1
kube_configmap_info{namespace="monitoring",configmap="monitoring-kube-prometheus-k8s-resources-workload"} 1

kube_pod_status_phase{namespace="paperless",pod="paperless-7b95864585-sfjft",uid="3f0c549c-725a-410b-b867-76f7591b1402",phase="Pending"} 0
kube_pod_status_phase{namespace="paperless",pod="paperless-7b95864585-sfjft",uid="3f0c549c-725a-410b-b867-76f7591b1402",phase="Succeeded"} 0
kube_pod_status_phase{namespace="paperless",pod="paperless-7b95864585-sfjft",uid="3f0c549c-725a-410b-b867-76f7591b1402",phase="Failed"} 0
kube_pod_status_phase{namespace="paperless",pod="paperless-7b95864585-sfjft",uid="3f0c549c-725a-410b-b867-76f7591b1402",phase="Unknown"} 0
kube_pod_status_phase{namespace="paperless",pod="paperless-7b95864585-sfjft",uid="3f0c549c-725a-410b-b867-76f7591b1402",phase="Running"} 1
kube_pod_status_phase{namespace="traefik",pod="traefik-78599579c7-k4wnr",uid="43f933cb-409e-4b50-a7e3-5a85a51637b9",phase="Pending"} 0
kube_pod_status_phase{namespace="traefik",pod="traefik-78599579c7-k4wnr",uid="43f933cb-409e-4b50-a7e3-5a85a51637b9",phase="Succeeded"} 0
kube_pod_status_phase{namespace="traefik",pod="traefik-78599579c7-k4wnr",uid="43f933cb-409e-4b50-a7e3-5a85a51637b9",phase="Failed"} 0
kube_pod_status_phase{namespace="traefik",pod="traefik-78599579c7-k4wnr",uid="43f933cb-409e-4b50-a7e3-5a85a51637b9",phase="Unknown"} 0
kube_pod_status_phase{namespace="traefik",pod="traefik-78599579c7-k4wnr",uid="43f933cb-409e-4b50-a7e3-5a85a51637b9",phase="Running"} 1
`

func Test_MetricsReconciler_is_uninitialied_after_reation(t *testing.T) {
	mr := reconciler.NewDummyMetricsReconciler(MOCK_METRICS)

	if mr.GetSatus() != reconciler.ReconcilerStatusInvalid {
		t.Error("Fresh MetricsReconciler should be invalid")
	}
}

func Test_MetricsReconciler_should_be_ok_after_running_reconcile_once(t *testing.T) {
	mr := reconciler.NewDummyMetricsReconciler(MOCK_METRICS)
	mr.AddMetric("unittest", reconciler.MetricSelect{Name: "fake_pods_total"})
	mr.Reconcile()

	if mr.GetSatus() != reconciler.ReconcilerStatusOK {
		t.Error("Reconciler should be ok")
	}
	if mr.Metrics["unittest"].Value != 4 {
		t.Error("Metric fake_pods_total should be 4. Value: ", mr.Metrics["unittest"].Value)
	}
}

// Function Process Metric

func Test_MetricsReconciler_ProcessMetric_should_read_a_simple_metric(t *testing.T) {
	value, ok := reconciler.ProcessMetric(MOCK_METRICS, reconciler.MetricSelect{Name: "fake_up"})
	if !ok {
		t.Error("should return an OK boolean: ", ok)
	}
	if value != 1 {
		t.Error("Metric fake_up should be 1. Value: ", value)
	}
}

func Test_MetricsReconciler_ProcessMetric_should_return_not_ok_if_metrics_does_not_exist(t *testing.T) {
	_, ok := reconciler.ProcessMetric(MOCK_METRICS, reconciler.MetricSelect{Name: "non_existent_metric"})
	if ok != false {
		t.Error("A non-existent metric should return ok == false")
	}
}

func Test_MetricsReconciler_ProcessMetric_adds_values_if_multiple_metrics_match(t *testing.T) {
	value, _ := reconciler.ProcessMetric(MOCK_METRICS, reconciler.MetricSelect{Name: "fake_pods_total"})
	if value != 4 {
		t.Error("Metric fake_pods_total should be 4. Value: ", value)
	}
}

func Test_MetricsReconciler_ProcessMetric_can_filter_by_label(t *testing.T) {
	value, _ := reconciler.ProcessMetric(MOCK_METRICS, reconciler.MetricSelect{Name: "fake_pods_total", Labels: map[string]string{"namespace": "ns1"}})
	if value != 2 {
		t.Error("Metric fake_pods_total{ns1} should be 2. Value: ", value)
	}
}

func Test_MetricsReconciler_ProcessMetric_matching_labels_should_be_added_even_if_some_labels_differ(t *testing.T) {
	value, _ := reconciler.ProcessMetric(MOCK_METRICS, reconciler.MetricSelect{Name: "fake_pods_total", Labels: map[string]string{"deployment": "name"}})
	if value != 4 {
		t.Error("Metric fake_pods_total{deployment=name} should be 4. Value: ", value)
	}
}

// Function: FindMetricLabels

func Test_MetricsReconciler_FindMetricLabel_should_return_empty_for_no_labels(t *testing.T) {
	labels := reconciler.FindMetricLabels("metric 1")
	if len(labels) != 0 {
		t.Error("Metric 'metric 1' has no labels. Labels: ", labels)
	}
}

func Test_MetricsReconciler_FindMetricLabel_should_return_empty_for_empty_label_list(t *testing.T) {
	labels := reconciler.FindMetricLabels("metric{} 1")
	if len(labels) != 0 {
		t.Error("Metric 'metric 1' has no labels. Labels: ", labels)
	}
}

func Test_MetricsReconciler_FindMetricLabel_should_return_labels_for_a_single_label(t *testing.T) {
	labels := reconciler.FindMetricLabels(`metric{key="value"} 1`)
	if len(labels) != 1 {
		t.Error("Metric 'metric{key=\"value\"}' has 1 label. Labels: ", labels)
	}
	if labels["key"] != "value" {
		t.Error("Label 'key' should be 'value'. Labels: ", labels)
	}
}

func Test_MetricsReconciler_FindMetricLabel_should_return_labels_for_multiple_labels(t *testing.T) {
	labels := reconciler.FindMetricLabels(`metric{key="value",key2="value2"} 1`)
	if labels["key"] != "value" {
		t.Error("Label 'key' should be 'value'. Labels: ", labels)
	}
	if labels["key2"] != "value2" {
		t.Error("Label 'key2' should be 'value'. Labels: ", labels)
	}
}

func Test_MetricsReconciler_FindMetricLabel_should_ignore_missing_quotes(t *testing.T) {
	labels := reconciler.FindMetricLabels(`metric{key=value} 1`)
	if labels["key"] != "value" {
		t.Error("Label 'key' should be 'value'. Labels: ", labels)
	}
}

func Test_MetricsReconciler_FindMetricLabel_should_be_able_to_parse_a_real_example(t *testing.T) {
	labels := reconciler.FindMetricLabels(`kube_pod_status_phase{namespace="paperless",pod="paperless-0",phase="Failed"} 0
`)
	if labels["namespace"] != "paperless" {
		t.Error("Label 'namespace' should be 'paperless'. Labels: ", labels)
	}
	if labels["pod"] != "paperless-0" {
		t.Error("Label 'pod' should be 'paperless-0'. Labels: ", labels)
	}
	if labels["phase"] != "Failed" {
		t.Error("Label 'phase' should be 'failed'. Labels: ", labels)
	}
}
