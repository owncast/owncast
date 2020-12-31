import React, { useContext, useEffect } from 'react';
import { Typography, Form } from 'antd';

import TextField, { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD, TEXTFIELD_TYPE_TEXTAREA } from './components/config/form-textfield';

import { ServerStatusContext } from '../utils/server-status-context';

const { Title } = Typography;

export default function ConfigServerDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { ffmpegPath, streamKey, webServerPort } = serverConfig;

  const streamDetails = {
    ffmpegPath, streamKey, webServerPort
  };
  
  useEffect(() => {
    form.setFieldsValue({...streamDetails});
  }, [serverStatusData]);


  const handleResetValue = (fieldName: string) => {
    form.setFieldsValue({ [fieldName]: streamDetails[fieldName]});
  }

  const extraProps = {
    handleResetValue,
    initialValues: streamDetails,
  };

  console.log(streamDetails)
  return (
    <>
      <Title level={2}>Edit your Server&apos;s details</Title>

      <div className="config-public-details-container">
        <Form
          form={form}
          layout="vertical"
        >
          <TextField fieldName="streamKey" type={TEXTFIELD_TYPE_PASSWORD} {...extraProps} />
          <TextField fieldName="ffmpegPath" type={TEXTFIELD_TYPE_TEXTAREA} {...extraProps} />
          <TextField fieldName="webServerPort" type={TEXTFIELD_TYPE_NUMBER} {...extraProps} />
        </Form>
      </div>      
    </>
  ); 
}


