import React from 'react';
import { Typography } from 'antd';
import EditServerDetails from './components/config/edit-server-details';

const { Title } = Typography;

export default function ConfigServerDetails() {
  return (
    <div className="config-server-details-form">
      <Title level={2}>Edit your Server&apos;s details</Title>

      <div className="config-server-details-container">
        <EditServerDetails />
        
      </div>      
    </div>
  ); 
}


