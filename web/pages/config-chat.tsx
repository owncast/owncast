import { Typography } from 'antd';
import React, { useContext, useEffect, useState } from 'react';
import { TEXTFIELD_TYPE_TEXTAREA } from '../components/config/form-textfield';
import TextFieldWithSubmit from '../components/config/form-textfield-with-submit';
import ToggleSwitch from '../components/config/form-toggleswitch';
import { UpdateArgs } from '../types/config-section';
import {
  FIELD_PROPS_DISABLE_CHAT,
  TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES,
  TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE,
} from '../utils/config-constants';
import { ServerStatusContext } from '../utils/server-status-context';

export default function ConfigChat() {
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);

  const { serverConfig } = serverStatusData || {};
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

  function handleChatForbiddenUsernamesChange(args: UpdateArgs) {
    const updatedForbiddenUsernameList = args.value.split(',');
    handleFieldChange({ fieldName: args.fieldName, value: updatedForbiddenUsernameList });
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
          fieldName="forbiddenUsernames"
          {...TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES}
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.forbiddenUsernames}
          onChange={handleChatForbiddenUsernamesChange}
        />
        <TextFieldWithSubmit
          fieldName="welcomeMessage"
          {...TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE}
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.welcomeMessage}
          onChange={handleFieldChange}
        />
      </div>
    </div>
  );
}
