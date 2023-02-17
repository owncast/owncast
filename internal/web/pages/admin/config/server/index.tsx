import React, { ReactElement } from 'react';
import { Tabs } from 'antd';

import StreamKeys from '../../../../components/admin/config/server/StreamKeys';
import ServerConfig from '../../../../components/admin/config/server/ServerConfig';
import StorageConfig from '../../../../components/admin/config/server/StorageConfig';

import { AdminLayout } from '../../../../components/layouts/AdminLayout';

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

PublicFacingDetails.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
