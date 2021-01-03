import React, { useContext, useEffect } from 'react';
import { Typography, Form } from 'antd';

import TextField, { TEXTFIELD_TYPE_TEXTAREA } from './components/config/form-textfield';

import EditInstanceTags from './components/config/edit-tags';
import EditDirectoryDetails from './components/config/edit-directory';

import { ServerStatusContext } from '../utils/server-status-context';
import { TEXTFIELD_DEFAULTS, postConfigUpdateToAPI } from './components/config/constants';

const { Title } = Typography;

export default function PublicFacingDetails() {
  const [form] = Form.useForm();

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { instanceDetails, yp } = serverConfig;
  const { instanceDetails: instanceDetailsDefaults, yp: ypDefaults } = TEXTFIELD_DEFAULTS;

  const initialValues = {
    ...instanceDetails,
    ...yp,
  };

  const defaultFields = {
    ...instanceDetailsDefaults,
    ...ypDefaults,
  };

  useEffect(() => {
    form.setFieldsValue(initialValues);
  }, [instanceDetails]);

  const handleResetValue = (fieldName: string) => {
    const defaultValue = defaultFields[fieldName] && defaultFields[fieldName].defaultValue || '';

    form.setFieldsValue({ [fieldName]: initialValues[fieldName] || defaultValue });
  }

  // if instanceUrl is empty, we should also turn OFF the `enabled` field of directory.
  const handleSubmitInstanceUrl = () => {
    if (form.getFieldValue('instanceUrl') === '') {
      if (yp.enabled === true) {
        const { apiPath } = TEXTFIELD_DEFAULTS.yp.enabled;
        postConfigUpdateToAPI({
          apiPath,
          data: { value: false },
        });
      }
    }
  }

  const extraProps = {
    handleResetValue,
    initialValues,
    configPath: 'instanceDetails',
  };

  return (
    <div className="config-public-details-form">
      <Title level={2}>Edit your public facing instance details</Title>

      <div className="config-public-details-container">
        <div className="text-fields">
          <Form
            form={form}
            layout="vertical"
          >
            <TextField
              fieldName="instanceUrl"
              {...extraProps}
              configPath="yp"
              onSubmit={handleSubmitInstanceUrl}
            />

            <TextField fieldName="title" {...extraProps} />
            <TextField fieldName="streamTitle" {...extraProps} />
            <TextField fieldName="name" {...extraProps} />
            <TextField fieldName="summary" type={TEXTFIELD_TYPE_TEXTAREA} {...extraProps} />
            <TextField fieldName="logo" {...extraProps} />
          </Form>
        </div>
        <div className="misc-fields">
          {/* add social handles comp
          <br/>
          add tags comp */}
          <EditInstanceTags />
          <br/><br/>
          <EditDirectoryDetails />          
        </div>
      </div>      
    </div>
  ); 
}


