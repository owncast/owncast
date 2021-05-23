import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { Table, Typography } from 'antd';
import { getGithubRelease } from '../utils/apis';

const { Title } = Typography;

function AssetTable(assets) {
  const data = Object.values(assets) as object[];

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (text, entry) => <a href={entry.browser_download_url}>{text}</a>,
    },
    {
      title: 'Size',
      dataIndex: 'size',
      key: 'size',
      render: text => `${(text / 1024 / 1024).toFixed(2)} MB`,
    },
  ];

  return (
    <Table
      dataSource={data}
      columns={columns}
      rowKey={record => record.id}
      size="large"
      pagination={false}
    />
  );
}

export default function Logs() {
  const [release, setRelease] = useState({
    html_url: '',
    name: '',
    created_at: null,
    body: '',
    assets: [],
  });

  const getRelease = async () => {
    try {
      const result = await getGithubRelease();
      setRelease(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    getRelease();
  }, []);

  if (!release) {
    return null;
  }

  return (
    <div className="upgrade-page">
      <Title level={2}>
        <a href={release.html_url}>{release.name}</a>
      </Title>
      <Title level={5}>{new Date(release.created_at).toDateString()}</Title>
      <ReactMarkdown>{release.body}</ReactMarkdown>
      <h3>Downloads</h3>
      <AssetTable {...release.assets} />
    </div>
  );
}
