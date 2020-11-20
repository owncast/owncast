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
      title: "#",
      dataIndex: "key",
      key: "key"
    },
    {
      title: "Video bitrate",
      dataIndex: "videoBitrate",
      key: "videoBitrate",
      render: (bitrate) =>
        !bitrate ? "Same as source" : `${bitrate} kbps`,
    },
    {
      title: "Framerate",
      dataIndex: "framerate",
      key: "framerate",
      render: (framerate) =>
        !framerate ? "Same as source" : `${framerate} fps`,
    },
    {
      title: "Encoder preset",
      dataIndex: "encoderPreset",
      key: "framerate",
      render: (preset) =>
        !preset ? "n/a" : preset,
    },
    {
      title: "Audio bitrate",
      dataIndex: "audioBitrate",
      key: "audioBitrate",
      render: (bitrate) =>
        !bitrate ? "Same as source" : `${bitrate} kbps`,
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
      key: "segmentLength"
    },
    {
      name: "Number of segments",
      value: config.videoSettings.numberOfPlaylistItems,
      key: "numberOfSegments"
    },
  ];

  const videoQualityVariantData = config.videoSettings.videoQualityVariants.map(function(variant, index) {
    return {
      key: index,
      ...variant
    }
  });

  return (
    <div>
      <Title>Video configuration</Title>
      <Table
        pagination={false}
        columns={videoQualityColumns}
        dataSource={videoQualityVariantData}
      />

      <Table
        pagination={false}
        columns={miscVideoSettingsColumns}
        dataSource={miscVideoSettings}
        rowKey={row => row.name}
      />
      <br/>
      <Title level={5}>Learn more about configuring Owncast <a href="https://owncast.online/docs/configuration">by visiting the documentation.</a></Title>
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

