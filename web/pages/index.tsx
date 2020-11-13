/* eslint-disable no-console */
/*
Will display an overview with the following datasources:
1. Current broadcaster.
2. Viewer count.
3. Video settings.

TODO: Link each overview value to the sub-page that focuses on it.
*/

import React, { useState, useEffect, useContext } from "react";
import { Row, Col, Skeleton, Result, List, Typography, Card, Statistic } from "antd";
import { UserOutlined, ClockCircleOutlined } from "@ant-design/icons";
import { formatDistanceToNow, formatRelative } from "date-fns";
import { ServerStatusContext } from "../utils/server-status-context";
import StatisticItem from "./components/statistic"
import LogTable from "./components/log-table";
import Offline from './offline-notice';

import {
  SERVER_CONFIG,
  LOGS_WARN,
  fetchData,
  FETCH_INTERVAL,
} from "../utils/apis";
import { formatIPAddress, isEmptyObject } from "../utils/format";
import { INITIAL_SERVER_CONFIG_STATE } from "./update-server-config";

const { Title } = Typography;


export default function Home() {
  const serverStatusData = useContext(ServerStatusContext);
  const { broadcaster } = serverStatusData || {};
  const { remoteAddr, streamDetails } = broadcaster || {};

  // Pull in the server config so we can show config overview.
  const [configData, setServerConfig] = useState(INITIAL_SERVER_CONFIG_STATE);
  const getConfig = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      setServerConfig(result);
      console.log("CONFIG", result);
    } catch (error) {
      console.log(error);
    }
  };

  const [logsData, setLogs] = useState([]);
  const getLogs = async () => {
    try {
      const result = await fetchData(LOGS_WARN);
      setLogs(result);
      console.log("LOGS", result);
    } catch (error) {
      console.log("==== error", error);
    }
  };
  const getMoreStats = () => {
    getLogs();
    getConfig();
  }
  useEffect(() => {
    let intervalId = null;
    intervalId = setInterval(getMoreStats, FETCH_INTERVAL);
    return () => {
      clearInterval(intervalId);
    }
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
  const videoQualitySettings = configData?.videoSettings?.videoQualityVariants?.map((setting, index) => {
    const { audioPassthrough, audioBitrate, videoBitrate, framerate } = setting;
    const audioSetting =
      audioPassthrough || audioBitrate === 0
        ? `${streamDetails.audioBitrate} kpbs (passthrough)`
        : `${audioBitrate} kbps`;
    
    let settingTitle = 'Outbound Stream Details';
    settingTitle = (videoQualitySettings?.length > 1) ?
      `${settingTitle} ${index + 1}` : settingTitle;
    return (
      <Card title={settingTitle} type="inner">
        <StatisticItem
          title="Outbound Video Stream"
          value={`${videoBitrate} kbps, ${framerate} fps`}
          prefix={null}
        />
        <StatisticItem
          title="Outbound Audio Stream"
          value={audioSetting}
          prefix={null}
        />
      </Card>
    );
  });

  const { viewerCount, sessionMaxViewerCount } = serverStatusData;
  const streamVideoDetailString = `${streamDetails.videoCodec} ${streamDetails.videoBitrate} kbps ${streamDetails.width}x${streamDetails.height}`;
  const streamAudioDetailString = `${streamDetails.audioCodec} ${streamDetails.audioBitrate} kbps`;

  const broadcastDate = new Date(broadcaster.time);

  return (
    <div className="home-container">
      <Title>Stream Overview</Title>

      <div className="sections-container">


        <div className="section online-status-section">
          <Card title="Stream is online" type="inner">
            <Statistic
              title={`Stream started ${formatRelative(
                broadcastDate,
                Date.now()
              )}`}
              value={formatDistanceToNow(broadcastDate)}
              prefix={<ClockCircleOutlined />}
            />
            <Statistic
              title="Viewers"
              value={viewerCount}
              prefix={<UserOutlined />}
            />
            <Statistic
              title="Peak viewer count"
              value={sessionMaxViewerCount}
              prefix={<UserOutlined />}
            />
          </Card>
          </div>

        <div className="section stream-details-section">

          <div className="details outbound-details">
            {videoQualitySettings}
          </div>

          <div className="details other-details">
            <Card title="Inbound Stream Details" type="inner">
              <StatisticItem
                title="Input"
                value={formatIPAddress(remoteAddr)}
                prefix={null}
              />
              <StatisticItem
                title="Inbound Video Stream"
                value={streamVideoDetailString}
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
                  title="Stream key"
                  value={configData.streamKey}
                  prefix={null}
                />
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

      {logsData.length ? (
        <>
          <Title level={2}>Stream Logs</Title>
          <LogTable logs={logsData} pageSize={5} />
        </>
      ): null}
    </div>
  );
}
