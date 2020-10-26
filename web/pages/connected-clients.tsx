import React, { useState, useEffect, useContext } from 'react';
import { Table } from 'antd';
import { formatDistanceToNow } from "date-fns";
import { BroadcastStatusContext } from './utils/broadcast-status-context';

import { CONNECTED_CLIENTS, fetchData, FETCH_INTERVAL } from './utils/apis';

/*
geo data looks like this
  "geo": {
    "countryCode": "US",
    "regionName": "California",
    "timeZone": "America/Los_Angeles"
  }
*/

export default function ConnectedClients() {
  const context = useContext(BroadcastStatusContext);
  const { broadcastActive } = context || {};

  const [clients, setClients] = useState([]);
  const getInfo = async () => {
    try {
      const result = await fetchData(CONNECTED_CLIENTS);
      console.log("result",result)
      setClients(result);
    } catch (error) {
      console.log("==== error", error)
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();
    if (broadcastActive) {
      getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
      // returned function will be called on component unmount 
      return () => {
        clearInterval(getStatusIntervalId);
      }
    }
    return () => [];    
  }, []);
  
  if (!clients.length) {
    return "no clients";
  }

  // todo - check to see if broadcast active has changed. if so, start polling.

  const columns = [
    {
      title: "User name",
      dataIndex: "username",
      key: "username",
      render: (username) => username || "-",
      sorter: (a, b) => a.username - b.username,
      sortDirections: ["descend", "ascend"],
    },
    {
      title: "Messages sent",
      dataIndex: "messageCount",
      key: "messageCount",
      sorter: (a, b) => a.messageCount - b.messageCount,
      sortDirections: ["descend", "ascend"],
    },
    {
      title: "Connected Time",
      dataIndex: "connectedAt",
      key: "connectedAt",
      render: (time) => formatDistanceToNow(new Date(time)),
    },
    {
      title: "User Agent",
      dataIndex: "userAgent",
      key: "userAgent",
    },
    {
      title: "Location",
      dataIndex: "geo",
      key: "geo",
      render: (geo) => geo && `${geo.regionName}, ${geo.countryCode}`,
    },
  ];
  
  return (
    <div>
      <h2>Connected Clients</h2>
      <Table dataSource={clients} columns={columns} />;
    </div>
  );
}
