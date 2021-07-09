import { Col, Collapse, Row, Typography } from 'antd';
import React from 'react';
import VideoCodecSelector from '../components/config/video-codec-selector';
import VideoLatency from '../components/config/video-latency';
import VideoVariantsTable from '../components/config/video-variants-table';

const { Panel } = Collapse;
const { Title } = Typography;

export default function ConfigVideoSettings() {
  return (
    <div className="config-video-variants">
      <Title>Video configuration</Title>
      <p className="description">
        Before changing your video configuration{' '}
        <a
          href="https://owncast.online/docs/video?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          visit the video documentation
        </a>{' '}
        to learn how it impacts your stream performance. The general rule is to start conservatively
        by having one middle quality stream output variant and experiment with adding more of varied
        qualities.
      </p>

      <Row gutter={[16, 16]}>
        <Col md={24} lg={12}>
          <div className="form-module variants-table-module">
            <VideoVariantsTable />
          </div>
        </Col>
        <Col md={24} lg={12}>
          <div className="form-module latency-module">
            <VideoLatency />
          </div>

          <Collapse className="advanced-settings codec-module">
            <Panel header="Advanced Settings" key="1">
              <div className="form-module variants-table-module">
                <VideoCodecSelector />
              </div>
            </Panel>
          </Collapse>
        </Col>
      </Row>
    </div>
  );
}
