import React from 'react';
import { Typography } from 'antd';

import VideoVariantsTable from './components/config/video-variants-table';
import VideoLatency from './components/config/video-latency';

const { Title } = Typography;

export default function ConfigVideoSettings() {
  return (
    <div className="config-video-variants">
      <Title level={2}>Video configuration</Title>
      <p>
        Learn more about configuring Owncast{' '}
        <a href="https://owncast.online/docs/configuration">by visiting the documentation.</a>
      </p>

    <p>
      <VideoVariantsTable />
    </p>

    <p>
      <VideoLatency />
    </p>
    
    </div>
  );
}
