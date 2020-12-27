import React, { useContext } from 'react';
import { Typography, Input } from 'antd';


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
        <div className="text-fields" role="form">
            Server Name
            <Input placeholder="Owncast" value={name} />


        </div>
        <div className="misc-optionals">
          add social handles
          <br/>
          add tags 
          
        </div>
      </div>      
    </>
  ); 
}

