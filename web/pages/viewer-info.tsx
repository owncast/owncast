import React, { useState, useEffect, useContext } from 'react';
import { Row, Col, Typography, Menu, Dropdown, Spin, Alert } from 'antd';
import { DownOutlined, UserOutlined } from '@ant-design/icons';
import { getUnixTime, sub } from 'date-fns';
import Chart from '../components/chart';
import StatisticItem from '../components/statistic';
import ViewerTable from '../components/viewer-table';

import { ServerStatusContext } from '../utils/server-status-context';

import { VIEWERS_OVER_TIME, ACTIVE_VIEWER_DETAILS, fetchData } from '../utils/apis';

const FETCH_INTERVAL = 60 * 1000; // 1 min

export default function ViewersOverTime() {
  const context = useContext(ServerStatusContext);
  const { online, broadcaster, viewerCount, overallPeakViewerCount, sessionPeakViewerCount } =
    context || {};
  let streamStart;
  if (broadcaster && broadcaster.time) {
    streamStart = new Date(broadcaster.time);
  }

  const times = [
    { title: 'Current stream', start: streamStart },
    { title: 'Last 12 hours', start: sub(new Date(), { hours: 12 }) },
    { title: 'Last 24 hours', start: sub(new Date(), { hours: 24 }) },
    { title: 'Last 7 days', start: sub(new Date(), { days: 7 }) },
    { title: 'Last 30 days', start: sub(new Date(), { days: 30 }) },
    { title: 'Last 3 months', start: sub(new Date(), { months: 3 }) },
    { title: 'Last 6 months', start: sub(new Date(), { months: 6 }) },
  ];

  const [loadingChart, setLoadingChart] = useState(true);
  const [viewerInfo, setViewerInfo] = useState([]);
  const [viewerDetails, setViewerDetails] = useState([]);
  const [timeWindowStart, setTimeWindowStart] = useState(times[1]);

  const getInfo = async () => {
    try {
      const url = `${VIEWERS_OVER_TIME}?windowStart=${getUnixTime(timeWindowStart.start)}`;
      const result = await fetchData(url);
      setViewerInfo(result);
      setLoadingChart(false);
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
  }, [online, timeWindowStart]);

  const onTimeWindowSelect = ({ key }) => {
    setTimeWindowStart(times[key]);
  };

  const menu = (
    <Menu>
      {online && streamStart && (
        <Menu.Item key="0" onClick={onTimeWindowSelect}>
          {times[0].title}
        </Menu.Item>
      )}
      {times.slice(1).map((time, index) => (
        // The array is hard coded, so it's safe to use the index as a key.
        // eslint-disable-next-line react/no-array-index-key
        <Menu.Item key={index + 1} onClick={onTimeWindowSelect}>
          {time.title}
        </Menu.Item>
      ))}
    </Menu>
  );

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
      <Dropdown overlay={menu} trigger={['click']}>
        <button
          type="button"
          style={{ float: 'right', background: 'transparent', border: 'unset' }}
        >
          {timeWindowStart.title} <DownOutlined />
        </button>
      </Dropdown>
      {viewerInfo.length > 0 && (
        <Chart title="Viewers" data={viewerInfo} color="#2087E2" unit="viewers" />
      )}

      <ViewerTable data={viewerDetails} />

      <Spin spinning={loadingChart}>
        {viewerInfo.length === 0 && (
          <Alert
            banner
            message="Please wait"
            description="No viewer data has been collected yet."
            type="info"
          />
        )}
      </Spin>
    </>
  );
}
