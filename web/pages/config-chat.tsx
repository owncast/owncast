import { Typography } from 'antd';
import React, { useContext, useEffect, useState } from 'react';
import { TEXTFIELD_TYPE_TEXTAREA } from '../components/config/form-textfield';
import TextFieldWithSubmit from '../components/config/form-textfield-with-submit';
import ToggleSwitch from '../components/config/form-toggleswitch';
import EditValueArray from '../components/config/edit-string-array';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../utils/input-statuses';

import { UpdateArgs } from '../types/config-section';
import {
  API_CHAT_FORBIDDEN_USERNAMES,
  API_CHAT_SUGGESTED_USERNAMES,
  FIELD_PROPS_CHAT_JOIN_MESSAGES_ENABLED,
  FIELD_PROPS_DISABLE_CHAT,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  TEXTFIELD_PROPS_CHAT_FORBIDDEN_USERNAMES,
  TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES,
  TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE,
} from '../utils/config-constants';
import { ServerStatusContext } from '../utils/server-status-context';

export default function ConfigChat() {
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const [forbiddenUsernameSaveState, setForbiddenUsernameSaveState] = useState(null);
  const [suggestedUsernameSaveState, setSuggestedUsernameSaveState] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const {
    chatDisabled,
    chatJoinMessagesEnabled,
    forbiddenUsernames,
    instanceDetails,
    suggestedUsernames,
  } = serverConfig;
  const { welcomeMessage } = instanceDetails;

  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  };

  function handleChatDisableChange(disabled: boolean) {
    handleFieldChange({ fieldName: 'chatDisabled', value: !disabled });
  }

  function handleChatJoinMessagesEnabledChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'chatJoinMessagesEnabled', value: enabled });
  }

  function resetForbiddenUsernameState() {
    setForbiddenUsernameSaveState(null);
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
        setForbiddenUsernameSaveState(STATUS_SUCCESS);
        setTimeout(resetForbiddenUsernameState, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setForbiddenUsernameSaveState(createInputStatus(STATUS_ERROR, message));
        setTimeout(resetForbiddenUsernameState, RESET_TIMEOUT);
      },
    });
  }

  function handleDeleteForbiddenUsernameIndex(index: number) {
    formDataValues.forbiddenUsernames.splice(index, 1);
    saveForbiddenUsernames();
  }

  function handleCreateForbiddenUsername(tag: string) {
    formDataValues.forbiddenUsernames.push(tag);
    handleFieldChange({
      fieldName: 'forbiddenUsernames',
      value: formDataValues.forbiddenUsernames,
    });
    saveForbiddenUsernames();
  }

  function resetSuggestedUsernameState() {
    setSuggestedUsernameSaveState(null);
  }

  function saveSuggestedUsernames() {
    postConfigUpdateToAPI({
      apiPath: API_CHAT_SUGGESTED_USERNAMES,
      data: { value: formDataValues.suggestedUsernames },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'suggestedUsernames',
          value: formDataValues.suggestedUsernames,
        });
        setSuggestedUsernameSaveState(STATUS_SUCCESS);
        setTimeout(resetSuggestedUsernameState, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setForbiddenUsernameSaveState(createInputStatus(STATUS_ERROR, message));
        setTimeout(resetSuggestedUsernameState, RESET_TIMEOUT);
      },
    });
  }

  function handleDeleteSuggestedUsernameIndex(index: number) {
    formDataValues.suggestedUsernames.splice(index, 1);
    saveSuggestedUsernames();
  }

  function handleCreateSuggestedUsername(tag: string) {
    formDataValues.suggestedUsernames.push(tag);
    handleFieldChange({
      fieldName: 'suggestedUsernames',
      value: formDataValues.suggestedUsernames,
    });
    saveSuggestedUsernames();
  }

  function getSuggestedUsernamesLimitWarning(length: number): StatusState | null {
    if (length === 0)
      return createInputStatus('success', TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES.no_entries);
    if (length > 0 && length < 10)
      return createInputStatus('warning', TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES.min_not_reached);
    return null;
  }

  useEffect(() => {
    setFormDataValues({
      chatDisabled,
      chatJoinMessagesEnabled,
      forbiddenUsernames,
      suggestedUsernames,
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
          checked={!formDataValues.chatDisabled}
          reversed
          onChange={handleChatDisableChange}
        />
        <ToggleSwitch
          fieldName="chatJoinMessagesEnabled"
          {...FIELD_PROPS_CHAT_JOIN_MESSAGES_ENABLED}
          checked={formDataValues.chatJoinMessagesEnabled}
          onChange={handleChatJoinMessagesEnabledChange}
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
          handleDeleteIndex={handleDeleteForbiddenUsernameIndex}
          handleCreateString={handleCreateForbiddenUsername}
          submitStatus={createInputStatus(forbiddenUsernameSaveState)}
        />
        <br />
        <br />
        <EditValueArray
          title={TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES.label}
          placeholder={TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES.placeholder}
          description={TEXTFIELD_PROPS_CHAT_SUGGESTED_USERNAMES.tip}
          values={formDataValues.suggestedUsernames}
          handleDeleteIndex={handleDeleteSuggestedUsernameIndex}
          handleCreateString={handleCreateSuggestedUsername}
          submitStatus={createInputStatus(suggestedUsernameSaveState)}
          continuousStatusMessage={getSuggestedUsernamesLimitWarning(
            formDataValues.suggestedUsernames.length,
          )}
        />
      </div>
    </div>
  );
}
