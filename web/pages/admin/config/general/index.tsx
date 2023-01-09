import React from 'react';
import { Tabs } from 'antd';

import GeneralConfig from './components/GeneralConfig';
import AppearanceConfig from './components/AppearanceConfig';

export default function PublicFacingDetails() {
  return (
    <div className="config-public-details-page">
      <Tabs
        defaultActiveKey="1"
        centered
        items={[
          {
            label: `General`,
            key: '1',
            children: <GeneralConfig />,
          },
          {
            label: `Appearance`,
            key: '2',
            children: <AppearanceConfig />,
          },
        ]}
      />
    </div>
  );
}
