import React from 'react';
import { Typography } from 'antd';
import Link from 'next/link';

import EditInstanceDetails from './components/config/edit-instance-details';
import EditDirectoryDetails from './components/config/edit-directory';
import EditInstanceTags from './components/config/edit-tags';

const { Title } = Typography;

export default function PublicFacingDetails() {
  return (
    <>
      <Title level={2}>Edit your public facing instance details</Title>

      <div className="edit-public-details-container">
        <EditInstanceDetails />
        <EditInstanceTags />
        <EditDirectoryDetails />
      </div>
    </>
  );
}
