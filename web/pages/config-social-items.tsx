import React from 'react';
import { Typography } from 'antd';
import EditSocialLinks from '../components/config/edit-social-links';

const { Title } = Typography;

export default function ConfigSocialThings() {
  return (
    <div className="config-social-items">
      <Title level={2}>Social Items</Title>

      <EditSocialLinks />
    </div>
  );
}
