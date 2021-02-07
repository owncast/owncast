import React from 'react';
import { Typography } from 'antd';

import VideoVariantsTable from '../components/config/video-variants-table';
import VideoLatency from '../components/config/video-latency';

const { Title } = Typography;

export default function ConfigVideoSettings() {
  return (
    <div className="config-video-variants">
      <Title level={2}>Video configuration</Title>
      <p>
        Before changing your video configuration{' '}
        <a href="https://owncast.online/docs/encoding">visit the video documentation</a> to learn
        how it impacts your stream performance.
      </p>

      <p>
        <VideoVariantsTable />
      </p>
      <br />
      <hr />
      <br />
      <p>
        <VideoLatency />
      </p>
    </div>
  );
}
