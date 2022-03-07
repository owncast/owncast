import React, { useState, useEffect, useContext } from 'react';
import { Row, Col, Typography } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import Chart from '../components/chart';
import StatisticItem from '../components/statistic';
import ViewerTable from '../components/viewer-table';

import { ServerStatusContext } from '../utils/server-status-context';

import { VIEWERS_OVER_TIME, ACTIVE_VIEWER_DETAILS, fetchData } from '../utils/apis';

const FETCH_INTERVAL = 60 * 1000; // 1 min

export default function ViewersOverTime() {
  const context = useContext(ServerStatusContext);
  const { online, viewerCount, overallPeakViewerCount, sessionPeakViewerCount } = context || {};

  const [viewerInfo, setViewerInfo] = useState([]);
  const [viewerDetails, setViewerDetails] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(VIEWERS_OVER_TIME);
      setViewerInfo(result);
    } catch (error) {
      console.log('==== error', error);
    }

    try {
      const result = await fetchData(ACTIVE_VIEWER_DETAILS);
      setViewerDetails(result);
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

  if (!viewerInfo.length) {
    return 'no info';
  }

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
      <ViewerTable data={viewerDetails} />
    </>
  );
}
