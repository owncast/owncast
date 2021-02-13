import React from 'react';
import { Typography } from 'antd';
import EditServerDetails from '../components/config/edit-server-details';

const { Title } = Typography;

export default function ConfigServerDetails() {
  return (
    <div className="config-server-details-form">
      <Title level={2} className="page-title">
        Server Settings
      </Title>
      <p className="description">
        You should change your stream key from the default and keep it safe. For most people
        it&apos;s likely the other settings will not need to be changed.
      </p>
      <div className="form-module config-server-details-container">
        <EditServerDetails />
      </div>
    </div>
  );
}
