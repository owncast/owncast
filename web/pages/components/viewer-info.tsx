import React, { useState, useEffect } from 'react';
import { VIEWERS_OVER_TIME, fetchData, FETCH_INTERVAL } from '../utils/apis';

export default function HardwareInfo() {
  const [viewerInfo, setViewerInfo] = useState({});

  const getInfo = async () => {
    try {
      const result = await fetchData(VIEWERS_OVER_TIME);
      console.log("viewers result", result)

      setViewerInfo({ ...result });

    } catch (error) {
      setViewerInfo({ ...viewerInfo, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();
    getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, []);

  return (
    <div>
      <h2>Viewers over time</h2>
      <div style={{border: '1px solid red', height: '300px', width: '100%', overflow:'auto'}}>
        {JSON.stringify(viewerInfo)}
      </div>
    </div>
  );
}
