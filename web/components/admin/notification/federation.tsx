import { Button, Typography } from 'antd';
import React, { useState, useContext, useEffect } from 'react';
import Link from 'next/link';
import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;

export const FediverseNotify = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { federation } = serverConfig || {};

  const { enabled } = federation || {};
  const [formDataValues, setFormDataValues] = useState<any>({});

  useEffect(() => {
    setFormDataValues({
      enabled,
    });
  }, [enabled]);

  return (
    <>
      <Title>Fediverse Social</Title>
      <p className="description">
        Enabling the Fediverse social features will not just alert people to when you go live, but
        also enable other functionality.
      </p>
      <p>
        Fediverse social features:{' '}
        <span style={{ color: federation.enabled ? 'green' : 'red' }}>
          {formDataValues.enabled ? 'Enabled' : 'Disabled'}
        </span>
      </p>

      <Link passHref href="/admin/config-federation/">
        <Button
          type="primary"
          style={{
            position: 'relative',
            marginLeft: 'auto',
            right: '0',
            marginTop: '20px',
          }}
        >
          Configure
        </Button>
      </Link>
    </>
  );
};
