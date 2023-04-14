import React, { ReactElement, useContext } from 'react';
import { Tabs } from 'antd';

import StreamKeys from '../../../../components/admin/config/server/StreamKeys';
import ServerConfig from '../../../../components/admin/config/server/ServerConfig';
import StorageConfig from '../../../../components/admin/config/server/StorageConfig';
import { ServerStatusContext } from '../../../../utils/server-status-context';

import { AdminLayout } from '../../../../components/layouts/AdminLayout';

export default function PublicFacingDetails() {
  const serverStatusData = useContext(ServerStatusContext);

  const { serverConfig } = serverStatusData || {};
  const { streamKeyOverridden } = serverConfig;

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
          !streamKeyOverridden && {
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
