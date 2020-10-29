import React, { useState, useEffect } from 'react';
import { Table, Typography, Input } from 'antd';
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from './utils/apis';
import { isEmptyObject } from './utils/format';
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
        <PageContent config={config} />

        {JSON.stringify(config)}
      </div>
    </div>
  ); 
}

