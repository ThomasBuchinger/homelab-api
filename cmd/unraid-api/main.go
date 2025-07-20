package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
)

type UnraidDiskMetrics struct {
	Name, ID, Partition   string
	Status                string
	Size                  int64
	Reads, Writes, Errors int64
}

type UnraidMetrics struct {
	ok                           bool
	sbState, sbSynced, sbSynced2 int64
	sbSyncErrs, sbSyncExit       int64

	mdState, mdResyncAction                                                   string
	mdResyncSize, mdResyncCorr, mdResync, mdResyncPos, mdResyncDt, mdResyncDb int64

	Disks map[string]UnraidDiskMetrics
}

var (
	unraidArrayUp     = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_array_up", Help: "Unraid Array ist started"})
	unraidArrayStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_array_status", Help: "Unraid Array Status string"}, []string{"status"})

	unraidParityRuning  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_running", Help: "Parity Check running?"})
	unraidParityLastRun = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_last_run", Help: "EXPERIMENTAL: Timestamp of the Last Parity Check"})
	unraidParityMode    = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_parity_mode", Help: "EXPERIMENTAL: Parity Check mode: check/write-corrections?"}, []string{"mode"})
	unraidParityErrors  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_errors", Help: "EXPERIMENTAL: Parity Errors"})

	unraidReadRunning = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_read_running", Help: "EXPERIMENTAL: ReadCheck running?"})
	unraidReadLastRun = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_read_last_run", Help: "EXPERIMENTAL: Timestamp of the Last ReadCheck"})
	unraidReadErrors  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_read_errors", Help: "EXPERIMENTAL: ReadCheck Errors (do not use, there is conflicting information if my calculations are correct)"})
	unraidReadExit    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_read_exit", Help: "EXPERIMENTAL: ReadCheck Exit Code (0 Success (with no errors?), -4 Aorted)"})

	unraidDiskSize   = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_size", Help: "Disk Size"}, []string{"name", "id", "partition"})
	unraidDiskReads  = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_reads", Help: "Disk Reads"}, []string{"name", "id", "partition"})
	unraidDiskWrites = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_writes", Help: "Disk Writes"}, []string{"name", "id", "partition"})

	unraidDiskErrors = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_errors", Help: "Disk Errors"}, []string{"name", "id", "partition"})
)

func setupUnraidApi(r *gin.Engine) *gin.Engine {

	r.GET("unraid/metrics", func(ctx *gin.Context) { updateMetrics() }, gin.WrapH(promhttp.Handler()))
	r.GET("unraid/mdstat", handleDebug)

	return r
}

func handleDebug(c *gin.Context) {
	data := readFile()
	if len(data) == 0 {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, data)
}

func readFile() string {
	data, err := os.ReadFile(viper.GetString("homelab.unraid.mdstat"))
	if err != nil {
		println(err.Error())
		return ""
	}
	return string(data)
}

func parseMdStatFile(data string) UnraidMetrics {
	metrics := UnraidMetrics{ok: true, Disks: map[string]UnraidDiskMetrics{}}
	getDisk := func(key string) (string, UnraidDiskMetrics) {
		key_arr := strings.Split(key, ".")
		if len(key_arr) == 1 {
			return "", UnraidDiskMetrics{}
		}
		disk, ok := metrics.Disks[key_arr[1]]
		if !ok {
			return key_arr[1], UnraidDiskMetrics{}
		}
		return key_arr[1], disk
	}

	for _, line := range strings.Split(data, "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		line_arr := strings.Split(line, "=")
		key, value := line_arr[0], line_arr[1]
		disk_num, disk := getDisk(key)
		var err error = nil

		if key == "mdState" {
			metrics.mdState = value
		}
		if key == "mdResyncAction" {
			metrics.mdResyncAction = value
		}
		if key == "sbState" {
			metrics.sbState, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "sbSynced" {
			metrics.sbSynced, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "sbSynced2" {
			metrics.sbSynced2, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "sbSyncErrs" {
			metrics.sbSyncErrs, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "sbSyncExit" {
			metrics.sbSyncExit, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResyncSize" {
			metrics.mdResyncSize, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResyncCorr" {
			metrics.mdResyncCorr, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResync" {
			metrics.mdResync, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResyncPos" {
			metrics.mdResyncPos, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResyncDt" {
			metrics.mdResyncDt, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}
		if key == "mdResyncDb" {
			metrics.mdResyncDb, err = strconv.ParseInt(value, 0, 0)
			metrics.ok = metrics.ok && err == nil
		}

		if strings.HasPrefix(key, "rdevId") {
			disk.ID = value
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevName") {
			disk.Name = value
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "diskName") {
			disk.Partition = value
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevStatus") {
			disk.Status = value
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevSize") {
			disk.Size, _ = strconv.ParseInt(value, 0, 0)
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevReads") {
			disk.Reads, _ = strconv.ParseInt(value, 0, 0)
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevWrites") {
			disk.Writes, _ = strconv.ParseInt(value, 0, 0)
			metrics.Disks[disk_num] = disk
		}
		if strings.HasPrefix(key, "rdevNumErrors") {
			disk.Errors, _ = strconv.ParseInt(value, 0, 0)
			metrics.Disks[disk_num] = disk
		}
	}
	return metrics
}

func updateMetrics() {
	metrics := parseMdStatFile(readFile())

	// Reset all metrics to default
	unraidArrayUp.Set(0)
	unraidArrayStatus.WithLabelValues("STARTED").Set(0)
	unraidArrayStatus.WithLabelValues("CHECKING").Set(0)
	unraidArrayStatus.WithLabelValues("RESYNC").Set(0)
	unraidArrayStatus.WithLabelValues("STOPPED").Set(0)

	unraidParityRuning.Set(0)
	unraidParityLastRun.Set(0)
	unraidParityMode.WithLabelValues("check").Set(0)
	unraidParityErrors.Set(0)

	unraidReadRunning.Set(0)
	unraidReadLastRun.Set(0)
	unraidReadErrors.Set(0)

	// Set Metrics
	if metrics.mdState == "STARTED" {
		unraidArrayUp.Set(1)
	}
	unraidArrayStatus.WithLabelValues(metrics.mdState).Set(1)

	if metrics.mdResyncPos != 0 {
		unraidParityRuning.Set(1)
	}
	unraidParityLastRun.Set(float64(metrics.sbSynced))
	unraidParityMode.WithLabelValues(metrics.mdResyncAction).Set(1)
	unraidParityErrors.Set(float64(metrics.mdResyncCorr))

	if metrics.sbSynced2 == 0 {
		unraidReadRunning.Set(1)
	}
	unraidReadLastRun.Set(float64(metrics.sbSynced2))
	unraidReadErrors.Set(float64(metrics.sbSyncErrs))
	unraidReadExit.Set(float64(metrics.sbSyncExit))

	for _, disk := range metrics.Disks {
		unraidDiskSize.WithLabelValues(disk.Name, disk.ID, disk.Partition).Set(float64(disk.Size))
		unraidDiskReads.WithLabelValues(disk.Name, disk.ID, disk.Partition).Set(float64(disk.Reads))
		unraidDiskWrites.WithLabelValues(disk.Name, disk.ID, disk.Partition).Set(float64(disk.Writes))
		unraidDiskErrors.WithLabelValues(disk.Name, disk.ID, disk.Partition).Set(float64(disk.Errors))
	}

}

func main() {
	common.SetupViperConfig()
	serverConfig := common.GetServerConfig()
	router := api.SetupDefaultRouter()
	router = setupUnraidApi(router)

	prometheus.MustRegister(unraidArrayStatus, unraidArrayUp, unraidDiskErrors, unraidDiskReads, unraidDiskSize, unraidDiskWrites, unraidParityErrors, unraidParityLastRun, unraidParityMode, unraidParityRuning, unraidReadErrors, unraidReadExit, unraidReadLastRun, unraidReadRunning)

	router.Run(common.GetEnvWithDefault("BIND_ADDR", serverConfig.BindAddr))
}
