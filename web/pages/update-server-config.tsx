/* eslint-disable react/prop-types */
import React, { useState, useEffect } from 'react';
import { Table, Typography, Input } from 'antd';
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from '../utils/apis';
import { isEmptyObject } from '../utils/format';
import KeyValueTable from "./components/key-value-table";

const { Title } = Typography;
const { TextArea } = Input;


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
      render: (url) => <a href={url}>{url}</a>
    },
  ];

  if (!config.instanceDetails?.socialHandles) {
    return null;
  }

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

function PageContent({ config }) {
  if (!config?.instanceDetails?.extraPageContent) {
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
        <InstanceDetails config={config} />
        <SocialHandles config={config} />
        <PageContent config={config} />
    </div>
  ); 
}

