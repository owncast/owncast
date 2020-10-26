/*
Will display an overview with the following datasources:
1. Current broadcaster.
2. Viewer count.
3. Video settings.

TODO: Link each overview value to the sub-page that focuses on it.
*/

import React, { useState, useEffect, useContext } from "react";
import { Statistic, Card, Row, Col, Skeleton, Empty, Typography } from "antd";
import { UserOutlined, ClockCircleOutlined } from "@ant-design/icons";
import { formatDistanceToNow, formatRelative } from "date-fns";
import { BroadcastStatusContext } from "./utils/broadcast-status-context";
import {
  STREAM_STATUS,
  SERVER_CONFIG,
  fetchData,
  FETCH_INTERVAL,
} from "./utils/apis";
import { formatIPAddress } from "./utils/format";

const { Title} = Typography;

function Item(title: string, value: string, prefix: Jsx.Element) {
  const valueStyle = { color: "#334", fontSize: "1.8rem" };

  return (
    <Col span={8}>
      <Card>
        <Statistic
          title={title}
          value={value}
          valueStyle={valueStyle}
          prefix={prefix}
        />
      </Card>
    </Col>
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
  };


  // Pull in the server config so we can show config overview.
  const [videoSettings, setVideoSettings] = useState([]);
  const getConfig = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      setVideoSettings(result.videoSettings.videoQualityVariants);
    } catch (error) {
      console.log(error);
    }
  };

  useEffect(() => {
    setInterval(getStats, FETCH_INTERVAL);
    getStats();
    getConfig();
  }, []);

  if (!stats) {
    return (
      <div>
        <Skeleton active />
        <Skeleton active />
        <Skeleton active />
      </div>
    );
  }

  if (!broadcaster) {
    return Offline();
  }

  const videoQualitySettings = videoSettings.map(function (setting, index) {
    const audioSetting =
      setting.audioPassthrough || setting.audioBitrate === 0
        ? `${streamDetails.audioBitrate} kpbs (passthrough)`
        : `${setting.audioBitrate} kbps`;
    
    return (
      <Row gutter={[16, 16]} key={index}>
        {Item("Output", `Video variant ${index}`, "")}
        {Item(
          "Outbound Video Stream",
          `${setting.videoBitrate} kbps ${setting.framerate} fps`,
          ""
        )}
        {Item("Outbound Audio Stream", audioSetting, "")}
      </Row>
    );
  });

  const { viewerCount, sessionMaxViewerCount, lastConnectTime } = stats;
  const streamVideoDetailString = `${streamDetails.width}x${streamDetails.height} ${streamDetails.videoBitrate} kbps ${streamDetails.framerate} fps `;
  const streamAudioDetailString = `${streamDetails.audioCodec} ${streamDetails.audioBitrate} kpbs`;

  return (
    <div>
      <Title>Server Overview</Title>
      <Row gutter={[16, 16]}>
        {Item(
          `Stream started ${formatRelative(
            new Date(lastConnectTime),
            new Date()
          )}`,
          formatDistanceToNow(new Date(lastConnectTime)),
          <ClockCircleOutlined />
        )}

        {Item("Viewers", viewerCount, <UserOutlined />)}
        {Item("Peak viewer count", sessionMaxViewerCount, <UserOutlined />)}
      </Row>

      <Row gutter={[16, 16]}>
        {Item("Input", formatIPAddress(remoteAddr), "")}
        {Item("Inbound Video Stream", streamVideoDetailString, "")}
        {Item("Inbound Audio Stream", streamAudioDetailString, "")}
      </Row>

      {videoQualitySettings}
    </div>
  );
}

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
