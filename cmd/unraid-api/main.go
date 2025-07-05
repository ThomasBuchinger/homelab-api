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
	ArrayStarted bool
	ArrayStatus  string

	ParityOK          bool
	ParityStatus      int64
	ParityErrors      int64
	ParityExitCode    int64
	ParityCorrections int64
	ParityLastRun     int64

	Disks map[string]UnraidDiskMetrics
}

var (
	unraidArrayUp           = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_array_up", Help: "Unraid Array ist started"})
	unraidArrayStatus       = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_array_status", Help: "Unraid Array Status string"}, []string{"status"})
	unraidParityOk          = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_ok", Help: "Was the last Parity Check without errors"})
	unraidParityStatus      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_pariy_status_code", Help: "Numeric Status code of the Parity check"})
	unraidParityErrors      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_errors", Help: "Number of Parity Errors"})
	unraidParityExitCode    = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_exit_code", Help: "Numeric Exit Code of the Parity Check"})
	unraidParityCorrections = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_corrections", Help: "Number of corrections found in the last Parity Check"})
	unraidParityLastRun     = prometheus.NewGauge(prometheus.GaugeOpts{Name: "unraid_parity_last_run", Help: "Timestamp of the Last Parity Check"})
	unraidDiskSize          = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_size", Help: "Disk Size"}, []string{"name", "id", "partition"})
	unraidDiskReads         = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_reads", Help: "Disk Reads"}, []string{"name", "id", "partition"})
	unraidDiskWrites        = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_writes", Help: "Disk Writes"}, []string{"name", "id", "partition"})
	unraidDiskErrors        = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "unraid_disk_errors", Help: "Disk Errors"}, []string{"name", "id", "partition"})
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
	metrics := UnraidMetrics{Disks: map[string]UnraidDiskMetrics{}}
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
			metrics.ArrayStatus = value
			metrics.ArrayStarted = value == "STARTED"
		}
		if key == "sbState" {
			metrics.ParityStatus, err = strconv.ParseInt(value, 0, 0)
			metrics.ParityOK = err == nil && metrics.ParityStatus == 1
		}
		if key == "sbSyncErrs" {
			metrics.ParityErrors, _ = strconv.ParseInt(value, 0, 0)
		}
		if key == "sbSyncExit" {
			metrics.ParityExitCode, _ = strconv.ParseInt(value, 0, 0)
		}
		if key == "sbSynced2" {
			metrics.ParityLastRun, _ = strconv.ParseInt(value, 0, 0)
		}
		if key == "mdResyncCorr" {
			metrics.ParityCorrections, _ = strconv.ParseInt(value, 0, 0)
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

	if metrics.ArrayStarted {
		unraidArrayUp.Set(1)
	} else {
		unraidArrayUp.Set(0)
	}
	unraidArrayStatus.WithLabelValues(metrics.ArrayStatus).Set(1)

	if metrics.ParityOK {
		unraidParityOk.Set(1)
	} else {
		unraidParityOk.Set(0)
	}
	unraidParityStatus.Set(float64(metrics.ParityStatus))
	unraidParityCorrections.Set(float64(metrics.ParityCorrections))
	unraidParityErrors.Set(float64(metrics.ParityErrors))
	unraidParityExitCode.Set(float64(metrics.ParityExitCode))
	unraidParityLastRun.Set(float64(metrics.ParityLastRun))

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

	prometheus.MustRegister(unraidArrayUp, unraidArrayStatus, unraidDiskErrors, unraidDiskReads, unraidDiskSize, unraidDiskWrites, unraidParityCorrections, unraidParityErrors, unraidParityExitCode, unraidParityLastRun, unraidParityOk, unraidParityStatus)

	router.Run(common.GetEnvWithDefault("BIND_ADDR", serverConfig.BindAddr))
}
