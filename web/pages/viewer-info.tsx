import React, { useState, useEffect } from 'react';
import {timeFormat} from 'd3-time-format';
import { LineChart, XAxis, YAxis, Line, Tooltip } from 'recharts';

import { VIEWERS_OVER_TIME, fetchData } from './utils/apis';

const FETCH_INTERVAL = 5 * 60 * 1000; // 5 mins

export default function ViewersOverTime() {
  const [viewerInfo, setViewerInfo] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(VIEWERS_OVER_TIME);
      setViewerInfo(result);
    } catch (error) {
      console.log("==== error", error)
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

  const timeFormatter = (tick) => {return timeFormat('%H:%M:%S')(new Date(tick));};

  const CustomizedTooltip = (props) => {
    const { active, payload, label } = props;
    if (active) {
      const numViewers = payload && payload[0] && payload[0].value;
      const time = timeFormatter(label);
      const message = `${numViewers} viewer(s) at ${time}`;
      return (
        <div className="custom-tooltip">
          <p className="label">{message}</p>
        </div>
      );
    }
    return null;
  };


  return (
    <div>
      <h2>Current Viewers</h2>
      <div className="chart-container">
        <LineChart width={800} height={400} data={viewerInfo}>
          <XAxis dataKey="time" tickFormatter={timeFormatter}/>
          <YAxis dataKey="value"/>
          <Tooltip
            content={<CustomizedTooltip />}
          />
          <Line
            type="monotone"
            dataKey="value"
            stroke="#ff84d8" 
            dot={{ stroke: 'red', strokeWidth: 2 }} 
          />
        </LineChart>
      </div>
    </div>
  );
}
