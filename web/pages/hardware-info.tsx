import { BulbOutlined, LaptopOutlined, SaveOutlined } from '@ant-design/icons';
import { Row } from 'antd';
import React, { useEffect, useState } from 'react';
import { fetchData, FETCH_INTERVAL, HARDWARE_STATS } from '../utils/apis';
import Chart from '../components/chart';
import StatisticItem from '../components/statistic';

interface TimedValue {
  time: Date;
  value: Number;
}

export default function HardwareInfo() {
  const [hardwareStatus, setHardwareStatus] = useState({
    cpu: Array<TimedValue>(),
    memory: Array<TimedValue>(),
    disk: Array<TimedValue>(),
    message: '',
  });

  const getHardwareStatus = async () => {
    try {
      const result = await fetchData(HARDWARE_STATS);
      setHardwareStatus({ ...result });
    } catch (error) {
      setHardwareStatus({ ...hardwareStatus, message: error.message });
    }
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    getHardwareStatus();
    getStatusIntervalId = setInterval(getHardwareStatus, FETCH_INTERVAL); // runs every 1 min.

    // returned function will be called on component unmount
    return () => {
      clearInterval(getStatusIntervalId);
    };
  }, []);

  if (!hardwareStatus.cpu) {
    return null;
  }

  const currentCPUUsage = hardwareStatus.cpu[hardwareStatus.cpu.length - 1]?.value;
  const currentRamUsage = hardwareStatus.memory[hardwareStatus.memory.length - 1]?.value;
  const currentDiskUsage = hardwareStatus.disk[hardwareStatus.disk.length - 1]?.value;

  const series = [
    {
      name: 'CPU',
      color: '#B63FFF',
      data: hardwareStatus.cpu,
    },
    {
      name: 'Memory',
      color: '#2087E2',
      data: hardwareStatus.memory,
    },
    {
      name: 'Disk',
      color: '#FF7700',
      data: hardwareStatus.disk,
    },
  ];

  return (
    <div>
      <div>
        <Row gutter={[16, 16]} justify="space-around">
          <StatisticItem
            title={series[0].name}
            value={`${currentCPUUsage}`}
            prefix={<LaptopOutlined style={{ color: series[0].color }} />}
            color={series[0].color}
            progress
            centered
          />
          <StatisticItem
            title={series[1].name}
            value={`${currentRamUsage}`}
            prefix={<BulbOutlined style={{ color: series[1].color }} />}
            color={series[1].color}
            progress
            centered
          />
          <StatisticItem
            title={series[2].name}
            value={`${currentDiskUsage}`}
            prefix={<SaveOutlined style={{ color: series[2].color }} />}
            color={series[2].color}
            progress
            centered
          />
        </Row>

        <Chart title="% used" dataCollections={series} color="#FF7700" unit="%" />
      </div>
    </div>
  );
}
