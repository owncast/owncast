import React from 'react';
import { Tabs } from 'antd';

import StreamKeys from './components/StreamKeys';
import ServerConfig from './components/ServerConfig';
import StorageConfig from './components/StorageConfig';

export default function PublicFacingDetails() {
  return (
    <div className="config-public-details-page">
      <Tabs
        defaultActiveKey="1"
        centered
        items={[
          {
            label: `Server Config`,
            key: '1',
            children: <ServerConfig />,
          },
          {
            label: `Stream Keys`,
            key: '2',
            children: <StreamKeys />,
          },
          {
            label: `S3 Object Storage`,
            key: '3',
            children: <StorageConfig />,
          },
        ]}
      />
    </div>
  );
}
