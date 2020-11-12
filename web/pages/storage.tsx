/* eslint-disable react/prop-types */
import React, { useState, useEffect } from "react";
import { SERVER_CONFIG, fetchData, FETCH_INTERVAL } from "../utils/apis";
import KeyValueTable from "./components/key-value-table";

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
  const [config, setConfig] = useState({});

  const getInfo = async () => {
    try {
      const result = await fetchData(SERVER_CONFIG);
      console.log("viewers result", result);

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
    };
  }, []);

  return (
    <div>
        <Storage config={config} />
    </div>
  );
}
