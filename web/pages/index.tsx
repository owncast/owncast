/*
Will display an overview with the following datasources:
1. Current broadcaster.
2. Viewer count.
3. Video settings.

TODO: Link each overview value to the sub-page that focuses on it.
*/

import React, { useState, useEffect, useContext } from "react";
import { Row, Skeleton, Empty, Typography } from "antd";
import { UserOutlined, ClockCircleOutlined } from "@ant-design/icons";
import { formatDistanceToNow, formatRelative } from "date-fns";
import { BroadcastStatusContext } from "../utils/broadcast-status-context";
import StatisticItem from "./components/statistic"
import LogTable from "./components/log-table";

import {
  STREAM_STATUS,
  SERVER_CONFIG,
  LOGS_WARN,
  fetchData,
  FETCH_INTERVAL,
} from "../utils/apis";
import { formatIPAddress, isEmptyObject } from "../utils/format";

const { Title } = Typography;

function Offline() {
  return (
    <div>
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description={
          <span>There is no stream currently active. Start one.</span>
        }
      />
    </div>
  );
}

export default function Stats() {
  const context = useContext(BroadcastStatusContext);
  const { broadcaster } = context || {};
  const { remoteAddr, streamDetails } = broadcaster || {};

  // Pull in the server status so we can show server overview.
  const [stats, setStats] = useState(null);
  const getStats = async () => {
    try {
      const result = await fetchData(STREAM_STATUS);
      setStats(result);
    } catch (error) {
      console.log(error);
    }

    getConfig();
    getLogs();
  };


  // Pull in the server config so we can show config overview.
  const [config, setConfig] = useState({
    streamKey: "",
    yp: {
      enabled: false,
    },
    videoSettings: {
      videoQualityVariants: [
        {
          audioPassthrough: false,
          videoBitrate: 0,
          audioBitrate: 0,
          framerate: 0,
        },
      ],
    },
  });
  const [logs, setLogs] = useState([]);

  const getConfig = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      setConfig(result);
    } catch (error) {
      console.log(error);
    }
  };

  const getLogs = async () => {
    try {
      const result = await fetchData(LOGS_WARN);
      setLogs(result);
    } catch (error) {
      console.log("==== error", error);
    }
  };

  useEffect(() => {
    setInterval(getStats, FETCH_INTERVAL);
    getStats();
  }, []);

  if (isEmptyObject(config) || isEmptyObject(stats)) {
    return (
      <div>
        <Skeleton active />
        <Skeleton active />
        <Skeleton active />
      </div>
    );
  }

  if (!broadcaster) {
    return <Offline />;
  }
  
  const videoSettings = config.videoSettings.videoQualityVariants;
  const videoQualitySettings = videoSettings.map((setting, index) => {
    const audioSetting =
      setting.audioPassthrough || setting.audioBitrate === 0
        ? `${streamDetails.audioBitrate} kpbs (passthrough)`
        : `${setting.audioBitrate} kbps`;
    
    return (
      // eslint-disable-next-line react/no-array-index-key
      <Row gutter={[16, 16]} key={index}>
        <StatisticItem
          title="Outbound Video Stream"
          value={`${setting.videoBitrate} kbps ${setting.framerate} fps`}
          prefix={null}
        />
        <StatisticItem
          title="Outbound Audio Stream"
          value={audioSetting}
          prefix={null}
        />
      </Row>
    );
  });

  const logTable = logs.length > 0 ? <LogTable logs={logs} pageSize={5} /> : null
  const { viewerCount, sessionMaxViewerCount, lastConnectTime } = stats;
  const streamVideoDetailString = `${streamDetails.videoCodec} ${streamDetails.videoBitrate} kbps ${streamDetails.width}x${streamDetails.height}`;
  const streamAudioDetailString = `${streamDetails.audioCodec} ${streamDetails.audioBitrate} kpbs`;

  return (
    <div>
      <Title>Server Overview</Title>
      <Row gutter={[16, 16]}>
        <StatisticItem
          title={`Stream started ${formatRelative(
            new Date(lastConnectTime),
            new Date()
          )}`}
          value={formatDistanceToNow(new Date(lastConnectTime))}
          prefix={<ClockCircleOutlined />}
        />
        <StatisticItem
          title="Viewers"
          value={viewerCount}
          prefix={<UserOutlined />}
        />
        <StatisticItem
          title="Peak viewer count"
          value={sessionMaxViewerCount}
          prefix={<UserOutlined />}
        />
      </Row>

      <Row gutter={[16, 16]}>
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
      </Row>

      {videoQualitySettings}

      <Row gutter={[16, 16]}>
        <StatisticItem
          title="Stream key"
          value={config.streamKey}
          prefix={null}
        />
        <StatisticItem
          title="Directory registration enabled"
          value={config.yp.enabled.toString()}
          prefix={null}
        />
      </Row>

      {logTable}
    </div>
  );
}
