import { Button, Typography } from 'antd';
import React, { useState, useContext, useEffect } from 'react';
import { ServerStatusContext } from '../../../utils/server-status-context';
import TextField, { TEXTFIELD_TYPE_PASSWORD } from '../form-textfield';
import FormStatusIndicator from '../form-status-indicator';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  TWITTER_CONFIG_FIELDS,
} from '../../../utils/config-constants';
import ToggleSwitch from '../form-toggleswitch';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../../utils/input-statuses';
import { UpdateArgs } from '../../../types/config-section';
import { TEXTFIELD_TYPE_TEXT } from '../form-textfield-with-submit';

const { Title } = Typography;

export default function ConfigNotify() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { notifications } = serverConfig || {};
  const { twitter } = notifications || {};

  const [formDataValues, setFormDataValues] = useState<any>({});
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  const [enableSaveButton, setEnableSaveButton] = useState<boolean>(false);

  useEffect(() => {
    const {
      enabled,
      apiKey,
      apiSecret,
      accessToken,
      accessTokenSecret,
      bearerToken,
      goLiveMessage,
    } = twitter || {};
    setFormDataValues({
      enabled,
      apiKey,
      apiSecret,
      accessToken,
      accessTokenSecret,
      bearerToken,
      goLiveMessage,
    });
  }, [twitter]);

  const canSave = (): boolean => {
    const {
      enabled,
      apiKey,
      apiSecret,
      accessToken,
      accessTokenSecret,
      bearerToken,
      goLiveMessage,
    } = formDataValues;

    return (
      enabled &&
      !!apiKey &&
      !!apiSecret &&
      !!accessToken &&
      !!accessTokenSecret &&
      !!bearerToken &&
      !!goLiveMessage
    );
  };

  // update individual values in state
  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });

    setEnableSaveButton(canSave());
  };

  // toggle switch.
  const handleSwitchChange = (switchEnabled: boolean) => {
    const previouslySaved = formDataValues.enabled;

    handleFieldChange({ fieldName: 'enabled', value: switchEnabled });

    return switchEnabled !== previouslySaved;
  };

  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
    setEnableSaveButton(false);
  };

  const save = async () => {
    const postValue = formDataValues;

    await postConfigUpdateToAPI({
      apiPath: '/notifications/twitter',
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'twitter',
          value: postValue,
          path: 'notifications',
        });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Updated.'));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  return (
    <>
      <Title>Twitter</Title>
      <p className="description reduced-margins">
        Let your Twitter followers know each time you go live.
      </p>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <p className="description reduced-margins">
          <a href="https://owncast.online/docs/notifications" target="_blank" rel="noreferrer">
            Read how to configure your Twitter account
          </a>{' '}
          to support posting from Owncast.
        </p>
        <p className="description reduced-margins">
          <a
            href="https://developer.twitter.com/en/portal/dashboard"
            target="_blank"
            rel="noreferrer"
          >
            And then get your Twitter developer credentials
          </a>{' '}
          to fill in below.
        </p>
      </div>

      <ToggleSwitch
        apiPath=""
        fieldName="enabled"
        label="Enable Twitter"
        onChange={handleSwitchChange}
        checked={formDataValues.enabled}
      />
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.apiKey}
          required
          value={formDataValues.apiKey}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.apiSecret}
          type={TEXTFIELD_TYPE_PASSWORD}
          required
          value={formDataValues.apiSecret}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.accessToken}
          required
          value={formDataValues.accessToken}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.accessTokenSecret}
          type={TEXTFIELD_TYPE_PASSWORD}
          required
          value={formDataValues.accessTokenSecret}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.bearerToken}
          required
          value={formDataValues.bearerToken}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...TWITTER_CONFIG_FIELDS.goLiveMessage}
          type={TEXTFIELD_TYPE_TEXT}
          required
          value={formDataValues.goLiveMessage}
          onChange={handleFieldChange}
        />
      </div>
      <Button
        type="primary"
        onClick={save}
        style={{
          display: enableSaveButton ? 'inline-block' : 'none',
          position: 'relative',
          marginLeft: 'auto',
          right: '0',
          marginTop: '20px',
        }}
      >
        Save
      </Button>
      <FormStatusIndicator status={submitStatus} />
    </>
  );
}
