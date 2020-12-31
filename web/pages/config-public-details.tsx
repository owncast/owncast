import React, { useContext, useEffect } from 'react';
import { Typography, Form } from 'antd';

import TextField, { TEXTFIELD_TYPE_TEXTAREA } from './components/config/form-textfield';

import EditInstanceTags from './components/config/tags';

import { ServerStatusContext } from '../utils/server-status-context';

const { Title } = Typography;

export default function PublicFacingDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { instanceDetails = {} } = serverConfig;
  console.log(serverConfig)
  
  useEffect(() => {
    form.setFieldsValue({...instanceDetails});
  }, [instanceDetails]);


  const handleResetValue = (fieldName: string) => {
    form.setFieldsValue({ [fieldName]: instanceDetails[fieldName]});
  }

  const extraProps = {
    handleResetValue,
    initialValues: instanceDetails,
  };

  return (
    <>
      <Title level={2}>Edit your public facing instance details</Title>

      <div className="config-public-details-container">
        <div className="text-fields">
          <Form
            form={form}
            layout="vertical"
          >
            <TextField fieldName="name" {...extraProps} />
            <TextField fieldName="summary" type={TEXTFIELD_TYPE_TEXTAREA} {...extraProps} />
            <TextField fieldName="title" {...extraProps} />
            <TextField fieldName="streamTitle" {...extraProps} />
          </Form>
        </div>
        <div className="misc-fields">
          {/* add social handles comp
          <br/>
          add tags comp */}
          <EditInstanceTags />
          
          
        </div>
      </div>      
    </>
  ); 
}


