import React, { useContext } from 'react';
import { Typography, Input, Form } from 'antd';

import TextField from './form-textfield';


import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;

export default function PublicFacingDetails() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setConfigField } = serverStatusData || {};

  const { instanceDetails = {},  } = serverConfig;

  const { name, summary, title } = instanceDetails;

  return (
    <>
      <Title level={2}>Edit your public facing instance details</Title>
      <div className="config-public-details-container">
        <div className="text-fields">
          <Form name="text-fields" labelCol={{ span: 8 }} wrapperCol={{ span: 16 }}>
            <TextField
              label="name"
              value={name}
            />
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

