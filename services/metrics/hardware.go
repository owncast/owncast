package metrics

import (
	"time"

	"github.com/owncast/owncast/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	log "github.com/sirupsen/logrus"
)

// Max number of metrics we want to keep.
const maxCollectionValues = 300

func (m *Metrics) collectCPUUtilization() {
	if len(m.metrics.CPUUtilizations) > maxCollectionValues {
		m.metrics.CPUUtilizations = m.metrics.CPUUtilizations[1:]
	}

	v, err := cpu.Percent(0, false)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Default to zero but try to use the cumulative values of all the CPUs
	// if values exist.
	value := 0.0
	if len(v) > 0 {
		value = v[0]
	}

	metricValue := models.TimestampedValue{time.Now(), value}
	m.metrics.CPUUtilizations = append(m.metrics.CPUUtilizations, metricValue)
	m.cpuUsage.Set(metricValue.Value)
}

func (m *Metrics) collectRAMUtilization() {
	if len(m.metrics.RAMUtilizations) > maxCollectionValues {
		m.metrics.RAMUtilizations = m.metrics.RAMUtilizations[1:]
	}

	memoryUsage, _ := mem.VirtualMemory()
	metricValue := models.TimestampedValue{time.Now(), memoryUsage.UsedPercent}
	m.metrics.RAMUtilizations = append(m.metrics.RAMUtilizations, metricValue)
}

func (m *Metrics) collectDiskUtilization() {
	path := "./"
	diskUse, _ := disk.Usage(path)

	if len(m.metrics.DiskUtilizations) > maxCollectionValues {
		m.metrics.DiskUtilizations = m.metrics.DiskUtilizations[1:]
	}

	metricValue := models.TimestampedValue{time.Now(), diskUse.UsedPercent}
	m.metrics.DiskUtilizations = append(m.metrics.DiskUtilizations, metricValue)
}
