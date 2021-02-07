import React, { useState, useEffect } from 'react';
import LogTable from '../components/log-table';

import { LOGS_ALL, fetchData } from '../utils/apis';

const FETCH_INTERVAL = 5 * 1000; // 5 sec

export default function Logs() {
  const [logs, setLogs] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(LOGS_ALL);
      setLogs(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    setInterval(getInfo, FETCH_INTERVAL);
    getInfo();

    getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
    // returned function will be called on component unmount
    return () => {
      clearInterval(getStatusIntervalId);
    };
  }, []);

  return <LogTable logs={logs} pageSize={20} />;
}
