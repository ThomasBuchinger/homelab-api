package metricscraper

import (
	"fmt"
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

	Status             common.ReconcilerStatus
	Reason             string
	MockScrapeResult   string // If this string exists, do not scrape for real, but use string as scrape-result
	MetricsUrl         string
	AuthUser, AuthPass string

	Metrics map[string]Metric
}

const MetricRootKey string = "__root__"

type Metric struct {
	Name   string
	Labels map[string]string
	Value  float64

	GroupBy       string
	GroupedValues map[string]float64
}

func NewMetricsReconciler(metrics_url string) *MetricsReconciler {
	conf := common.GetServerConfig()
	return &MetricsReconciler{
		Logger: *conf.RootLogger.Named("MetricsReconciler"),

		Status:           common.ReconcilerStatusNew,
		Reason:           "Did not Run",
		MetricsUrl:       metrics_url,
		MockScrapeResult: "",
		AuthUser:         "",
		AuthPass:         "",

		Metrics: make(map[string]Metric),
	}
}

func NewDummyMetricsReconciler(content string) *MetricsReconciler {
	ret := NewMetricsReconciler("http://metrics.example.com")
	ret.MockScrapeResult = content
	return ret
}

func ProcessMetric(scrapeResult string, metric Metric) (map[string]float64, bool) {
	ret_map, key, ok := map[string]float64{}, MetricRootKey, false

	for _, v := range strings.Split(scrapeResult, "\n") {
		if strings.HasPrefix(v, "#") || strings.TrimSpace(v) == "" || len(strings.Split(v, " ")) != 2 {
			continue
		}
		key = MetricRootKey // Reset key variable for each line

		name := strings.TrimSpace(strings.Split(v, " ")[0])
		value, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(v, " ")[1]), 64)
		if err != nil {
			continue
		}
		if strings.Contains(name, metric.Name) {
			all_checks_passed := true
			discovered_labels := FindMetricLabels(name)
			groupKey, groupKeyExists := discovered_labels[metric.GroupBy]
			if metric.GroupBy != "" && groupKeyExists {
				key = groupKey
			}

			for k, expected := range metric.Labels {
				if val, ok := discovered_labels[k]; !ok || val != expected {
					all_checks_passed = false
				}

			}
			if all_checks_passed {
				ret_map[key] += value
				ok = true
			}

		}
	}

	return ret_map, ok
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

func (r *MetricsReconciler) GetSatus() common.ReconcilerStatus {
	return r.Status
}
func (r *MetricsReconciler) GetReason() string {
	return r.Reason
}

func (r *MetricsReconciler) AddMetric(key string, metric Metric) *MetricsReconciler {
	r.Metrics[key] = metric
	return r
}

func (r *MetricsReconciler) SetMetricValue(key string, value_map map[string]float64) {
	oldValue := r.Metrics[key]
	for _, v := range value_map {
		oldValue.Value += v
	}
	oldValue.GroupedValues = value_map
	r.Metrics[key] = oldValue
}

func (r *MetricsReconciler) ResetMetric(key string) {
	r.Metrics[key] = Metric{
		Name:          r.Metrics[key].Name,
		Labels:        r.Metrics[key].Labels,
		GroupBy:       r.Metrics[key].GroupBy,
		Value:         0,
		GroupedValues: make(map[string]float64),
	}
}

func scrapeTarget(url string, user string, pass string) (string, common.ReconcilerStatus, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return "", common.ReconcilerStatusConfigInvalid, err
	}
	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", common.ReconcilerStatusDown, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", common.ReconcilerStatusError, fmt.Errorf("%s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", common.ReconcilerStatusError, err
	}

	return string(body), common.ReconcilerStatusOK, nil
}

func (r *MetricsReconciler) Reconcile() (bool, int) {
	var err error = nil
	var scrapeResult string = ""
	var status common.ReconcilerStatus
	for key := range r.Metrics {
		r.ResetMetric(key)
	}

	if r.MockScrapeResult == "" {
		r.Logger.Debugln("Fetching Metrics from: ", r.MetricsUrl)
		scrapeResult, status, err = scrapeTarget(r.MetricsUrl, r.AuthUser, r.AuthPass)
	} else {
		r.Logger.Debugln("Using MockScrapeResult")
		scrapeResult = r.MockScrapeResult
		status = common.ReconcilerStatusOK
	}
	if err != nil {
		r.Logger.Errorf("Failed to scrape target '%s': %s", r.MetricsUrl, err.Error())
		r.Status = status
		r.Reason = err.Error()
		return false, 60
	}

	// var ok bool
	for key, metric := range r.Metrics {
		value_map, _ := ProcessMetric(scrapeResult, metric)

		r.Logger.Debugln("Metric: ", key, "Value: ", value_map)
		r.SetMetricValue(key, value_map)
	}
	r.Status = common.ReconcilerStatusOK
	r.Reason = "ok"

	return true, 600
}
