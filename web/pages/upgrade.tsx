import React, { useState, useEffect } from "react";
import ReactMarkdown from "react-markdown";
import { Table, Tag } from "antd";

const API_URL = "https://api.github.com/repos/owncast/owncast/releases/latest";

export default function Logs() {
  const [release, setRelease] = useState({
    html_url: "",
    name: "",
    created_at: null,
    body: "",
    assets: [],

  });

  const getRelease = async () => {
    try {
      const result = await fetchData(API_URL);
      setRelease(result);
    } catch (error) {
      console.log("==== error", error);
    }
  };

  useEffect(() => {
    getRelease();
  }, []);

  if (!release) {
    return null;
  }

  return (
    <div>
      <h2>
        <a href={release.html_url}>{release.name}</a>
      </h2>
      <h1>{release.created_at}</h1>
      <ReactMarkdown>{release.body}</ReactMarkdown>;

      <h3>Downloads</h3>
      <AssetTable {...release.assets} />
    </div>
  );
}

function AssetTable(assets) {
  const data = [];

  for (const key in assets) {
    data.push(assets[key]);
  }

  const columns = [
    {
      title: "Name",
      dataIndex: "name",
      key: "name",
      render: (text, entry) =>
        <a href={entry.browser_download_url}>{text}</a>,
    },
    {
      title: "Size",
      dataIndex: "size",
      key: "size",
      render: (text, entry) =>
        `${(text/1024/1024).toFixed(2)} MB`
    },
  ];

  return <Table dataSource={data} columns={columns} rowKey="id" size="large" />;
}

async function fetchData(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      const message = `An error has occured: ${response.status}`;
      throw new Error(message);
    }
    const json = await response.json();
    return json;
  } catch (error) {
    console.log(error);
  }
  return {};
}
