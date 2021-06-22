import React, { useState, useContext, useEffect } from 'react';
import { Typography } from 'antd';
import {
  FIELD_PROPS_DISABLE_CHAT,
  TEXTFIELD_PROPS_CHAT_USERNAME_BLOCKLIST,
} from '../utils/config-constants';
import { ServerStatusContext } from '../utils/server-status-context';
import ToggleSwitch from '../components/config/form-toggleswitch';
import { UpdateArgs } from '../types/config-section';
import { TEXTFIELD_TYPE_TEXTAREA } from '../components/config/form-textfield';
import TextFieldWithSubmit from '../components/config/form-textfield-with-submit';

export default function ConfigChat() {
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);

  const { serverConfig } = serverStatusData || {};
  const { chatDisabled } = serverConfig;
  const { usernameBlocklist } = serverConfig;

  function handleChatDisableChange(disabled: boolean) {
    handleFieldChange({ fieldName: 'chatDisabled', value: disabled });
  }

  function handleChatUsernameBlockListChange(args: UpdateArgs) {
    handleFieldChange({ fieldName: 'usernameBlocklist', value: args.value });
  }

  useEffect(() => {
    setFormDataValues({
      chatDisabled,
      usernameBlocklist,
    });
  }, [serverConfig]);

  if (!formDataValues) {
    return null;
  }

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };
  
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
          fieldName="usernameBlocklist"
          {...TEXTFIELD_PROPS_CHAT_USERNAME_BLOCKLIST}
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.usernameBlocklist}
          initialValue={usernameBlocklist}
          onChange={handleChatUsernameBlockListChange}
        />
      </div>
    </div>
  );
}
