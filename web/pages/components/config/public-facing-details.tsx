import React, { useContext, useEffect } from 'react';
import { Typography, Form, Input } from 'antd';

import TextField from './form-textfield';

import { ServerStatusContext } from '../../../utils/server-status-context';

import { UpdateArgs } from '../../../types/config-section';

const { Title } = Typography;

export default function PublicFacingDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setConfigField } = serverStatusData || {};

  const { instanceDetails = {},  } = serverConfig;

  const { name, summary, title } = instanceDetails;
  
  useEffect(() => {
    form.setFieldsValue({...instanceDetails});
  }, [instanceDetails]);
  

  return (
    <>
      <Title level={2}>Edit your public facing instance details</Title>

      <div className="config-public-details-container">
        <div className="text-fields">
          <Form
            form={form}
            layout="vertical"
          >
            <TextField fieldName="name" />
            <TextField fieldName="summary" />
          </Form>
        </div>
        <div className="misc-optionals">
          add social handles comp
          <br/>
          add tags comp
          
        </div>
      </div>      
    </>
  ); 
}

