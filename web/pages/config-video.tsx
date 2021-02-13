import React from 'react';
import { Typography } from 'antd';

import VideoVariantsTable from '../components/config/video-variants-table';
import VideoLatency from '../components/config/video-latency';

const { Title } = Typography;

export default function ConfigVideoSettings() {
  return (
    <div className="config-video-variants">
      <Title level={2} className="page-title">
        Video configuration
      </Title>
      <p className="description">
        Before changing your video configuration{' '}
        <a href="https://owncast.online/docs/encoding">visit the video documentation</a> to learn
        how it impacts your stream performance.
      </p>

      <div className="row">
        <div className="form-module variants-table-module">
          <VideoVariantsTable />
        </div>

        <div className="form-module latency-module">
          <VideoLatency />
        </div>
      </div>
    </div>
  );
}
