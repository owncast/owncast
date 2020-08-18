package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Max number of metrics we want to keep.
const maxCollectionValues = 500

func collectCPUUtilization() {
	if len(Metrics.CPUUtilizations) > maxCollectionValues {
		Metrics.CPUUtilizations = Metrics.CPUUtilizations[1:]
	}

	value, err := cpu.Percent(0, false)
	if err != nil {
		panic(err)
	}

	Metrics.CPUUtilizations = append(Metrics.CPUUtilizations, int(value[0]))
}

func collectRAMUtilization() {
	if len(Metrics.RAMUtilizations) > maxCollectionValues {
		Metrics.RAMUtilizations = Metrics.RAMUtilizations[1:]
	}

	value, _ := mem.VirtualMemory()
	Metrics.RAMUtilizations = append(Metrics.RAMUtilizations, int(value.UsedPercent))
}
