import React, { useContext } from 'react';
import { Table, Typography, Input } from 'antd';
import { isEmptyObject } from '../utils/format';
import KeyValueTable from "./components/key-value-table";
import { ServerStatusContext } from '../utils/server-status-context';
import adminStyles from '../styles/styles.module.scss';

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
      <div className={adminStyles.configSection}>
        <Title level={2}>Social Handles</Title>
        <Table
          pagination={false}
          columns={columns}
          dataSource={config.instanceDetails.socialHandles}
          rowKey="platform"
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
      value: instanceDetails.logo?.large,
    },
    {
      name: "Tags",
      value: instanceDetails.tags?.join(", "),
    },
    {
      name: "NSFW",
      value: instanceDetails.nsfw?.toString(),
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
      <div className={adminStyles.configSection}>
        <KeyValueTable title="Server details" data={data} />
      </div>
      <div className={adminStyles.configSection}>
        <KeyValueTable title="Server configuration" data={configData} />
      </div>
    </>
  );
}

function PageContent({ config }) {
  if (!config?.instanceDetails?.extraPageContent) {
    return null;
  }
  return (
    <div className={adminStyles.configSection}>
      <Title level={2}>Page content</Title>
      <TextArea
        disabled rows={4}
        value={config.instanceDetails.extraPageContent}
      />
    </div>
  );
}

export default function ServerConfig() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig: config } = serverStatusData || {};

  return (
    <>
      <InstanceDetails config={config} />
      <SocialHandles config={config} />
      <PageContent config={config} />
      
      <Title level={5}>Learn more about configuring Owncast <a href="https://owncast.online/docs/configuration">by visiting the documentation.</a></Title>
    </>
  ); 
}

