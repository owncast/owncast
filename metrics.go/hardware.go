package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func collectCPUUtilization() {
	value, err := cpu.Percent(0, false)
	if err != nil {
		panic(err)
	}

	Metrics.CPUUtilizations = append(Metrics.CPUUtilizations, int(value[0]))
}

func collectRAMUtilization() {
	value, _ := mem.VirtualMemory()
	Metrics.RAMUtilizations = append(Metrics.RAMUtilizations, int(value.UsedPercent))
}
