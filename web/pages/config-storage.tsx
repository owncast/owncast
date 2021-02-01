import React, { useContext, useState, useEffect } from 'react';
import { ServerStatusContext } from '../utils/server-status-context';
import { Typography, Switch, Input, Button } from 'antd';
import {
  postConfigUpdateToAPI,
  API_S3_INFO,
} from './components/config/constants';
const { Title } = Typography;

function Storage({ config }) {
  if (!config || !config.s3) {
    return null;
  }

  const [endpoint, setEndpoint] = useState(config.s3.endpoint);
  const [accessKey, setAccessKey] = useState(config.s3.accessKey);
  const [secret, setSecret] = useState(config.s3.secret);
  const [bucket, setBucket] = useState(config.s3.bucket);
  const [region, setRegion] = useState(config.s3.region);
  const [acl, setAcl] = useState(config.s3.acl);
  const [servingEndpoint, setServingEndpoint] = useState(config.s3.servingEndpoint);
  const [enabled, setEnabled] = useState(config.s3.enabled);

  function storageEnabledChanged(storageEnabled) {
    setEnabled(storageEnabled);
  }

  function endpointChanged(e) {
    setEndpoint(e.target.value)
  }

  function accessKeyChanged(e) {
    setAccessKey(e.target.value)
  }

  function secretChanged(e) {
    setSecret(e.target.value)
  }

  function bucketChanged(e) {
    setBucket(e.target.value)
  }

  function regionChanged(e) {
    setRegion(e.target.value)
  }

  function aclChanged(e) {
    setAcl(e.target.value)
  }

  function servingEndpointChanged(e) {
    setServingEndpoint(e.target.value)
  }

  async function save() {
    const payload = {
      value: {
        enabled: enabled,
        endpoint: endpoint,
        accessKey: accessKey,
        secret: secret,
        bucket: bucket,
        region: region,
        acl: acl,
        servingEndpoint: servingEndpoint,
      }
    };

    try {
      await postConfigUpdateToAPI({apiPath: API_S3_INFO, data: payload});
    } catch(e) {
      console.error(e);
    }
  }

  const table = enabled ? (
    <>
    <br></br>
    endpoint <Input defaultValue={endpoint} value={endpoint} onChange={endpointChanged} />
    access key<Input label="Access key" defaultValue={accessKey} value={accessKey} onChange={accessKeyChanged} />
    secret <Input label="Secret" defaultValue={secret} value={secret} onChange={secretChanged} />
    bucket <Input label="Bucket" defaultValue={bucket} value={bucket} onChange={bucketChanged} />
    region <Input label="Region" defaultValue={region} value={region} onChange={regionChanged} />
    advanced<br></br>
    acl <Input label="ACL" defaultValue={acl} value={acl} onChange={aclChanged} />
    serving endpoint <Input label="Serving endpoint" defaultValue={servingEndpoint} value={servingEndpoint} onChange={servingEndpointChanged} />
    <Button onClick={save}>Save</Button>
    </>
  ): null;

  return (
    <>
      <Title level={2}>Storage</Title>
      Enabled: 
      <Switch checked={enabled} onChange={storageEnabledChanged} />
      { table }
    </>
  );
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
