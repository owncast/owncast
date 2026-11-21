import { Col, Row, Typography } from 'antd';
import { ReactElement, useContext, useEffect, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import { TEXTFIELD_TYPE_TEXTAREA } from '../../components/admin/TextField';
import { TextFieldWithSubmit } from '../../components/admin/TextFieldWithSubmit';
import { ToggleSwitch } from '../../components/admin/ToggleSwitch';
import { EditValueArray } from '../../components/admin/EditValueArray';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';

import { UpdateArgs } from '../../types/config-section';
import {
  API_CHAT_FORBIDDEN_USERNAMES,
  API_CHAT_SUGGESTED_USERNAMES,
  FIELD_PROPS_CHAT_JOIN_MESSAGES_ENABLED,
  FIELD_PROPS_ENABLE_CHAT_SLUR_FILTER,
  CHAT_ESTABLISHED_USER_MODE,
  FIELD_PROPS_DISABLE_CHAT,
  FIELD_PROPS_CHAT_REQUIRE_AUTHENTICATION,
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE,
  FIELD_PROPS_ENABLE_SPAM_PROTECTION,
} from '../../utils/config-constants';
import { ServerStatusContext } from '../../utils/server-status-context';

import { AdminLayout } from '../../components/layouts/AdminLayout';
import { Localization } from '../../types/localization';

export default function ConfigChat() {
  const { t } = useTranslation();
  const { Title } = Typography;
  const [formDataValues, setFormDataValues] = useState(null);
  const [forbiddenUsernameSaveState, setForbiddenUsernameSaveState] = useState<StatusState>(null);
  const [suggestedUsernameSaveState, setSuggestedUsernameSaveState] = useState<StatusState>(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const {
    chatDisabled,
    chatJoinMessagesEnabled,
    forbiddenUsernames,
    instanceDetails,
    suggestedUsernames,
    chatEstablishedUserMode,
    chatSpamProtectionEnabled,
    chatSlurFilterEnabled,
    chatRequireAuthentication,
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

  function handleEstablishedUserModeChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'chatEstablishedUserMode', value: enabled });
  }

  function handleChatSpamProtectionChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'chatSpamProtectionEnabled', value: enabled });
  }

  function handleChatSlurFilterChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'chatSlurFilterEnabled', value: enabled });
  }

  function handleChatRequireAuthenticationChange(enabled: boolean) {
    handleFieldChange({ fieldName: 'chatRequireAuthentication', value: enabled });
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
        setForbiddenUsernameSaveState(createInputStatus(STATUS_SUCCESS));
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
        setSuggestedUsernameSaveState(createInputStatus(STATUS_SUCCESS));
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
      return createInputStatus('success', t(Localization.Admin.ConfigChat.defaultUsernamesNone));
    if (length > 0 && length < 10)
      return createInputStatus(
        'warning',
        t(Localization.Admin.ConfigChat.defaultUsernamesMinNotReached),
      );
    return null;
  }

  useEffect(() => {
    setFormDataValues({
      chatDisabled,
      chatJoinMessagesEnabled,
      forbiddenUsernames,
      suggestedUsernames,
      welcomeMessage,
      chatEstablishedUserMode,
      chatSpamProtectionEnabled,
      chatSlurFilterEnabled,
      chatRequireAuthentication,
    });
  }, [serverConfig]);

  if (!formDataValues) {
    return null;
  }

  return (
    <div className="config-server-details-form">
      <Title>{t(Localization.Admin.ConfigChat.title)}</Title>
      <Row gutter={[45, 16]}>
        <Col md={24} lg={12}>
          <div className="form-module">
            <ToggleSwitch
              fieldName="chatDisabled"
              {...FIELD_PROPS_DISABLE_CHAT}
              label={t(Localization.Admin.ConfigChat.chatLabel)}
              tip={t(Localization.Admin.ConfigChat.chatTip)}
              checked={!formDataValues.chatDisabled}
              reversed
              onChange={handleChatDisableChange}
            />
            <ToggleSwitch
              fieldName="chatJoinMessagesEnabled"
              {...FIELD_PROPS_CHAT_JOIN_MESSAGES_ENABLED}
              label={t(Localization.Admin.ConfigChat.joinMessagesLabel)}
              tip={t(Localization.Admin.ConfigChat.joinMessagesTip)}
              checked={formDataValues.chatJoinMessagesEnabled}
              onChange={handleChatJoinMessagesEnabledChange}
            />
            <TextFieldWithSubmit
              fieldName="welcomeMessage"
              {...TEXTFIELD_PROPS_SERVER_WELCOME_MESSAGE}
              label={t(Localization.Admin.ConfigChat.welcomeMessageLabel)}
              tip={t(Localization.Admin.ConfigChat.welcomeMessageTip)}
              type={TEXTFIELD_TYPE_TEXTAREA}
              value={formDataValues.welcomeMessage}
              initialValue={welcomeMessage}
              onChange={handleFieldChange}
            />
            <br />
            <br />
            <EditValueArray
              title={t(Localization.Admin.ConfigChat.forbiddenUsernamesTitle)}
              placeholder={t(Localization.Admin.ConfigChat.usernamePlaceholder)}
              description={t(Localization.Admin.ConfigChat.forbiddenUsernamesTip)}
              values={formDataValues.forbiddenUsernames}
              handleDeleteIndex={handleDeleteForbiddenUsernameIndex}
              handleCreateString={handleCreateForbiddenUsername}
              submitStatus={forbiddenUsernameSaveState}
            />
            <br />
            <br />
            <EditValueArray
              title={t(Localization.Admin.ConfigChat.defaultUsernamesTitle)}
              placeholder={t(Localization.Admin.ConfigChat.usernamePlaceholder)}
              description={t(Localization.Admin.ConfigChat.defaultUsernamesTip)}
              values={formDataValues.suggestedUsernames}
              handleDeleteIndex={handleDeleteSuggestedUsernameIndex}
              handleCreateString={handleCreateSuggestedUsername}
              submitStatus={suggestedUsernameSaveState}
              continuousStatusMessage={getSuggestedUsernamesLimitWarning(
                formDataValues.suggestedUsernames.length,
              )}
            />
          </div>
        </Col>
        <Col md={24} lg={12}>
          <div className="form-module">
            <ToggleSwitch
              fieldName="chatSpamProtectionEnabled"
              {...FIELD_PROPS_ENABLE_SPAM_PROTECTION}
              label={t(Localization.Admin.ConfigChat.spamProtectionLabel)}
              tip={t(Localization.Admin.ConfigChat.spamProtectionTip)}
              checked={formDataValues.chatSpamProtectionEnabled}
              onChange={handleChatSpamProtectionChange}
            />
            <ToggleSwitch
              fieldName="chatEstablishedUserMode"
              {...CHAT_ESTABLISHED_USER_MODE}
              label={t(Localization.Admin.ConfigChat.establishedUsersLabel)}
              tip={t(Localization.Admin.ConfigChat.establishedUsersTip)}
              checked={formDataValues.chatEstablishedUserMode}
              onChange={handleEstablishedUserModeChange}
            />
            <ToggleSwitch
              fieldName="chatSlurFilterEnabled"
              {...FIELD_PROPS_ENABLE_CHAT_SLUR_FILTER}
              label={t(Localization.Admin.ConfigChat.slurFilterLabel)}
              tip={t(Localization.Admin.ConfigChat.slurFilterTip)}
              checked={formDataValues.chatSlurFilterEnabled}
              onChange={handleChatSlurFilterChange}
            />
            <ToggleSwitch
              fieldName="chatRequireAuthentication"
              {...FIELD_PROPS_CHAT_REQUIRE_AUTHENTICATION}
              label={t(Localization.Admin.ConfigChat.requireAuthenticationLabel)}
              tip={t(Localization.Admin.ConfigChat.requireAuthenticationTip)}
              checked={formDataValues.chatRequireAuthentication}
              onChange={handleChatRequireAuthenticationChange}
            />
          </div>
        </Col>
      </Row>
    </div>
  );
}

ConfigChat.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
