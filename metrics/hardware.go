package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	log "github.com/sirupsen/logrus"
)

// Max number of metrics we want to keep.
const maxCollectionValues = 500

func collectCPUUtilization() {
	if len(metrics.CPUUtilizations) > maxCollectionValues {
		metrics.CPUUtilizations = metrics.CPUUtilizations[1:]
	}

	v, err := cpu.Percent(0, false)
	if err != nil {
		log.Errorln(err)
		return
	}

	metricValue := timestampedValue{time.Now(), int(v[0])}
	metrics.CPUUtilizations = append(metrics.CPUUtilizations, metricValue)
	cpuUsage.Set(float64(metricValue.Value))
}

func collectRAMUtilization() {
	if len(metrics.RAMUtilizations) > maxCollectionValues {
		metrics.RAMUtilizations = metrics.RAMUtilizations[1:]
	}

	memoryUsage, _ := mem.VirtualMemory()
	metricValue := timestampedValue{time.Now(), int(memoryUsage.UsedPercent)}
	metrics.RAMUtilizations = append(metrics.RAMUtilizations, metricValue)
}

func collectDiskUtilization() {
	path := "./"
	diskUse, _ := disk.Usage(path)

	if len(metrics.DiskUtilizations) > maxCollectionValues {
		metrics.DiskUtilizations = metrics.DiskUtilizations[1:]
	}

	metricValue := timestampedValue{time.Now(), int(diskUse.UsedPercent)}
	metrics.DiskUtilizations = append(metrics.DiskUtilizations, metricValue)
}
