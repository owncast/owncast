import React, { ReactElement } from 'react';
import { Typography } from 'antd';
import EditSocialLinks from '../../components/admin/config/general/EditSocialLinks';

import { AdminLayout } from '../../components/layouts/AdminLayout';

const { Title } = Typography;

export default function ConfigSocialThings() {
  return (
    <div className="config-social-items">
      <Title>Social Items</Title>

      <EditSocialLinks />
    </div>
  );
}

ConfigSocialThings.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
