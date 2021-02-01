import { Typography } from 'antd';
import React from 'react';
import EditStorage from './components/config/edit-storage';

const { Title } = Typography;

export default function ConfigStorageInfo() {
  return (
    <>
      <Title level={2}>Storage</Title>
      <EditStorage />
    </>
  );
}
