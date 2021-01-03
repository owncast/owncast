// Note: references to "yp" in the app are likely related to Owncast Directory
import React, { useContext, useEffect } from 'react';
import { Typography, Form } from 'antd';

import ToggleSwitch from './form-toggleswitch';

import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;

export default function EditYPDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { yp, instanceDetails } = serverConfig;
  const { nsfw } = instanceDetails;
  const { enabled, instanceUrl } = yp;

  const initialValues = {
    ...yp,
    enabled,
    nsfw,
  };

  const hasInstanceUrl = instanceUrl !== '';
  
  useEffect(() => {
    form.setFieldsValue(initialValues);
  }, [yp]);

  const extraProps = {
    initialValues,
    disabled: !hasInstanceUrl,
  };

  // TODO: DISABLE THIS SECTION UNTIL instanceURL is populated
  return (
    <div className="config-directory-details-form">
      <Title level={3}>Owncast Directory Settings</Title>
      
      <p>Would you like to appear in the <a href="https://directory.owncast.online" target="_blank" rel="noreferrer"><strong>Owncast Directory</strong></a>?</p>

      <p><em>NOTE: You will need to have a URL specified in the <code>Instance URL</code> field to be able to use this.</em></p>

      <div className="config-yp-container">
        <Form
          form={form}
          layout="vertical"
        >
          <ToggleSwitch fieldName="enabled" configPath="yp" {...extraProps}/>
          <ToggleSwitch fieldName="nsfw" configPath="instanceDetails" {...extraProps} />
        </Form>
      </div>      
    </div>
  ); 
}


