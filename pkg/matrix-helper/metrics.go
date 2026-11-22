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
	metricMatrixUserCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrixhelper_matrix_user_created_total",
		Help: "Total number of Matrix users registered by the matrix-management task",
	})
	metricMatrixRoomsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "matrixhelper_matrix_rooms_created_total",
		Help: "Total number of Matrix rooms created by the matrix-management task",
	})
	metricMatrixConnected = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matrixhelper_matrix_api_connected",
		Help: "1 when the bot is logged in to the Matrix API",
	})
	metricMatrixCheckTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "matrixhelper_matrix_check_timestamp",
		Help: "Unix timestamp of the last completed matrix-management run",
	})
)
