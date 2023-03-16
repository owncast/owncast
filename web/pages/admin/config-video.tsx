import { Col, Collapse, Row, Typography } from 'antd';
import React, { ReactElement } from 'react';
import { CodecSelector as VideoCodecSelector } from '../../components/admin/CodecSelector';
import { VideoLatency } from '../../components/admin/VideoLatency';
import { CurrentVariantsTable } from '../../components/admin/CurrentVariantsTable';

import { AdminLayout } from '../../components/layouts/AdminLayout';

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

      <Row gutter={[45, 16]}>
        <Col md={24} lg={12}>
          <div className="form-module variants-table-module">
            <CurrentVariantsTable />
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

ConfigVideoSettings.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
