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
	if len(Metrics.CPUUtilizations) > maxCollectionValues {
		Metrics.CPUUtilizations = Metrics.CPUUtilizations[1:]
	}

	v, err := cpu.Percent(0, false)
	if err != nil {
		log.Errorln(err)
		return
	}

	metricValue := timestampedValue{time.Now(), int(v[0])}
	Metrics.CPUUtilizations = append(Metrics.CPUUtilizations, metricValue)
	cpuUsage.Set(float64(metricValue.Value))
}

func collectRAMUtilization() {
	if len(Metrics.RAMUtilizations) > maxCollectionValues {
		Metrics.RAMUtilizations = Metrics.RAMUtilizations[1:]
	}

	memoryUsage, _ := mem.VirtualMemory()
	metricValue := timestampedValue{time.Now(), int(memoryUsage.UsedPercent)}
	Metrics.RAMUtilizations = append(Metrics.RAMUtilizations, metricValue)
}

func collectDiskUtilization() {
	path := "./"
	diskUse, _ := disk.Usage(path)

	if len(Metrics.DiskUtilizations) > maxCollectionValues {
		Metrics.DiskUtilizations = Metrics.DiskUtilizations[1:]
	}

	metricValue := timestampedValue{time.Now(), int(diskUse.UsedPercent)}
	Metrics.DiskUtilizations = append(Metrics.DiskUtilizations, metricValue)
}
