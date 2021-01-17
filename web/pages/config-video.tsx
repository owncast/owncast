import React, { useContext } from 'react';
import { Typography } from 'antd';
import { ServerStatusContext } from '../utils/server-status-context';

import VideoVariantsTable from './components/config/video-variants-table';
import VideoSegmentsEditor from './components/config/video-segments-editor';

const { Title } = Typography;

export default function VideoConfig() {
  // const [form] = Form.useForm();
  // const serverStatusData = useContext(ServerStatusContext);
  // const { serverConfig } = serverStatusData || {};
  // const { videoSettings } = serverConfig || {};
  return (
    <div className="config-video-variants">
      <Title level={2}>Video configuration</Title>
      <p>Learn more about configuring Owncast <a href="https://owncast.online/docs/configuration">by visiting the documentation.</a></p>

        {/* <div style={{ wordBreak: 'break-word'}}>
          {JSON.stringify(videoSettings)}
        </div> */}
        
        <VideoSegmentsEditor />

        <br /><br />

        <VideoVariantsTable />
    </div>
  ); 
}

