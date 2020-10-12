import React, { useState, useEffect } from 'react';
import { VIEWERS_OVER_TIME, fetchData } from '../utils/apis';


const FETCH_INTERVAL = 5 * 60 * 1000; // 5 mins

export default function ViewersOverTime() {
  const [viewerInfo, setViewerInfo] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(VIEWERS_OVER_TIME);
      console.log("viewers result", result)

      setViewerInfo(result);

    } catch (error) {
      console.log("==== error", error)
      // setViewerInfo({ ...viewerInfo, message: error.message });
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


  const formattedData = viewerInfo.map(viewer => ({
    x: (new Date(viewer.time)).toLocaleTimeString(),
    y: viewer.value,
  }));

  return (
    <div>
      <h2>Viewers over time</h2>
      <p>Time on X axis, # Viewer on Y</p>
      <div style={{border: '1px solid red', height: '300px', width: '100%', overflow:'auto'}}>
        {JSON.stringify(formattedData)}
      </div>
    </div>
  );
}
