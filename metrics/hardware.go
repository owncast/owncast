package metrics

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Max number of metrics we want to keep.
const maxCollectionValues = 500

func collectCPUUtilization() {
	if len(Metrics.CPUUtilizations) > maxCollectionValues {
		Metrics.CPUUtilizations = Metrics.CPUUtilizations[1:]
	}

	v, err := cpu.Percent(0, false)
	if err != nil {
		panic(err)
	}

	metricValue := timestampedValue{time.Now(), int(v[0])}
	Metrics.CPUUtilizations = append(Metrics.CPUUtilizations, metricValue)
}

func collectRAMUtilization() {
	if len(Metrics.RAMUtilizations) > maxCollectionValues {
		Metrics.RAMUtilizations = Metrics.RAMUtilizations[1:]
	}

	memoryUsage, _ := mem.VirtualMemory()
	metricValue := timestampedValue{time.Now(), int(memoryUsage.UsedPercent)}
	Metrics.RAMUtilizations = append(Metrics.RAMUtilizations, metricValue)
}
