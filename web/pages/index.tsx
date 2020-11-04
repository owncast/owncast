/*
Will display an overview with the following datasources:
1. Current broadcaster.
2. Viewer count.
3. Video settings.

TODO: Link each overview value to the sub-page that focuses on it.
*/

import React, { useState, useEffect, useContext } from "react";
import { Row, Skeleton, Result, List, Typography, Card } from "antd";
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
import OwncastLogo from "./components/logo"

const { Title } = Typography;

function Offline() {
  const data = [
    {
      title: "Send some test content",
      content: (
        <div>
          With any video you have around you can pass it to the test script and start streaming it.
          <blockquote>
            <em>./test/ocTestStream.sh yourVideo.mp4</em>
          </blockquote>
        </div>
      ),
    },
    {
      title: "Use your broadcasting software",
      content: (
        <div>
          <a href="https://owncast.online/docs/broadcasting/">Learn how to point your existing software to your new server and start streaming your content.</a>
        </div>
      )
    },
    {
      title: "Something else",
    },
  ];
  return (
    <div>
      <Result
        icon={<OwncastLogo />}
        title="No stream is active."
        subTitle="You should start one."
      />

      <List
        grid={{
          gutter: 16,
          xs: 1,
          sm: 2,
          md: 4,
          lg: 4,
          xl: 6,
          xxl: 3,
        }}
        dataSource={data}
        renderItem={(item) => (
          <List.Item>
            <Card title={item.title}>{item.content}</Card>
          </List.Item>
        )}
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
