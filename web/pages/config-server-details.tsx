import React, { useContext, useEffect } from 'react';
import { Typography, Form } from 'antd';

import TextField, { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD, TEXTFIELD_TYPE_TEXTAREA } from './components/config/form-textfield';

import { ServerStatusContext } from '../utils/server-status-context';
import { TEXTFIELD_DEFAULTS } from './components/config/constants';

const { Title } = Typography;

export default function ConfigServerDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { ffmpegPath, streamKey, webServerPort, rtmpServerPort } = serverConfig;

  const initialValues = {
    ffmpegPath,
    streamKey,
    webServerPort,
    rtmpServerPort,
  };
  
  useEffect(() => {
    form.setFieldsValue(initialValues);
  }, [serverStatusData]);

  const handleResetValue = (fieldName: string) => {
    const defaultValue = TEXTFIELD_DEFAULTS[fieldName] && TEXTFIELD_DEFAULTS[fieldName].defaultValue || '';

    form.setFieldsValue({ [fieldName]: initialValues[fieldName] || defaultValue });
  }

  const extraProps = {
    handleResetValue,
    initialValues,
    configPath: '',
  };
  return (
    <>
      <Title level={2}>Edit your Server&apos;s details</Title>

      <div className="config-public-details-container">
        <Form
          form={form}
          layout="vertical"
        >
          <TextField fieldName="streamKey" type={TEXTFIELD_TYPE_PASSWORD} {...extraProps} />
          <TextField fieldName="ffmpegPath" {...extraProps} />
          <TextField fieldName="webServerPort" type={TEXTFIELD_TYPE_NUMBER} {...extraProps} />
          <TextField fieldName="rtmpServerPort" type={TEXTFIELD_TYPE_NUMBER} {...extraProps} />
        </Form>
      </div>      
    </>
  ); 
}


