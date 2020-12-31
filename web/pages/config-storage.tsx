import React, { useContext } from "react";
import KeyValueTable from "./components/key-value-table";
import { ServerStatusContext } from '../utils/server-status-context';
import { Typography } from 'antd';
import Link from 'next/link';

const { Title } = Typography;

function Storage({ config }) {
  if (!config || !config.s3) {
    return null;
  }

  if (!config.s3.enabled) {
    return (
      <div>
      <Title>External Storage</Title>
      <p>
        You are currently using the <Link href="/hardware-info">local storage of this Owncast server</Link> to store and distribute video.
      </p>
      <p>
      Owncast can use S3-compatible external storage providers to offload the responsibility of disk and bandwidth utilization from your own server.
      </p>

      <p>
        Visit our <a href="https://owncast.online/docs/s3/">storage documentation</a> to learn how to configure this.
      </p>
      </div>
    );
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
  return <KeyValueTable title="External Storage" data={data} />;
}

export default function ServerConfig() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig: config } = serverStatusData || {};

  return (
    <div>
        <Storage config={config} />
    </div>
  );
}
