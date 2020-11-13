/* eslint-disable react/prop-types */
import React, { useContext } from 'react';
import { Table, Typography } from 'antd';
import { ServerStatusContext } from '../utils/server-status-context';

const { Title } = Typography;


function VideoVariants({ config }) {
  if (!config || !config.videoSettings) {
    return null;
  }

  const videoQualityColumns = [
    {
      title: "Video bitrate",
      dataIndex: "videoBitrate",
      key: "videoBitrate",
      render: (bitrate) =>
        bitrate === 0 || !bitrate ? "Passthrough" : `${bitrate} kbps`,
    },
    {
      title: "Framerate",
      dataIndex: "framerate",
      key: "framerate",
    },
    {
      title: "Encoder preset",
      dataIndex: "encoderPreset",
      key: "framerate",
    },
    {
      title: "Audio bitrate",
      dataIndex: "audioBitrate",
      key: "audioBitrate",
      render: (bitrate) =>
        bitrate === 0 || !bitrate ? "Passthrough" : `${bitrate} kbps`,
    },
  ];

  const miscVideoSettingsColumns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
    },
    {
      title: "Value",
      dataIndex: "value",
      key: "value",
    },
  ];

  const miscVideoSettings = [
    {
      name: "Segment length",
      value: config.videoSettings.segmentLengthSeconds,
    },
    {
      name: "Number of segments",
      value: config.videoSettings.numberOfPlaylistItems,
    },
  ];

  return (
    <div>
      <Title>Video configuration</Title>
      <Table
        pagination={false}
        columns={videoQualityColumns}
        dataSource={config.videoSettings.videoQualityVariants}
      />

      <Table
        pagination={false}
        columns={miscVideoSettingsColumns}
        dataSource={miscVideoSettings}
      />
    </div>
  );
}

export default function VideoConfig() {  
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig: config } = serverStatusData || {};

  return (
    <div>
        <VideoVariants config={config} />
    </div>
  ); 
}

