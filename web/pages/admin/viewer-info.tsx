import React, { useState, useEffect, useContext, ReactElement } from 'react';
import { Row, Col, Typography, MenuProps, Dropdown, Spin, Alert } from 'antd';
import { getUnixTime, sub } from 'date-fns';
import dynamic from 'next/dynamic';
import { Chart } from '../../components/admin/Chart';
import { StatisticItem } from '../../components/admin/StatisticItem';
import { ViewerTable } from '../../components/admin/ViewerTable';

import { ServerStatusContext } from '../../utils/server-status-context';

import { VIEWERS_OVER_TIME, ACTIVE_VIEWER_DETAILS, fetchData } from '../../utils/apis';

import { AdminLayout } from '../../components/layouts/AdminLayout';

// Lazy loaded components

const DownOutlined = dynamic(() => import('@ant-design/icons/DownOutlined'), {
  ssr: false,
});

const UserOutlined = dynamic(() => import('@ant-design/icons/UserOutlined'), {
  ssr: false,
});

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

  const offset: number = online && streamStart ? 0 : 1;
  const items: MenuProps['items'] = times.slice(offset).map((time, index) => ({
    key: index + offset,
    label: time.title,
    onClick: onTimeWindowSelect,
  }));

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
            title={online ? 'Max viewers this stream' : 'Max viewers last stream'}
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
      {!viewerInfo.length && (
        <Alert
          style={{ marginTop: '10px' }}
          banner
          message="Please wait"
          description="No viewer data has been collected yet."
          type="info"
        />
      )}

      <Spin spinning={!viewerInfo.length || loadingChart}>
        {viewerInfo.length > 0 && (
          <Chart
            title="Viewers"
            data={viewerInfo}
            color="#2087E2"
            unit="viewers"
            minYValue={0}
            yStepSize={1}
          />
        )}

        <Dropdown menu={{ items }} trigger={['click']}>
          <button
            type="button"
            style={{
              position: 'absolute',
              top: '5px',
              right: '35px',
              background: 'transparent',
              border: 'unset',
            }}
          >
            {timeWindowStart.title} <DownOutlined />
          </button>
        </Dropdown>
        <ViewerTable data={viewerDetails} />
      </Spin>
    </>
  );
}

ViewersOverTime.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
