import React from 'react';
import { Typography } from 'antd';
import Link from 'next/link';

import EditInstanceDetails from '../components/config/edit-instance-details';
import EditDirectoryDetails from '../components/config/edit-directory';
import EditInstanceTags from '../components/config/edit-tags';

const { Title } = Typography;

export default function PublicFacingDetails() {
  return (
    <>
      <Title level={2}>General Settings</Title>
      <p>
        The following are displayed on your site to describe your stream and its content.{' '}
        <a href="https://owncast.online/docs/website/">Learn more.</a>
      </p>
      <div className="edit-public-details-container">
        <EditInstanceDetails />
        <EditInstanceTags />
        <EditDirectoryDetails />
      </div>
    </>
  );
}
