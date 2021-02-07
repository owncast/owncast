/*
Will display an overview with the following datasources:
1. Current broadcaster.
2. Viewer count.
3. Video settings.

TODO: Link each overview value to the sub-page that focuses on it.
*/

import React, { useState, useEffect, useContext } from 'react';
import { Skeleton, Card, Statistic } from 'antd';
import { UserOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { formatDistanceToNow, formatRelative } from 'date-fns';
import { ServerStatusContext } from '../utils/server-status-context';
import StatisticItem from '../components/statistic';
import LogTable from '../components/log-table';
import Offline from './offline-notice';

import { LOGS_WARN, fetchData, FETCH_INTERVAL } from '../utils/apis';
import { formatIPAddress, isEmptyObject } from '../utils/format';
import { UpdateArgs } from '../types/config-section';

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
    return <Offline logs={logsData} />;
  }

  // map out settings
  const videoQualitySettings = serverStatusData?.currentBroadcast?.outputSettings?.map(
    (setting, index) => {
      const { audioPassthrough, videoPassthrough, audioBitrate, videoBitrate, framerate } = setting;

      const audioSetting = audioPassthrough
        ? `${streamDetails.audioCodec || 'Unknown'}, ${streamDetails.audioBitrate} kbps`
        : `${audioBitrate || 'Unknown'} kbps`;

      const videoSetting = videoPassthrough
        ? `${streamDetails.videoBitrate || 'Unknown'} kbps, ${streamDetails.framerate} fps ${
            streamDetails.width
          } x ${streamDetails.height}`
        : `${videoBitrate || 'Unknown'} kbps, ${framerate} fps`;

      let settingTitle = 'Outbound Stream Details';
      settingTitle =
        videoQualitySettings?.length > 1 ? `${settingTitle} ${index + 1}` : settingTitle;
      return (
        <Card title={settingTitle} type="inner" key={`${settingTitle}${index}`}>
          <StatisticItem title="Outbound Video Stream" value={videoSetting} prefix={null} />
          <StatisticItem title="Outbound Audio Stream" value={audioSetting} prefix={null} />
        </Card>
      );
    },
  );

  // inbound
  const { viewerCount, sessionPeakViewerCount } = serverStatusData;

  const streamAudioDetailString = `${streamDetails.audioCodec}, ${
    streamDetails.audioBitrate || 'Unknown'
  } kbps`;

  const broadcastDate = new Date(broadcaster.time);

  return (
    <div className="home-container">
      <div className="sections-container">
        <div className="section online-status-section">
          <Card title="Stream is online" type="inner">
            <Statistic
              title={`Stream started ${formatRelative(broadcastDate, Date.now())}`}
              value={formatDistanceToNow(broadcastDate)}
              prefix={<ClockCircleOutlined />}
            />
            <Statistic title="Viewers" value={viewerCount} prefix={<UserOutlined />} />
            <Statistic
              title="Peak viewer count"
              value={sessionPeakViewerCount}
              prefix={<UserOutlined />}
            />
          </Card>
        </div>

        <div className="section stream-details-section">
          <div className="details outbound-details">{videoQualitySettings}</div>

          <div className="details other-details">
            <Card title="Inbound Stream Details" type="inner">
              <StatisticItem
                title="Input"
                value={`${encoder} ${formatIPAddress(remoteAddr)}`}
                prefix={null}
              />
              <StatisticItem
                title="Inbound Video Stream"
                value={streamDetails}
                formatter={streamDetailsFormatter}
                prefix={null}
              />
              <StatisticItem
                title="Inbound Audio Stream"
                value={streamAudioDetailString}
                prefix={null}
              />
            </Card>

            <div className="server-detail">
              <Card title="Server Config" type="inner">
                <StatisticItem
                  title="Directory registration enabled"
                  value={configData.yp.enabled.toString()}
                  prefix={null}
                />
              </Card>
            </div>
          </div>
        </div>
      </div>
      <LogTable logs={logsData} pageSize={5} />
    </div>
  );
}
