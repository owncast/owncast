import React, { ReactElement } from 'react';
import { Tabs } from 'antd';
import { useTranslation } from 'next-export-i18n';

import GeneralConfig from '../../../../components/admin/config/general/GeneralConfig';
import AppearanceConfig from '../../../../components/admin/config/general/AppearanceConfig';

import { AdminLayout } from '../../../../components/layouts/AdminLayout';
import { EditCustomJavascript } from '../../../../components/admin/EditCustomJavascript';
import { Localization } from '../../../../types/localization';

export default function PublicFacingDetails() {
  const { t } = useTranslation();

  return (
    <div className="config-public-details-page">
      <Tabs
        defaultActiveKey="1"
        centered
        items={[
          {
            label: t(Localization.Admin.Config.general),
            key: '1',
            children: <GeneralConfig />,
          },
          {
            label: t(Localization.Admin.Config.appearance),
            key: '2',
            children: <AppearanceConfig />,
          },
          {
            label: t(Localization.Admin.Config.customScripting),
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
