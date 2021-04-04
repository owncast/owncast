import React, { useState, useEffect, useContext } from 'react';
import { Skeleton, Card, Statistic, Row, Col } from 'antd';
import { UserOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { formatDistanceToNow, formatRelative } from 'date-fns';
import { ServerStatusContext } from '../utils/server-status-context';
import LogTable from '../components/log-table';
import Offline from './offline-notice';

import { LOGS_WARN, fetchData, FETCH_INTERVAL } from '../utils/apis';
import { formatIPAddress, isEmptyObject } from '../utils/format';
import NewsFeed from '../components/news-feed';

function streamDetailsFormatter(streamDetails) {
  return (
    <ul className="statistics-list">
      <li>
        {streamDetails.videoCodec || 'Unknown'} @ {streamDetails.videoBitrate || 'Unknown'} kbps
      </li>
      <li>{streamDetails.framerate || 'Unknown'} fps</li>
      <li>
        {streamDetails.width} x {streamDetails.height}
      </li>
    </ul>
  );
}

export default function Home() {
  const serverStatusData = useContext(ServerStatusContext);
  const { broadcaster, serverConfig: configData } = serverStatusData || {};
  const { remoteAddr, streamDetails } = broadcaster || {};

  const encoder = streamDetails?.encoder || 'Unknown encoder';

  const [logsData, setLogs] = useState([]);
  const getLogs = async () => {
    try {
      const result = await fetchData(LOGS_WARN);
      setLogs(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };
  const getMoreStats = () => {
    getLogs();
  };

  useEffect(() => {
    getMoreStats();

    let intervalId = null;
    intervalId = setInterval(getMoreStats, FETCH_INTERVAL);

    return () => {
      clearInterval(intervalId);
    };
  }, []);

  if (isEmptyObject(configData) || isEmptyObject(serverStatusData)) {
    return (
      <>
        <Skeleton active />
        <Skeleton active />
        <Skeleton active />
      </>
    );
  }

  if (!broadcaster) {
    return <Offline logs={logsData} config={configData} />;
  }

  // map out settings
  const videoQualitySettings = serverStatusData?.currentBroadcast?.outputSettings?.map(setting => {
    const { audioPassthrough, videoPassthrough, audioBitrate, videoBitrate, framerate } = setting;

    const audioSetting = audioPassthrough
      ? `${streamDetails.audioCodec || 'Unknown'}, ${streamDetails.audioBitrate} kbps`
      : `${audioBitrate || 'Unknown'} kbps`;

    const videoSetting = videoPassthrough
      ? `${streamDetails.videoBitrate || 'Unknown'} kbps, ${streamDetails.framerate} fps ${
          streamDetails.width
        } x ${streamDetails.height}`
      : `${videoBitrate || 'Unknown'} kbps, ${framerate} fps`;

    return (
      <div className="stream-details-item-container">
        <Statistic
          className="stream-details-item"
          title="Outbound Video Stream"
          value={videoSetting}
        />
        <Statistic
          className="stream-details-item"
          title="Outbound Audio Stream"
          value={audioSetting}
        />
      </div>
    );
  });

  // inbound
  const { viewerCount, sessionPeakViewerCount } = serverStatusData;

  const streamAudioDetailString = `${streamDetails.audioCodec}, ${
    streamDetails.audioBitrate || 'Unknown'
  } kbps`;

  const broadcastDate = new Date(broadcaster.time);

  return (
    <div className="home-container">
      <div className="sections-container">
        <div className="online-status-section">
          <Card size="small" type="inner" className="online-details-card">
            <Row gutter={[16, 16]} align="middle">
              <Col span={8} sm={24} md={8}>
                <Statistic
                  title={`Stream started ${formatRelative(broadcastDate, Date.now())}`}
                  value={formatDistanceToNow(broadcastDate)}
                  prefix={<ClockCircleOutlined />}
                />
              </Col>
              <Col span={8} sm={24} md={8}>
                <Statistic title="Viewers" value={viewerCount} prefix={<UserOutlined />} />
              </Col>
              <Col span={8} sm={24} md={8}>
                <Statistic
                  title="Peak viewer count"
                  value={sessionPeakViewerCount}
                  prefix={<UserOutlined />}
                />
              </Col>
            </Row>
          </Card>
        </div>

        <Row gutter={[16, 16]} className="section stream-details-section">
          <Col className="stream-details" span={12} sm={24} md={24} lg={12}>
            <Card
              size="small"
              title="Outbound Stream Details"
              type="inner"
              className="outbound-details"
            >
              {videoQualitySettings}
            </Card>

            <Card size="small" title="Inbound Stream Details" type="inner">
              <Statistic
                className="stream-details-item"
                title="Input"
                value={`${encoder} ${formatIPAddress(remoteAddr)}`}
              />
              <Statistic
                className="stream-details-item"
                title="Inbound Video Stream"
                value={streamDetails}
                formatter={streamDetailsFormatter}
              />
              <Statistic
                className="stream-details-item"
                title="Inbound Audio Stream"
                value={streamAudioDetailString}
              />
            </Card>
          </Col>
          <Col span={12} xs={24} sm={24} md={24} lg={12}>
            <NewsFeed />
          </Col>
        </Row>
      </div>
      <br />
      <LogTable logs={logsData} pageSize={5} />
    </div>
  );
}
