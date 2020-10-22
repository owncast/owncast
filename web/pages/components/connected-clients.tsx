import React, { useState, useEffect } from 'react';
import { Table } from 'antd';

import { CONNECTED_CLIENTS, fetchData, FETCH_INTERVAL } from '../utils/apis';

/*
geo data looks like this
  "geo": {
    "countryCode": "US",
    "regionName": "California",
    "timeZone": "America/Los_Angeles"
  }
*/

export default function HardwareInfo() {
  const [clients, setClients] = useState([]);
  const getInfo = async () => {
    try {
      const result = await fetchData(CONNECTED_CLIENTS);
      setClients(result);

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
  
  const columns = [
    {
      title: 'User name',
      dataIndex: 'username',
      key: 'username',
      render: username => username || '-',
      sorter: (a, b) => a.username - b.username,
      sortDirections: ['descend', 'ascend'],  
    },
    {
      title: 'Messages sent',
      dataIndex: 'messageCount',
      key: 'messageCount',
      sorter: (a, b) => a.messageCount - b.messageCount,
      sortDirections: ['descend', 'ascend'],  
    },
    {
      title: 'Connected Time',
      dataIndex: 'connectedAt',
      key: 'connectedAt',
      render: time => (Date.now() - (new Date(time).getTime())) / 1000 / 60,
    },
    {
      title: 'User Agent',
      dataIndex: 'userAgent',
      key: 'userAgent', 
    },
    {
      title: 'Location',
      dataIndex: 'geo',
      key: 'geo', 
      render: geo => geo && `${geo.regionName}, ${geo.countryCode}`,
    },
  ];
  
  console.log({clients})
  return (
    <div>
      <h2>Connected Clients</h2>
      <Table dataSource={clients} columns={columns} />;
    </div>
  );
}
