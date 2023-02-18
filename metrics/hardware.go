package metrics

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	log "github.com/sirupsen/logrus"
)

// Max number of metrics we want to keep.
const maxCollectionValues = 300

func (s *Service) collectCPUUtilization() {
	if len(s.Metrics.CPUUtilizations) > maxCollectionValues {
		s.Metrics.CPUUtilizations = s.Metrics.CPUUtilizations[1:]
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

	metricValue := TimestampedValue{time.Now(), value}
	s.Metrics.CPUUtilizations = append(s.Metrics.CPUUtilizations, metricValue)
	cpuUsage.Set(metricValue.Value)
}

func (s *Service) collectRAMUtilization() {
	if len(s.Metrics.RAMUtilizations) > maxCollectionValues {
		s.Metrics.RAMUtilizations = s.Metrics.RAMUtilizations[1:]
	}

	memoryUsage, _ := mem.VirtualMemory()
	metricValue := TimestampedValue{time.Now(), memoryUsage.UsedPercent}
	s.Metrics.RAMUtilizations = append(s.Metrics.RAMUtilizations, metricValue)
}

func (s *Service) collectDiskUtilization() {
	path := "./"
	diskUse, _ := disk.Usage(path)

	if len(s.Metrics.DiskUtilizations) > maxCollectionValues {
		s.Metrics.DiskUtilizations = s.Metrics.DiskUtilizations[1:]
	}

	metricValue := TimestampedValue{time.Now(), diskUse.UsedPercent}
	s.Metrics.DiskUtilizations = append(s.Metrics.DiskUtilizations, metricValue)
}
