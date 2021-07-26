import { Typography } from 'antd';
import React, { useContext, useEffect, useState } from 'react';
import { TEXTFIELD_TYPE_TEXTAREA } from '../components/config/form-textfield';
import TextFieldWithSubmit from '../components/config/form-textfield-with-submit';
import ToggleSwitch from '../components/config/form-toggleswitch';
import EditValueArray from '../components/config/edit-string-array';
import { createInputStatus, STATUS_ERROR, STATUS_SUCCESS } from '../utils/input-statuses';

import { UpdateArgs } from '../types/config-section';
import {
  FIELD_PROPS_DISABLE_CHAT,
  TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES,
  TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE,
  API_CHAT_FORBIDDEN_USERNAMES,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
} from '../utils/config-constants';
import { ServerStatusContext } from '../utils/server-status-context';

export default function ConfigChat() {
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const [forbiddenUsernameSaveState, setForbidenUsernameSaveState] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { chatDisabled, forbiddenUsernames } = serverConfig;
  const { instanceDetails } = serverConfig;
  const { welcomeMessage } = instanceDetails;

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  function handleChatDisableChange(disabled: boolean) {
    handleFieldChange({ fieldName: 'chatDisabled', value: disabled });
  }

  function resetForbiddenUsernameState() {
    setForbidenUsernameSaveState(null);
  }

  function saveForbiddenUsernames() {
    postConfigUpdateToAPI({
      apiPath: API_CHAT_FORBIDDEN_USERNAMES,
      data: { value: formDataValues.forbiddenUsernames },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'forbiddenUsernames',
          value: formDataValues.forbiddenUsernames,
        });
        setForbidenUsernameSaveState(STATUS_SUCCESS);
        setTimeout(resetForbiddenUsernameState, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setForbidenUsernameSaveState(createInputStatus(STATUS_ERROR, message));
        setTimeout(resetForbiddenUsernameState, RESET_TIMEOUT);
      },
    });
  }

  function handleDeleteUsername(index: number) {
    formDataValues.forbiddenUsernames.splice(index, 1);
    saveForbiddenUsernames();
  }

  function handleCreateUsername(tag: string) {
    formDataValues.forbiddenUsernames.push(tag);
    handleFieldChange({
      fieldName: 'forbiddenUsernames',
      value: formDataValues.forbiddenUsernames,
    });
    saveForbiddenUsernames();
  }

  useEffect(() => {
    setFormDataValues({
      chatDisabled,
      forbiddenUsernames,
      welcomeMessage,
    });
  }, [serverConfig]);

  if (!formDataValues) {
    return null;
  }

  return (
    <div className="config-server-details-form">
      <Title>Chat Settings</Title>
      <div className="form-module config-server-details-container">
        <ToggleSwitch
          fieldName="chatDisabled"
          {...FIELD_PROPS_DISABLE_CHAT}
          checked={formDataValues.chatDisabled}
          onChange={handleChatDisableChange}
        />
        <TextFieldWithSubmit
          fieldName="welcomeMessage"
          {...TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE}
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.welcomeMessage}
          initialValue={welcomeMessage}
          onChange={handleFieldChange}
        />
        <br />
        <br />
        <EditValueArray
          title={TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES.label}
          placeholder={TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES.placeholder}
          description={TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES.tip}
          values={formDataValues.forbiddenUsernames}
          handleDeleteIndex={handleDeleteUsername}
          handleCreateString={handleCreateUsername}
          submitStatus={createInputStatus(forbiddenUsernameSaveState)}
        />
      </div>
    </div>
  );
}
