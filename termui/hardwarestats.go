package termui

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func getMemoryUsage() int {
	v, _ := mem.VirtualMemory()
	return int(v.UsedPercent)
}

func getCPUUsage() int {
	v, err := cpu.Percent(0, false)
	if err != nil {
		panic(err)
	}
	return int(v[0])
}
