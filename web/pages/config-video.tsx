import React, { useContext, useEffect } from 'react';
import { Typography, Form, Slider } from 'antd';
import { ServerStatusContext } from '../utils/server-status-context';

import VideoVariantsTable from './components/config/video-variants-table';
import TextField, { TEXTFIELD_TYPE_NUMBER } from './components/config/form-textfield';
import { TEXTFIELD_DEFAULTS } from './components/config/constants';

const { Title } = Typography;

export default function VideoConfig() {
  const [form] = Form.useForm();
  const serverStatusData = useContext(ServerStatusContext);
  // const { serverConfig } = serverStatusData || {};
  // const { videoSettings } = serverConfig || {};
  // const { numberOfPlaylistItems, segmentLengthSeconds } = videoSettings || {};

  const videoSettings = serverStatusData?.serverConfig?.videoSettings;
  const { numberOfPlaylistItems, segmentLengthSeconds } = videoSettings || {};
  const initialValues = {
    numberOfPlaylistItems,
    segmentLengthSeconds,
  };

  useEffect(() => {
    form.setFieldsValue(initialValues);
  }, [serverStatusData]);

  const handleResetValue = (fieldName: string) => {
    const defaultValue = TEXTFIELD_DEFAULTS.videoSettings[fieldName] && TEXTFIELD_DEFAULTS.videoSettings[fieldName].defaultValue || '';

    form.setFieldsValue({ [fieldName]: initialValues[fieldName] || defaultValue });
  }

  const extraProps = {
    handleResetValue,
    initialValues: videoSettings,
    configPath: 'videoSettings',
  };

  return (
    <div className="config-video-variants">
      <Title level={2}>Video configuration</Title>
      <Title level={5}>Learn more about configuring Owncast <a href="https://owncast.online/docs/configuration">by visiting the documentation.</a></Title>

        <div style={{ wordBreak: 'break-word'}}>
          {JSON.stringify(videoSettings)}
        </div>
        <div className="config-video-misc">
          <Form
            form={form}
            layout="vertical"
          >
            <TextField fieldName="numberOfPlaylistItems" type={TEXTFIELD_TYPE_NUMBER} {...extraProps} />
            <TextField fieldName="segmentLengthSeconds" type={TEXTFIELD_TYPE_NUMBER} {...extraProps} />
          </Form>
        </div>  

        <VideoVariantsTable />
    </div>
  ); 
}

