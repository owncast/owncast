/* eslint-disable react/prop-types */
import React, { useState, useEffect } from 'react';
import { Table, Typography } from 'antd';
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from '../utils/apis';

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
  const [config, setConfig] = useState({});

  const getInfo = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      console.log("viewers result", result)

      setConfig({ ...result });

    } catch (error) {
      setConfig({ ...config, message: error.message });
    }
  };
  
  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();
    getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
  
    // returned function will be called on component unmount 
    return () => {
      clearInterval(getStatusIntervalId);
    }
  }, []);

  return (
    <div>
      <h2>Server Config</h2>
      <p>
        Display this data all pretty, most things will be editable in the
        future, not now.
      </p>
      <div
        style={{
          border: "1px solid pink",
          width: "100%",
          overflow: "auto",
        }}
      >
        <VideoVariants config={config} />
      </div>
    </div>
  ); 
}

