import React from 'react';
import { Typography } from 'antd';
import Link from 'next/link';

import EditInstanceDetails from './components/config/edit-instance-details';

const { Title } = Typography;

export default function PublicFacingDetails() {
  return (
    <>
      <Title level={2}>Edit your public facing instance details</Title>

      <div className={`publicDetailsContainer`}>
        <div className={`textFieldsSection`}>
          <EditInstanceDetails />

          <Link href="/admin/config-page-content">
            <a>Edit your extra page content here.</a>
          </Link>
        </div>
      </div>
    </>
  );
}
