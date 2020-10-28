import React, { useState, useEffect } from 'react';
import { Table, Typography, Input } from 'antd';
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from './utils/apis';
import { isEmptyObject } from './utils/format';

const { Title } = Typography;
const { TextArea } = Input;

function KeyValueTable({ title, data }) {
  const columns = [
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
  
  return (
    <div>
      <Title>{title}</Title>
      <Table pagination={false} columns={columns} dataSource={data} />
    </div>
  );
}

function SocialHandles({ config }) {
  if (!config) {
    return null;
  }

  const columns = [
    {
      title: "Platform",
      dataIndex: "platform",
      key: "platform",
    },
    {
      title: "URL",
      dataIndex: "url",
      key: "url",
      render: (url) => `<a href="${url}">${url}</a>`
    },
  ];

  return (
    <div>
      <Title>Social Handles</Title>
      <Table
        pagination={false}
        columns={columns}
        dataSource={config.instanceDetails.socialHandles}
      />
    </div>
  );
}

function InstanceDetails({ config }) {
  console.log(config)
  if (!config || isEmptyObject(config)) {
    return null;
  }

  const { instanceDetails = {}, yp, streamKey, ffmpegPath, webServerPort } = config;
  
  const data = [
    {
      name: "Server name",
      value: instanceDetails.name,
    },
    {
      name: "Title",
      value: instanceDetails.title,
    },
    {
      name: "Summary",
      value: instanceDetails.summary,
    },
    {
      name: "Logo",
      value: instanceDetails.logo.large,
    },
    {
      name: "Tags",
      value: instanceDetails.tags.join(", "),
    },
    {
      name: "NSFW",
      value: instanceDetails.nsfw.toString(),
    },
    {
      name: "Shows in Owncast directory",
      value: yp.enabled.toString(),
    },
  ];

  const configData = [
    {
      name: "Stream key",
      value: streamKey,
    },
    {
      name: "ffmpeg path",
      value: ffmpegPath,
    },
    {
      name: "Web server port",
      value: webServerPort,
    },
  ];

  return (
    <>
      <KeyValueTable title="Server details" data={data} />
      <KeyValueTable title="Server configuration" data={configData} />
    </>
  );
}

function VideoVariants({ config }) {
  if (!config) {
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

function Storage({ config }) {
  if (!config) {
    return null;
  }

  const data = [
    {
      name: "Enabled",
      value: config.s3.enabled.toString(),
    },
    {
      name: "Endpoint",
      value: config.s3.endpoint,
    },
    {
      name: "Access Key",
      value: config.s3.accessKey,
    },
    {
      name: "Secret",
      value: config.s3.secret,
    },
    {
      name: "Bucket",
      value: config.s3.bucket,
    },
    {
      name: "Region",
      value: config.s3.region,
    },
  ];
  return <KeyValueTable title="External Storage" data={data} />
}

function PageContent({ config }) {
  if (!config) {
    return null;
  }
  return (
    <div>
      <Title>Page content</Title>
      <TextArea
        disabled rows={4}
        value={config.instanceDetails.extraPageContent}
      />
    </div>
  );
}

export default function ServerConfig() {  
  const [config, setConfig] = useState();

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
        <InstanceDetails config={config} />
        <SocialHandles config={config} />
        <VideoVariants config={config} />
        <Storage config={config} />
        <PageContent config={config} />

        {JSON.stringify(config)}
      </div>
    </div>
  ); 
}

