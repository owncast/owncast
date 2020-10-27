import React, { useState, useEffect } from 'react';
import { timeFormat } from "d3-time-format";
import { HARDWARE_STATS, fetchData, FETCH_INTERVAL } from './utils/apis';
import Chart from './components/chart.tsx'

export default function HardwareInfo() {
  const [hardwareStatus, setHardwareStatus] = useState({});

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

    return (
      <div>
        <div>
          <h2>Hardware Info</h2>
          <div className="chart-container">
            <h3>CPU</h3>
            <Chart data={hardwareStatus.cpu} color="#FF7700" unit="%" />
            <h3>Memory</h3>
            <Chart data={hardwareStatus.memory} color="#004777" unit="%" />
            <h3>Disk</h3>
            <Chart data={hardwareStatus.disk} color="#A9E190" unit="%" />
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