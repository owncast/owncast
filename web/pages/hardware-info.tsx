import React, { useState, useEffect } from 'react';
import { HARDWARE_STATS, fetchData, FETCH_INTERVAL } from './utils/apis';

export default function HardwareInfo() {
  const [hardwareStatus, setHardwareStatus] = useState({});

  const getHardwareStatus = async () => {
    try {
      const result = await fetchData(HARDWARE_STATS);
      console.log("hardare result", result)

      setHardwareStatus({ ...result });

    } catch (error) {
      setHardwareStatus({ ...hardwareStatus, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getHardwareStatus();
    getStatusIntervalId = setInterval(getHardwareStatus, FETCH_INTERVAL); //runs every 1 min.
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, []);

  return (
    <div>
      <h2>Hardware Info</h2>
      <p>cpu:[], disk: [], memory: [];  value = %age.</p>
      <p>the times should be the same for each, though milliseconds differ</p>
      <div style={{border: '1px solid blue', height: '300px', width: '100%', overflow:'auto'}}>
        {JSON.stringify(hardwareStatus)}
      </div>
    </div>
  );
}
