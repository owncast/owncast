import React, { ReactElement } from 'react';
import { Tabs } from 'antd';

import GeneralConfig from '../../../../components/admin/config/general/GeneralConfig';
import AppearanceConfig from '../../../../components/admin/config/general/AppearanceConfig';

import { AdminLayout } from '../../../../components/layouts/AdminLayout';
import { EditCustomJavascript } from '../../../../components/admin/EditCustomJavascript';

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
          {
            label: `Custom Scripting`,
            key: '3',
            children: <EditCustomJavascript />,
          },
        ]}
      />
    </div>
  );
}

PublicFacingDetails.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
