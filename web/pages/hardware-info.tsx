/* eslint-disable no-array-constructor */
import React, { useState, useEffect } from 'react';
import { Row } from "antd";
import {LaptopOutlined, BulbOutlined, SaveOutlined} from "@ant-design/icons"
import { HARDWARE_STATS, fetchData, FETCH_INTERVAL } from '../utils/apis';
import Chart from './components/chart';
import StatisticItem from "./components/statistic";

interface TimedValue {
  time: Date,
  value: Number
}

export default function HardwareInfo() {
  const [hardwareStatus, setHardwareStatus] = useState({
    cpu: Array<TimedValue>(),
    memory: Array<TimedValue>(),
    disk: Array<TimedValue>(),
    message: "",
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
    }
  }, []);

 
  if (!hardwareStatus.cpu) {
    return null;
  }

  const currentCPUUsage = hardwareStatus.cpu[hardwareStatus.cpu.length - 1]?.value;
  const currentRamUsage =
    hardwareStatus.memory[hardwareStatus.memory.length - 1]?.value;
  const currentDiskUsage =
    hardwareStatus.disk[hardwareStatus.disk.length - 1]?.value;
  
const series = [
  {
    name: "CPU",
    color: "#FF7700",
    data: hardwareStatus.cpu,
  },
  {
    name: "Memory",
    color: "#004777",
    data: hardwareStatus.memory,
  },
  {
    name: "Disk",
    color: "#A9E190",
    data: hardwareStatus.disk,
  },
];
  
    return (
      <div>
        <div>
          <h2>Hardware Info</h2>
          <Row gutter={[16, 16]}>
            <StatisticItem
              title="CPU used"
              value={`${currentCPUUsage} %`}
              prefix={<LaptopOutlined />}
            />
            <StatisticItem
              title="Memory used"
              value={`${currentRamUsage} %`}
              prefix={<BulbOutlined />}
            />
            <StatisticItem
              title="Disk used"
              value={`${currentDiskUsage} %`}
              prefix={<SaveOutlined />}
            />
          </Row>

          <div className="chart-container">
            <Chart dataCollections={series} color="#FF7700" unit="%" />
          </div>
        </div>
        <p>cpu:[], disk: [], memory: []; value = %age.</p>
        <p>the times should be the same for each, though milliseconds differ</p>
        <div
          style={{
            border: "1px solid blue",
            height: "300px",
            width: "100%",
            overflow: "auto",
          }}
        >
          {JSON.stringify(hardwareStatus)}
        </div>
      </div>
    );
}