/* eslint-disable react/prop-types */
import React, { useContext } from "react";
import KeyValueTable from "./components/key-value-table";
import { ServerStatusContext } from '../utils/server-status-context';

function Storage({ config }) {
  if (!config || !config.s3) {
    return null;
  }

  if (!config.s3.enabled) {
    return (
      <h3>
        Local storage is being used. Enable external S3 storage if you want
        to use it.  TODO: Make this message somewhat more informative.  Point to documentation or something.
      </h3>
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
