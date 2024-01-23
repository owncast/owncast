import React, { useState, useEffect, ReactElement } from 'react';
import { LogTable } from '../../components/admin/LogTable';

import { LOGS_ALL, fetchData } from '../../utils/apis';

import { AdminLayout } from '../../components/layouts/AdminLayout';

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

  return <LogTable logs={logs} initialPageSize={20} />;
}

Logs.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
