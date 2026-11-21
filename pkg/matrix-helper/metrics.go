package matrixhelper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricChangeTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrixhelper_clientsecret_change_total",
		Help: "Total number of client_secret changes applied to the MAS config secret",
	})
	metricMasRestartTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrixhelper_clientsecret_mas_restart_total",
		Help: "Total number of MAS deployment restarts triggered by the client_secret sync",
	})
	metricCheckTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matrixhelper_clientsecret_check_timestamp",
		Help: "Unix timestamp of the last completed client_secret check",
	})
	metricChangeTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matrixhelper_clientsecret_change_timestamp",
		Help: "Unix timestamp of the last applied client_secret change",
	})
	metricMasRestartTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matrixhelper_clientsecret_mas_restart_timestamp",
		Help: "Unix timestamp of the last MAS deployment restart triggered by the client_secret sync",
	})
)
