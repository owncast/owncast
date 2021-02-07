import { Typography } from 'antd';
import React from 'react';
import EditStorage from '../components/config/edit-storage';

const { Title } = Typography;

export default function ConfigStorageInfo() {
  return (
    <>
      <Title level={2}>Storage</Title>
      <p>
        Owncast supports optionally using external storage providers to distribute your video. Learn
        more about this by visiting our{' '}
        <a href="https://owncast.online/docs/storage/">Storage Documentation</a>.
      </p>
      <p>
        Configuring this incorrectly will likely cause your video to be unplayable. Double check the
        documentation for your storage provider on how to configure the bucket you created for
        Owncast.
      </p>
      <EditStorage />
    </>
  );
}
