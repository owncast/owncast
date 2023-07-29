import { Row, Col, Typography, Alert, Spin } from 'antd';
import React, { ReactElement, useEffect, useState } from 'react';
import dynamic from 'next/dynamic';
import { fetchData, FETCH_INTERVAL, HARDWARE_STATS } from '../../utils/apis';
import { Chart } from '../../components/admin/Chart';
import { StatisticItem } from '../../components/admin/StatisticItem';

import { AdminLayout } from '../../components/layouts/AdminLayout';

// Lazy loaded components

const BulbOutlined = dynamic(() => import('@ant-design/icons/BulbOutlined'), {
  ssr: false,
});

const LaptopOutlined = dynamic(() => import('@ant-design/icons/LaptopOutlined'), {
  ssr: false,
});

const SaveOutlined = dynamic(() => import('@ant-design/icons/SaveOutlined'), {
  ssr: false,
});

export default function HardwareInfo() {
  const [hardwareStatus, setHardwareStatus] = useState({
    cpu: [], // Array<TimedValue>(),
    memory: [], // Array<TimedValue>(),
    disk: [], // Array<TimedValue>(),
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
    return (
      <div>
        <Typography.Title>Hardware Info</Typography.Title>

        <Alert
          style={{ marginTop: '10px' }}
          banner
          message="Please wait"
          description="No hardware details have been collected yet."
          type="info"
        />
        <Spin spinning style={{ width: '100%', margin: '10px' }} />
      </div>
    );
  }

  const currentCPUUsage = hardwareStatus.cpu[hardwareStatus.cpu.length - 1]?.value;
  const currentRamUsage = hardwareStatus.memory[hardwareStatus.memory.length - 1]?.value;
  const currentDiskUsage = hardwareStatus.disk[hardwareStatus.disk.length - 1]?.value;

  const series = [
    {
      name: 'CPU',
      color: '#B63FFF',
      data: hardwareStatus.cpu,
      pointStyle: 'rect',
    },
    {
      name: 'Memory',
      color: '#2087E2',
      data: hardwareStatus.memory,
      pointStyle: 'circle',
    },
    {
      name: 'Disk',
      color: '#FF7700',
      data: hardwareStatus.disk,
      pointStyle: 'rectRounded',
    },
  ];

  return (
    <>
      <Typography.Title>Hardware Info</Typography.Title>
      <br />
      <div>
        <Row gutter={[16, 16]} justify="space-around">
          <Col>
            <StatisticItem
              title={series[0].name}
              value={Math.round(currentCPUUsage) || 0}
              prefix={<LaptopOutlined style={{ color: series[0].color }} />}
              color={series[0].color}
              progress
              centered
            />
          </Col>
          <Col>
            <StatisticItem
              title={series[1].name}
              value={Math.round(currentRamUsage) || 0}
              prefix={<BulbOutlined style={{ color: series[1].color }} />}
              color={series[1].color}
              progress
              centered
            />
          </Col>
          <Col>
            <StatisticItem
              title={series[2].name}
              value={Math.round(currentDiskUsage) || 0}
              prefix={<SaveOutlined style={{ color: series[2].color }} />}
              color={series[2].color}
              progress
              centered
            />
          </Col>
        </Row>

        <Chart title="% used" dataCollections={series} color="#FF7700" unit="%" />
      </div>
    </>
  );
}

HardwareInfo.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
