import React, { useState, useContext, useEffect } from 'react';
import TextFieldWithSubmit, {
  TEXTFIELD_TYPE_TEXTAREA,
  TEXTFIELD_TYPE_URL,
} from './form-textfield-with-submit';

import { ServerStatusContext } from '../../utils/server-status-context';
import {
  postConfigUpdateToAPI,
  TEXTFIELD_PROPS_INSTANCE_URL,
  TEXTFIELD_PROPS_SERVER_NAME,
  TEXTFIELD_PROPS_SERVER_SUMMARY,
  TEXTFIELD_PROPS_LOGO,
  API_YP_SWITCH,
} from '../../utils/config-constants';

import { UpdateArgs } from '../../types/config-section';

export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { instanceDetails, yp } = serverConfig;

  useEffect(() => {
    setFormDataValues({
      ...instanceDetails,
      ...yp,
    });
  }, [instanceDetails, yp]);

  if (!formDataValues) {
    return null;
  }

  // if instanceUrl is empty, we should also turn OFF the `enabled` field of directory.
  const handleSubmitInstanceUrl = () => {
    if (formDataValues.instanceUrl === '') {
      if (yp.enabled === true) {
        postConfigUpdateToAPI({
          apiPath: API_YP_SWITCH,
          data: { value: false },
        });
      }
    }
  };

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  return (
    <div className={`publicDetailsContainer`}>
      <div className={`textFieldsSection`}>
        <TextFieldWithSubmit
          fieldName="instanceUrl"
          {...TEXTFIELD_PROPS_INSTANCE_URL}
          value={formDataValues.instanceUrl}
          initialValue={yp.instanceUrl}
          type={TEXTFIELD_TYPE_URL}
          onChange={handleFieldChange}
          onSubmit={handleSubmitInstanceUrl}
        />

        <TextFieldWithSubmit
          fieldName="name"
          {...TEXTFIELD_PROPS_SERVER_NAME}
          value={formDataValues.name}
          initialValue={instanceDetails.name}
          onChange={handleFieldChange}
        />
        <TextFieldWithSubmit
          fieldName="summary"
          {...TEXTFIELD_PROPS_SERVER_SUMMARY}
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.summary}
          initialValue={instanceDetails.summary}
          onChange={handleFieldChange}
        />
        <TextFieldWithSubmit
          fieldName="logo"
          {...TEXTFIELD_PROPS_LOGO}
          value={formDataValues.logo}
          initialValue={instanceDetails.logo}
          onChange={handleFieldChange}
        />
      </div>
    </div>
  );
}
