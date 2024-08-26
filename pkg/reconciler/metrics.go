package reconciler

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
)

type MetricsReconciler struct {
	Logger zap.SugaredLogger

	Status           ReconcilerStatus
	Reason           string
	MetricsUrl       string
	MockScrapeResult string // If this string exists, do not scrape for real, but use string as scrape-result

	Metrics map[string]MetricSelect
}

type MetricSelect struct {
	Name   string
	Labels map[string]string
	Value  float64
}

func NewMetricsReconciler(metrics_url string) *MetricsReconciler {
	conf := common.GetServerConfig()
	return &MetricsReconciler{
		Logger: *conf.RootLogger.Named("MetricsReconciler"),

		Status:           ReconcilerStatusInvalid,
		Reason:           "Did not Run",
		MetricsUrl:       metrics_url,
		MockScrapeResult: "",

		Metrics: make(map[string]MetricSelect),
	}
}

func NewDummyMetricsReconciler(content string) *MetricsReconciler {
	ret := NewMetricsReconciler("http://metrics.example.com")
	ret.MockScrapeResult = content
	return ret
}
func scrapeTarget(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ProcessMetric(scrapeResult string, metric MetricSelect) (float64, bool) {
	var ret_value float64 = 0
	var ok bool = false

	for _, v := range strings.Split(scrapeResult, "\n") {
		if strings.HasPrefix(v, "#") || strings.TrimSpace(v) == "" || len(strings.Split(v, " ")) != 2 {
			continue
		}

		name := strings.TrimSpace(strings.Split(v, " ")[0])
		value, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(v, " ")[1]), 64)
		if err != nil {
			continue
		}
		if strings.Contains(name, metric.Name) {
			discovered_labels := FindMetricLabels(name)
			all_checks_passed := true

			for k, expected := range metric.Labels {
				if val, ok := discovered_labels[k]; !ok || val != expected {
					all_checks_passed = false
				}
			}
			if all_checks_passed {
				ret_value += value
				ok = true
			}

		}
	}

	return ret_value, ok
}

func FindMetricLabels(metric_line string) map[string]string {
	ret := make(map[string]string)

	match := regexp.MustCompile(`\{.+=.+(,.+=.+)*\}`).FindAllString(metric_line, 1)
	if len(match) == 0 {
		return ret
	}
	labels := strings.Split(strings.TrimRight(strings.TrimLeft(match[0], "{"), "}"), ",")

	for _, kvpair := range labels {
		label_name := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.Split(kvpair, "=")[0]), `"`), `"`)
		label_value := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(strings.Split(kvpair, "=")[1]), `"`), `"`)
		ret[label_name] = label_value
	}
	return ret
}

func (r *MetricsReconciler) GetSatus() ReconcilerStatus {
	return r.Status
}
func (r *MetricsReconciler) GetReason() string {
	return r.Reason
}

func (r *MetricsReconciler) AddMetric(key string, metric MetricSelect) *MetricsReconciler {
	r.Metrics[key] = metric
	return r
}

func (r *MetricsReconciler) SetMetricValue(key string, value float64) {
	oldValue := r.Metrics[key]
	oldValue.Value = value
	r.Metrics[key] = oldValue
}

func (r *MetricsReconciler) Reconcile() (bool, int) {
	var err error = nil
	var scrapeResult string = ""

	if r.MockScrapeResult == "" {
		r.Logger.Debugln("Fetching Metrics from: ", r.MetricsUrl)
		scrapeResult, err = scrapeTarget(r.MetricsUrl)
	} else {
		r.Logger.Debugln("Using MockScrapeResult")
		scrapeResult = r.MockScrapeResult
	}
	if err != nil {
		return false, 60
	}

	r.Status = ReconcilerStatusOK
	r.Reason = "ok"

	// var ok bool
	for key, metric := range r.Metrics {
		value, _ := ProcessMetric(scrapeResult, metric)
		r.Logger.Debugln("Metric: ", key, "Value: ", value)
		r.SetMetricValue(key, value)
	}

	return true, 600
}

const TMP_METRICS = `

kube_persistentvolume_status_phase{persistentvolume="syncthing-root",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root",phase="Bound"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root",phase="Released"} 1
kube_persistentvolume_status_phase{persistentvolume="syncthing-root",phase="Failed"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root-nasv3",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root-nasv3",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root-nasv3",phase="Bound"} 1
kube_persistentvolume_status_phase{persistentvolume="syncthing-root-nasv3",phase="Released"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-root-nasv3",phase="Failed"} 0
kube_persistentvolume_status_phase{persistentvolume="paperless-nfs-nasv3",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="paperless-nfs-nasv3",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="paperless-nfs-nasv3",phase="Bound"} 1
kube_persistentvolume_status_phase{persistentvolume="paperless-nfs-nasv3",phase="Released"} 0
kube_persistentvolume_status_phase{persistentvolume="paperless-nfs-nasv3",phase="Failed"} 0
kube_persistentvolume_status_phase{persistentvolume="pvc-6660f756-fc44-415c-b51d-b5c0ffaa96d9",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="pvc-6660f756-fc44-415c-b51d-b5c0ffaa96d9",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="pvc-6660f756-fc44-415c-b51d-b5c0ffaa96d9",phase="Bound"} 1
kube_persistentvolume_status_phase{persistentvolume="pvc-6660f756-fc44-415c-b51d-b5c0ffaa96d9",phase="Released"} 0
kube_persistentvolume_status_phase{persistentvolume="pvc-6660f756-fc44-415c-b51d-b5c0ffaa96d9",phase="Failed"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13",phase="Bound"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13",phase="Released"} 1
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13",phase="Failed"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13-nasv3",phase="Pending"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13-nasv3",phase="Available"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13-nasv3",phase="Bound"} 1
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13-nasv3",phase="Released"} 0
kube_persistentvolume_status_phase{persistentvolume="syncthing-bs13-nasv3",phase="Failed"} 0

kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-state-metrics-58fd4447c6-rzz5z",uid="77732041-59e7-46e7-af8a-5b8faadd6a4b",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="promtail-fdln7",uid="0f8fb597-1f8e-43e2-b819-818934a9e35c",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="promtail-fdln7",uid="0f8fb597-1f8e-43e2-b819-818934a9e35c",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="promtail-fdln7",uid="0f8fb597-1f8e-43e2-b819-818934a9e35c",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="promtail-fdln7",uid="0f8fb597-1f8e-43e2-b819-818934a9e35c",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="promtail-fdln7",uid="0f8fb597-1f8e-43e2-b819-818934a9e35c",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-application-controller-0",uid="0ba04e98-e8aa-451c-b363-184991c3945d",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-application-controller-0",uid="0ba04e98-e8aa-451c-b363-184991c3945d",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-application-controller-0",uid="0ba04e98-e8aa-451c-b363-184991c3945d",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-application-controller-0",uid="0ba04e98-e8aa-451c-b363-184991c3945d",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-application-controller-0",uid="0ba04e98-e8aa-451c-b363-184991c3945d",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-applicationset-controller-595fbb8ff4-d98tm",uid="1325f3e2-67ce-455c-bb53-bc47d61964e6",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-applicationset-controller-595fbb8ff4-d98tm",uid="1325f3e2-67ce-455c-bb53-bc47d61964e6",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-applicationset-controller-595fbb8ff4-d98tm",uid="1325f3e2-67ce-455c-bb53-bc47d61964e6",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-applicationset-controller-595fbb8ff4-d98tm",uid="1325f3e2-67ce-455c-bb53-bc47d61964e6",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-applicationset-controller-595fbb8ff4-d98tm",uid="1325f3e2-67ce-455c-bb53-bc47d61964e6",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-server-7b4ff9897b-zccp9",uid="ba505fa9-e715-4025-869a-31b1fd9f4806",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-server-7b4ff9897b-zccp9",uid="ba505fa9-e715-4025-869a-31b1fd9f4806",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-server-7b4ff9897b-zccp9",uid="ba505fa9-e715-4025-869a-31b1fd9f4806",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-server-7b4ff9897b-zccp9",uid="ba505fa9-e715-4025-869a-31b1fd9f4806",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-server-7b4ff9897b-zccp9",uid="ba505fa9-e715-4025-869a-31b1fd9f4806",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="kube-scheduler-prod",uid="25f09dbe-4b7e-4d80-8605-23d39ca9d234",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-scheduler-prod",uid="25f09dbe-4b7e-4d80-8605-23d39ca9d234",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-scheduler-prod",uid="25f09dbe-4b7e-4d80-8605-23d39ca9d234",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-scheduler-prod",uid="25f09dbe-4b7e-4d80-8605-23d39ca9d234",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-scheduler-prod",uid="25f09dbe-4b7e-4d80-8605-23d39ca9d234",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-gateway-6b75df658-56t8n",uid="6eab7ca5-4118-4429-b2c6-260709b27fb3",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-gateway-6b75df658-56t8n",uid="6eab7ca5-4118-4429-b2c6-260709b27fb3",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-gateway-6b75df658-56t8n",uid="6eab7ca5-4118-4429-b2c6-260709b27fb3",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-gateway-6b75df658-56t8n",uid="6eab7ca5-4118-4429-b2c6-260709b27fb3",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-gateway-6b75df658-56t8n",uid="6eab7ca5-4118-4429-b2c6-260709b27fb3",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-prometheus-operator-6cb8cc898-lftmd",uid="0c5cf962-eed5-4a58-9fce-f1b42e8c1924",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-prometheus-operator-6cb8cc898-lftmd",uid="0c5cf962-eed5-4a58-9fce-f1b42e8c1924",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-prometheus-operator-6cb8cc898-lftmd",uid="0c5cf962-eed5-4a58-9fce-f1b42e8c1924",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-prometheus-operator-6cb8cc898-lftmd",uid="0c5cf962-eed5-4a58-9fce-f1b42e8c1924",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-kube-prometheus-operator-6cb8cc898-lftmd",uid="0c5cf962-eed5-4a58-9fce-f1b42e8c1924",phase="Running"} 1
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
kube_pod_status_phase{namespace="external-http",pod="external-http-7df4b65568-sxcwv",uid="8dd9f3a9-6038-40ad-8217-034c931094f1",phase="Pending"} 0
kube_pod_status_phase{namespace="external-http",pod="external-http-7df4b65568-sxcwv",uid="8dd9f3a9-6038-40ad-8217-034c931094f1",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-http",pod="external-http-7df4b65568-sxcwv",uid="8dd9f3a9-6038-40ad-8217-034c931094f1",phase="Failed"} 0
kube_pod_status_phase{namespace="external-http",pod="external-http-7df4b65568-sxcwv",uid="8dd9f3a9-6038-40ad-8217-034c931094f1",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-http",pod="external-http-7df4b65568-sxcwv",uid="8dd9f3a9-6038-40ad-8217-034c931094f1",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="kube-flannel-rtpcf",uid="04134da0-4276-4caf-bc24-704f02625caf",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-flannel-rtpcf",uid="04134da0-4276-4caf-bc24-704f02625caf",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-flannel-rtpcf",uid="04134da0-4276-4caf-bc24-704f02625caf",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-flannel-rtpcf",uid="04134da0-4276-4caf-bc24-704f02625caf",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-flannel-rtpcf",uid="04134da0-4276-4caf-bc24-704f02625caf",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="alertmanager-monitoring-kube-prometheus-alertmanager-0",uid="f59f02e2-6a46-41f0-b5df-458585c1cb1f",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="alertmanager-monitoring-kube-prometheus-alertmanager-0",uid="f59f02e2-6a46-41f0-b5df-458585c1cb1f",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="alertmanager-monitoring-kube-prometheus-alertmanager-0",uid="f59f02e2-6a46-41f0-b5df-458585c1cb1f",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="alertmanager-monitoring-kube-prometheus-alertmanager-0",uid="f59f02e2-6a46-41f0-b5df-458585c1cb1f",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="alertmanager-monitoring-kube-prometheus-alertmanager-0",uid="f59f02e2-6a46-41f0-b5df-458585c1cb1f",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-chunks-cache-0",uid="2b42cab3-96f9-4358-8334-ed01053ad3a7",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-chunks-cache-0",uid="2b42cab3-96f9-4358-8334-ed01053ad3a7",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-chunks-cache-0",uid="2b42cab3-96f9-4358-8334-ed01053ad3a7",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-chunks-cache-0",uid="2b42cab3-96f9-4358-8334-ed01053ad3a7",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-chunks-cache-0",uid="2b42cab3-96f9-4358-8334-ed01053ad3a7",phase="Running"} 1
kube_pod_status_phase{namespace="external-homelabapi",pod="external-homelab-api-774cf74bf6-nkbwj",uid="bc1db237-6eec-4edb-b27e-f18d3f6d92f4",phase="Pending"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="external-homelab-api-774cf74bf6-nkbwj",uid="bc1db237-6eec-4edb-b27e-f18d3f6d92f4",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="external-homelab-api-774cf74bf6-nkbwj",uid="bc1db237-6eec-4edb-b27e-f18d3f6d92f4",phase="Failed"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="external-homelab-api-774cf74bf6-nkbwj",uid="bc1db237-6eec-4edb-b27e-f18d3f6d92f4",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="external-homelab-api-774cf74bf6-nkbwj",uid="bc1db237-6eec-4edb-b27e-f18d3f6d92f4",phase="Running"} 1
kube_pod_status_phase{namespace="gotify",pod="ntfy-5f978bb5fd-gld9q",uid="064eeaf5-a97d-479f-8671-251349954dc8",phase="Pending"} 0
kube_pod_status_phase{namespace="gotify",pod="ntfy-5f978bb5fd-gld9q",uid="064eeaf5-a97d-479f-8671-251349954dc8",phase="Succeeded"} 0
kube_pod_status_phase{namespace="gotify",pod="ntfy-5f978bb5fd-gld9q",uid="064eeaf5-a97d-479f-8671-251349954dc8",phase="Failed"} 0
kube_pod_status_phase{namespace="gotify",pod="ntfy-5f978bb5fd-gld9q",uid="064eeaf5-a97d-479f-8671-251349954dc8",phase="Unknown"} 0
kube_pod_status_phase{namespace="gotify",pod="ntfy-5f978bb5fd-gld9q",uid="064eeaf5-a97d-479f-8671-251349954dc8",phase="Running"} 1
kube_pod_status_phase{namespace="kubelet-serving-cert-approver",pod="kubelet-serving-cert-approver-7ddfd4469b-9nfjk",uid="4d97c4f8-dac6-4ba1-b83a-8e1632ab4eea",phase="Pending"} 0
kube_pod_status_phase{namespace="kubelet-serving-cert-approver",pod="kubelet-serving-cert-approver-7ddfd4469b-9nfjk",uid="4d97c4f8-dac6-4ba1-b83a-8e1632ab4eea",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kubelet-serving-cert-approver",pod="kubelet-serving-cert-approver-7ddfd4469b-9nfjk",uid="4d97c4f8-dac6-4ba1-b83a-8e1632ab4eea",phase="Failed"} 0
kube_pod_status_phase{namespace="kubelet-serving-cert-approver",pod="kubelet-serving-cert-approver-7ddfd4469b-9nfjk",uid="4d97c4f8-dac6-4ba1-b83a-8e1632ab4eea",phase="Unknown"} 0
kube_pod_status_phase{namespace="kubelet-serving-cert-approver",pod="kubelet-serving-cert-approver-7ddfd4469b-9nfjk",uid="4d97c4f8-dac6-4ba1-b83a-8e1632ab4eea",phase="Running"} 1
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-webhook-6cc4bd4fd4-hj2kz",uid="bd31d68b-2948-43d9-a80d-299ab1a96ed6",phase="Pending"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-webhook-6cc4bd4fd4-hj2kz",uid="bd31d68b-2948-43d9-a80d-299ab1a96ed6",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-webhook-6cc4bd4fd4-hj2kz",uid="bd31d68b-2948-43d9-a80d-299ab1a96ed6",phase="Failed"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-webhook-6cc4bd4fd4-hj2kz",uid="bd31d68b-2948-43d9-a80d-299ab1a96ed6",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-webhook-6cc4bd4fd4-hj2kz",uid="bd31d68b-2948-43d9-a80d-299ab1a96ed6",phase="Running"} 1
kube_pod_status_phase{namespace="spdf",pod="stirling-pdf-5d459cd97-7rnkj",uid="620d5a02-6139-4bf5-9a11-38b2fd9dcb1e",phase="Pending"} 0
kube_pod_status_phase{namespace="spdf",pod="stirling-pdf-5d459cd97-7rnkj",uid="620d5a02-6139-4bf5-9a11-38b2fd9dcb1e",phase="Succeeded"} 0
kube_pod_status_phase{namespace="spdf",pod="stirling-pdf-5d459cd97-7rnkj",uid="620d5a02-6139-4bf5-9a11-38b2fd9dcb1e",phase="Failed"} 0
kube_pod_status_phase{namespace="spdf",pod="stirling-pdf-5d459cd97-7rnkj",uid="620d5a02-6139-4bf5-9a11-38b2fd9dcb1e",phase="Unknown"} 0
kube_pod_status_phase{namespace="spdf",pod="stirling-pdf-5d459cd97-7rnkj",uid="620d5a02-6139-4bf5-9a11-38b2fd9dcb1e",phase="Running"} 1
kube_pod_status_phase{namespace="syncthing",pod="syncthing-6b85f6c797-2d8mt",uid="749430a0-059c-492f-8983-fdfc87c780de",phase="Pending"} 0
kube_pod_status_phase{namespace="syncthing",pod="syncthing-6b85f6c797-2d8mt",uid="749430a0-059c-492f-8983-fdfc87c780de",phase="Succeeded"} 0
kube_pod_status_phase{namespace="syncthing",pod="syncthing-6b85f6c797-2d8mt",uid="749430a0-059c-492f-8983-fdfc87c780de",phase="Failed"} 0
kube_pod_status_phase{namespace="syncthing",pod="syncthing-6b85f6c797-2d8mt",uid="749430a0-059c-492f-8983-fdfc87c780de",phase="Unknown"} 0
kube_pod_status_phase{namespace="syncthing",pod="syncthing-6b85f6c797-2d8mt",uid="749430a0-059c-492f-8983-fdfc87c780de",phase="Running"} 1
kube_pod_status_phase{namespace="external-homelabapi",pod="internal-homelab-api-b6ff6b777-hkls9",uid="f1df149f-5899-4845-89d8-cfad6d626f28",phase="Pending"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="internal-homelab-api-b6ff6b777-hkls9",uid="f1df149f-5899-4845-89d8-cfad6d626f28",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="internal-homelab-api-b6ff6b777-hkls9",uid="f1df149f-5899-4845-89d8-cfad6d626f28",phase="Failed"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="internal-homelab-api-b6ff6b777-hkls9",uid="f1df149f-5899-4845-89d8-cfad6d626f28",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-homelabapi",pod="internal-homelab-api-b6ff6b777-hkls9",uid="f1df149f-5899-4845-89d8-cfad6d626f28",phase="Running"} 1
kube_pod_status_phase{namespace="external-r2",pod="s3manager-55cf4f5547-k4dhk",uid="3f76b182-8186-43ac-9fad-a14e919fda69",phase="Pending"} 0
kube_pod_status_phase{namespace="external-r2",pod="s3manager-55cf4f5547-k4dhk",uid="3f76b182-8186-43ac-9fad-a14e919fda69",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-r2",pod="s3manager-55cf4f5547-k4dhk",uid="3f76b182-8186-43ac-9fad-a14e919fda69",phase="Failed"} 0
kube_pod_status_phase{namespace="external-r2",pod="s3manager-55cf4f5547-k4dhk",uid="3f76b182-8186-43ac-9fad-a14e919fda69",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-r2",pod="s3manager-55cf4f5547-k4dhk",uid="3f76b182-8186-43ac-9fad-a14e919fda69",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="kube-apiserver-prod",uid="ae191b48-6b6f-42c5-9f36-1133d85c31ee",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-apiserver-prod",uid="ae191b48-6b6f-42c5-9f36-1133d85c31ee",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-apiserver-prod",uid="ae191b48-6b6f-42c5-9f36-1133d85c31ee",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-apiserver-prod",uid="ae191b48-6b6f-42c5-9f36-1133d85c31ee",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-apiserver-prod",uid="ae191b48-6b6f-42c5-9f36-1133d85c31ee",phase="Running"} 1
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-cert-controller-7d675fdf6-mw6dr",uid="8590f198-ed3e-4bb6-8d8b-558b42509683",phase="Pending"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-cert-controller-7d675fdf6-mw6dr",uid="8590f198-ed3e-4bb6-8d8b-558b42509683",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-cert-controller-7d675fdf6-mw6dr",uid="8590f198-ed3e-4bb6-8d8b-558b42509683",phase="Failed"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-cert-controller-7d675fdf6-mw6dr",uid="8590f198-ed3e-4bb6-8d8b-558b42509683",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-cert-controller-7d675fdf6-mw6dr",uid="8590f198-ed3e-4bb6-8d8b-558b42509683",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="prometheus-monitoring-kube-prometheus-prometheus-0",uid="38a7a2a8-1e3b-46ad-a57b-027af692202f",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="prometheus-monitoring-kube-prometheus-prometheus-0",uid="38a7a2a8-1e3b-46ad-a57b-027af692202f",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="prometheus-monitoring-kube-prometheus-prometheus-0",uid="38a7a2a8-1e3b-46ad-a57b-027af692202f",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="prometheus-monitoring-kube-prometheus-prometheus-0",uid="38a7a2a8-1e3b-46ad-a57b-027af692202f",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="prometheus-monitoring-kube-prometheus-prometheus-0",uid="38a7a2a8-1e3b-46ad-a57b-027af692202f",phase="Running"} 1
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-5859d8dc69-2svf5",uid="284550f0-56b1-4638-bb0d-de67f9f694f8",phase="Pending"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-5859d8dc69-2svf5",uid="284550f0-56b1-4638-bb0d-de67f9f694f8",phase="Succeeded"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-5859d8dc69-2svf5",uid="284550f0-56b1-4638-bb0d-de67f9f694f8",phase="Failed"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-5859d8dc69-2svf5",uid="284550f0-56b1-4638-bb0d-de67f9f694f8",phase="Unknown"} 0
kube_pod_status_phase{namespace="external-secrets",pod="external-secrets-5859d8dc69-2svf5",uid="284550f0-56b1-4638-bb0d-de67f9f694f8",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-dex-server-fdddcc54d-ft56k",uid="e72a706f-1658-4e4b-9b3e-48c2f142b365",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-dex-server-fdddcc54d-ft56k",uid="e72a706f-1658-4e4b-9b3e-48c2f142b365",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-dex-server-fdddcc54d-ft56k",uid="e72a706f-1658-4e4b-9b3e-48c2f142b365",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-dex-server-fdddcc54d-ft56k",uid="e72a706f-1658-4e4b-9b3e-48c2f142b365",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-dex-server-fdddcc54d-ft56k",uid="e72a706f-1658-4e4b-9b3e-48c2f142b365",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-redis-6976fc7dfc-jmpcc",uid="62cbbdb6-160d-4984-a9f5-12a3258f9bea",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-redis-6976fc7dfc-jmpcc",uid="62cbbdb6-160d-4984-a9f5-12a3258f9bea",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-redis-6976fc7dfc-jmpcc",uid="62cbbdb6-160d-4984-a9f5-12a3258f9bea",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-redis-6976fc7dfc-jmpcc",uid="62cbbdb6-160d-4984-a9f5-12a3258f9bea",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-redis-6976fc7dfc-jmpcc",uid="62cbbdb6-160d-4984-a9f5-12a3258f9bea",phase="Running"} 1
kube_pod_status_phase{namespace="argocd",pod="argocd-repo-server-75c9f5866f-hdlj9",uid="f432df2a-83dd-47bf-b9c5-8d8136b85088",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-repo-server-75c9f5866f-hdlj9",uid="f432df2a-83dd-47bf-b9c5-8d8136b85088",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-repo-server-75c9f5866f-hdlj9",uid="f432df2a-83dd-47bf-b9c5-8d8136b85088",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-repo-server-75c9f5866f-hdlj9",uid="f432df2a-83dd-47bf-b9c5-8d8136b85088",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-repo-server-75c9f5866f-hdlj9",uid="f432df2a-83dd-47bf-b9c5-8d8136b85088",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="cilium-v6pmm",uid="d29e6abc-eb13-46a8-a773-d5f50dd5c236",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-v6pmm",uid="d29e6abc-eb13-46a8-a773-d5f50dd5c236",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-v6pmm",uid="d29e6abc-eb13-46a8-a773-d5f50dd5c236",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-v6pmm",uid="d29e6abc-eb13-46a8-a773-d5f50dd5c236",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-v6pmm",uid="d29e6abc-eb13-46a8-a773-d5f50dd5c236",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="kube-controller-manager-prod",uid="360ad1b3-53c1-4929-805e-5d44a2bad3d0",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-controller-manager-prod",uid="360ad1b3-53c1-4929-805e-5d44a2bad3d0",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-controller-manager-prod",uid="360ad1b3-53c1-4929-805e-5d44a2bad3d0",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-controller-manager-prod",uid="360ad1b3-53c1-4929-805e-5d44a2bad3d0",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-controller-manager-prod",uid="360ad1b3-53c1-4929-805e-5d44a2bad3d0",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="grafana-74d88d5968-jmnb8",uid="2c061e63-b3c1-4faf-a88b-ea916e66f092",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="grafana-74d88d5968-jmnb8",uid="2c061e63-b3c1-4faf-a88b-ea916e66f092",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="grafana-74d88d5968-jmnb8",uid="2c061e63-b3c1-4faf-a88b-ea916e66f092",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="grafana-74d88d5968-jmnb8",uid="2c061e63-b3c1-4faf-a88b-ea916e66f092",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="grafana-74d88d5968-jmnb8",uid="2c061e63-b3c1-4faf-a88b-ea916e66f092",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="cilium-operator-f4d6bd795-vxwgt",uid="e9453061-148b-494d-868f-d163aead18bd",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-operator-f4d6bd795-vxwgt",uid="e9453061-148b-494d-868f-d163aead18bd",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-operator-f4d6bd795-vxwgt",uid="e9453061-148b-494d-868f-d163aead18bd",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-operator-f4d6bd795-vxwgt",uid="e9453061-148b-494d-868f-d163aead18bd",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="cilium-operator-f4d6bd795-vxwgt",uid="e9453061-148b-494d-868f-d163aead18bd",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-t94wg",uid="9bf8a32b-63ec-4cc6-9064-4cc183ce086d",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-t94wg",uid="9bf8a32b-63ec-4cc6-9064-4cc183ce086d",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-t94wg",uid="9bf8a32b-63ec-4cc6-9064-4cc183ce086d",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-t94wg",uid="9bf8a32b-63ec-4cc6-9064-4cc183ce086d",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-t94wg",uid="9bf8a32b-63ec-4cc6-9064-4cc183ce086d",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-0",uid="840c275c-5ea3-41ec-953d-fda4ad59c5c0",phase="Pending"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-0",uid="840c275c-5ea3-41ec-953d-fda4ad59c5c0",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-0",uid="840c275c-5ea3-41ec-953d-fda4ad59c5c0",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-0",uid="840c275c-5ea3-41ec-953d-fda4ad59c5c0",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-0",uid="840c275c-5ea3-41ec-953d-fda4ad59c5c0",phase="Running"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-notifications-controller-8654456db7-bh9nj",uid="38259529-70cd-445b-ac68-8e20389bdbf1",phase="Pending"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-notifications-controller-8654456db7-bh9nj",uid="38259529-70cd-445b-ac68-8e20389bdbf1",phase="Succeeded"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-notifications-controller-8654456db7-bh9nj",uid="38259529-70cd-445b-ac68-8e20389bdbf1",phase="Failed"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-notifications-controller-8654456db7-bh9nj",uid="38259529-70cd-445b-ac68-8e20389bdbf1",phase="Unknown"} 0
kube_pod_status_phase{namespace="argocd",pod="argocd-notifications-controller-8654456db7-bh9nj",uid="38259529-70cd-445b-ac68-8e20389bdbf1",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-6kg8m",uid="c6590dc0-5079-4da1-b4b0-c35b1a7b5e94",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-6kg8m",uid="c6590dc0-5079-4da1-b4b0-c35b1a7b5e94",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-6kg8m",uid="c6590dc0-5079-4da1-b4b0-c35b1a7b5e94",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-6kg8m",uid="c6590dc0-5079-4da1-b4b0-c35b1a7b5e94",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="coredns-c78fdf99-6kg8m",uid="c6590dc0-5079-4da1-b4b0-c35b1a7b5e94",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="taint-cleaner-28744925-th6lk",uid="39b31d4a-5c45-4de6-8ea3-49e80363fd9e",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="taint-cleaner-28744925-th6lk",uid="39b31d4a-5c45-4de6-8ea3-49e80363fd9e",phase="Succeeded"} 1
kube_pod_status_phase{namespace="kube-system",pod="taint-cleaner-28744925-th6lk",uid="39b31d4a-5c45-4de6-8ea3-49e80363fd9e",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="taint-cleaner-28744925-th6lk",uid="39b31d4a-5c45-4de6-8ea3-49e80363fd9e",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="taint-cleaner-28744925-th6lk",uid="39b31d4a-5c45-4de6-8ea3-49e80363fd9e",phase="Running"} 0
kube_pod_status_phase{namespace="local-path-storage",pod="local-path-provisioner-859c9d8bc6-b8mgd",uid="12f3feba-12c5-41b0-b5ae-7539430b68f3",phase="Pending"} 0
kube_pod_status_phase{namespace="local-path-storage",pod="local-path-provisioner-859c9d8bc6-b8mgd",uid="12f3feba-12c5-41b0-b5ae-7539430b68f3",phase="Succeeded"} 0
kube_pod_status_phase{namespace="local-path-storage",pod="local-path-provisioner-859c9d8bc6-b8mgd",uid="12f3feba-12c5-41b0-b5ae-7539430b68f3",phase="Failed"} 0
kube_pod_status_phase{namespace="local-path-storage",pod="local-path-provisioner-859c9d8bc6-b8mgd",uid="12f3feba-12c5-41b0-b5ae-7539430b68f3",phase="Unknown"} 0
kube_pod_status_phase{namespace="local-path-storage",pod="local-path-provisioner-859c9d8bc6-b8mgd",uid="12f3feba-12c5-41b0-b5ae-7539430b68f3",phase="Running"} 1
kube_pod_status_phase{namespace="paperless",pod="redis-ephemeral-master-0",uid="bb1f1e1f-ae23-4328-997b-94a1b851eb2f",phase="Pending"} 0
kube_pod_status_phase{namespace="paperless",pod="redis-ephemeral-master-0",uid="bb1f1e1f-ae23-4328-997b-94a1b851eb2f",phase="Succeeded"} 0
kube_pod_status_phase{namespace="paperless",pod="redis-ephemeral-master-0",uid="bb1f1e1f-ae23-4328-997b-94a1b851eb2f",phase="Failed"} 0
kube_pod_status_phase{namespace="paperless",pod="redis-ephemeral-master-0",uid="bb1f1e1f-ae23-4328-997b-94a1b851eb2f",phase="Unknown"} 0
kube_pod_status_phase{namespace="paperless",pod="redis-ephemeral-master-0",uid="bb1f1e1f-ae23-4328-997b-94a1b851eb2f",phase="Running"} 1
kube_pod_status_phase{namespace="okd-console",pod="okd-console-7b87954845-kzrxc",uid="e9e1b380-61ef-4a21-bb55-f247dea1f962",phase="Pending"} 0
kube_pod_status_phase{namespace="okd-console",pod="okd-console-7b87954845-kzrxc",uid="e9e1b380-61ef-4a21-bb55-f247dea1f962",phase="Succeeded"} 0
kube_pod_status_phase{namespace="okd-console",pod="okd-console-7b87954845-kzrxc",uid="e9e1b380-61ef-4a21-bb55-f247dea1f962",phase="Failed"} 0
kube_pod_status_phase{namespace="okd-console",pod="okd-console-7b87954845-kzrxc",uid="e9e1b380-61ef-4a21-bb55-f247dea1f962",phase="Unknown"} 0
kube_pod_status_phase{namespace="okd-console",pod="okd-console-7b87954845-kzrxc",uid="e9e1b380-61ef-4a21-bb55-f247dea1f962",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="kube-proxy-w75vh",uid="ceb71d3f-9420-45e7-9ae2-2ce5c3d2a812",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-proxy-w75vh",uid="ceb71d3f-9420-45e7-9ae2-2ce5c3d2a812",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-proxy-w75vh",uid="ceb71d3f-9420-45e7-9ae2-2ce5c3d2a812",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-proxy-w75vh",uid="ceb71d3f-9420-45e7-9ae2-2ce5c3d2a812",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="kube-proxy-w75vh",uid="ceb71d3f-9420-45e7-9ae2-2ce5c3d2a812",phase="Running"} 1
kube_pod_status_phase{namespace="kube-system",pod="metrics-server-7ffbc6d68-rt8ft",uid="5f2b142d-3ad4-4246-9d74-db5badabaf97",phase="Pending"} 0
kube_pod_status_phase{namespace="kube-system",pod="metrics-server-7ffbc6d68-rt8ft",uid="5f2b142d-3ad4-4246-9d74-db5badabaf97",phase="Succeeded"} 0
kube_pod_status_phase{namespace="kube-system",pod="metrics-server-7ffbc6d68-rt8ft",uid="5f2b142d-3ad4-4246-9d74-db5badabaf97",phase="Failed"} 0
kube_pod_status_phase{namespace="kube-system",pod="metrics-server-7ffbc6d68-rt8ft",uid="5f2b142d-3ad4-4246-9d74-db5badabaf97",phase="Unknown"} 0
kube_pod_status_phase{namespace="kube-system",pod="metrics-server-7ffbc6d68-rt8ft",uid="5f2b142d-3ad4-4246-9d74-db5badabaf97",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-canary-ffkkt",uid="f154f3cd-0621-4f02-9abd-2a9057de1796",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-canary-ffkkt",uid="f154f3cd-0621-4f02-9abd-2a9057de1796",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-canary-ffkkt",uid="f154f3cd-0621-4f02-9abd-2a9057de1796",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-canary-ffkkt",uid="f154f3cd-0621-4f02-9abd-2a9057de1796",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-canary-ffkkt",uid="f154f3cd-0621-4f02-9abd-2a9057de1796",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="loki-results-cache-0",uid="ea87700d-5f56-481a-b501-bea89464fedb",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-results-cache-0",uid="ea87700d-5f56-481a-b501-bea89464fedb",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-results-cache-0",uid="ea87700d-5f56-481a-b501-bea89464fedb",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-results-cache-0",uid="ea87700d-5f56-481a-b501-bea89464fedb",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="loki-results-cache-0",uid="ea87700d-5f56-481a-b501-bea89464fedb",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="monitoring-prometheus-node-exporter-gkf4t",uid="767ae934-c3cb-450e-9caa-8c2a25916a05",phase="Pending"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-prometheus-node-exporter-gkf4t",uid="767ae934-c3cb-450e-9caa-8c2a25916a05",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-prometheus-node-exporter-gkf4t",uid="767ae934-c3cb-450e-9caa-8c2a25916a05",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-prometheus-node-exporter-gkf4t",uid="767ae934-c3cb-450e-9caa-8c2a25916a05",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="monitoring-prometheus-node-exporter-gkf4t",uid="767ae934-c3cb-450e-9caa-8c2a25916a05",phase="Running"} 1
kube_pod_status_phase{namespace="monitoring",pod="tempo-0",uid="0fb58858-f5b0-41cd-8f05-59aad53f87d2",phase="Pending"} 1
kube_pod_status_phase{namespace="monitoring",pod="tempo-0",uid="0fb58858-f5b0-41cd-8f05-59aad53f87d2",phase="Succeeded"} 0
kube_pod_status_phase{namespace="monitoring",pod="tempo-0",uid="0fb58858-f5b0-41cd-8f05-59aad53f87d2",phase="Failed"} 0
kube_pod_status_phase{namespace="monitoring",pod="tempo-0",uid="0fb58858-f5b0-41cd-8f05-59aad53f87d2",phase="Unknown"} 0
kube_pod_status_phase{namespace="monitoring",pod="tempo-0",uid="0fb58858-f5b0-41cd-8f05-59aad53f87d2",phase="Running"} 0
`
