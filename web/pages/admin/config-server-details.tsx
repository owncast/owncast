import React from 'react';
import { Typography } from 'antd';
import { EditInstanceDetails } from '../../components/config/EditInstanceDetails2';

const { Title } = Typography;

export default function ConfigServerDetails() {
  return (
    <div className="config-server-details-form">
      <Title>Server Settings</Title>
      <p className="description">
        You should change your stream key from the default and keep it safe. For most people
        it&apos;s likely the other settings will not need to be changed.
      </p>
      <div className="form-module config-server-details-container">
        <EditInstanceDetails />
      </div>
    </div>
  );
}
