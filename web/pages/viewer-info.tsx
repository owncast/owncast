import React, { useState, useEffect, useContext } from 'react';
import { Table, Row, Col, Typography } from 'antd';
import { formatDistanceToNow } from 'date-fns';
import { UserOutlined } from '@ant-design/icons';
import { SortOrder } from 'antd/lib/table/interface';
import Chart from '../components/chart';
import StatisticItem from '../components/statistic';

import { ServerStatusContext } from '../utils/server-status-context';

import { CONNECTED_CLIENTS, VIEWERS_OVER_TIME, fetchData } from '../utils/apis';

const FETCH_INTERVAL = 60 * 1000; // 1 min

export default function ViewersOverTime() {
  const context = useContext(ServerStatusContext);
  const { online, viewerCount, overallPeakViewerCount, sessionPeakViewerCount } = context || {};

  const [viewerInfo, setViewerInfo] = useState([]);
  const [clients, setClients] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(VIEWERS_OVER_TIME);
      setViewerInfo(result);
    } catch (error) {
      console.log('==== error', error);
    }

    try {
      const result = await fetchData(CONNECTED_CLIENTS);
      setClients(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();
    if (online) {
      getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
      // returned function will be called on component unmount
      return () => {
        clearInterval(getStatusIntervalId);
      };
    }

    return () => [];
  }, [online]);

  // todo - check to see if broadcast active has changed. if so, start polling.

  if (!viewerInfo.length) {
    return 'no info';
  }

  const columns = [
    {
      title: 'User name',
      dataIndex: 'username',
      key: 'username',
      render: username => username || '-',
      sorter: (a, b) => a.username - b.username,
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'Messages sent',
      dataIndex: 'messageCount',
      key: 'messageCount',
      sorter: (a, b) => a.messageCount - b.messageCount,
      sortDirections: ['descend', 'ascend'] as SortOrder[],
    },
    {
      title: 'Connected Time',
      dataIndex: 'connectedAt',
      key: 'connectedAt',
      render: time => formatDistanceToNow(new Date(time)),
      sorter: (a, b) => new Date(a.connectedAt).getTime() - new Date(b.connectedAt).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
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
      render: geo => (geo ? `${geo.regionName}, ${geo.countryCode}` : '-'),
    },
  ];

  return (
    <>
      <Typography.Title>Viewer Info</Typography.Title>
      <br />
      <Row gutter={[16, 16]} justify="space-around">
        {online && (
          <Col span={8} md={8}>
            <StatisticItem
              title="Current viewers"
              value={viewerCount.toString()}
              prefix={<UserOutlined />}
            />
          </Col>
        )}
        <Col md={online ? 8 : 12}>
          <StatisticItem
            title={online ? 'Max viewers this session' : 'Max viewers last session'}
            value={sessionPeakViewerCount.toString()}
            prefix={<UserOutlined />}
          />
        </Col>
        <Col md={online ? 8 : 12}>
          <StatisticItem
            title="All-time max viewers"
            value={overallPeakViewerCount.toString()}
            prefix={<UserOutlined />}
          />
        </Col>
      </Row>

      <Chart title="Viewers" data={viewerInfo} color="#2087E2" unit="" />
      {online && (
        <div>
          <Table dataSource={clients} columns={columns} rowKey={row => row.clientID} />
          <p>
            <Typography.Text type="secondary">
              Visit the{' '}
              <a href="https://owncast.online/docs/viewers/?source=admin">documentation</a> to
              configure additional details about your viewers.
            </Typography.Text>{' '}
          </p>
        </div>
      )}
    </>
  );
}
