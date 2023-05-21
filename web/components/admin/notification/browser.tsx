import { Button, Typography } from 'antd';
import React, { useState, useContext, useEffect } from 'react';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { TextField, TEXTFIELD_TYPE_TEXTAREA } from '../TextField';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  BROWSER_PUSH_CONFIG_FIELDS,
} from '../../../utils/config-constants';
import { ToggleSwitch } from '../ToggleSwitch';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../../utils/input-statuses';
import { UpdateArgs } from '../../../types/config-section';
import { FormStatusIndicator } from '../FormStatusIndicator';

const { Title } = Typography;

export const BrowserNotify = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { notifications } = serverConfig || {};
  const { browser } = notifications || {};

  const { enabled, goLiveMessage } = browser || {};

  const [formDataValues, setFormDataValues] = useState<any>({});
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  const [enableSaveButton, setEnableSaveButton] = useState<boolean>(false);

  useEffect(() => {
    setFormDataValues({
      enabled,
      goLiveMessage,
    });
  }, [notifications, browser]);

  const canSave = (): boolean => true;

  // update individual values in state
  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    console.log(fieldName, value);
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });

    setEnableSaveButton(canSave());
  };

  // toggle switch.
  const handleSwitchChange = (switchEnabled: boolean) => {
    // setShouldDisplayForm(storageEnabled);
    handleFieldChange({ fieldName: 'enabled', value: switchEnabled });
  };

  let resetTimer = null;
  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  const save = async () => {
    const postValue = formDataValues;

    await postConfigUpdateToAPI({
      apiPath: '/notifications/browser',
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'browser',
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
      <Title>Browser Alerts</Title>
      <p className="description reduced-margins">
        Viewers can opt into being notified when you go live with their browser.
      </p>
      <p className="description reduced-margins">Not all browsers support this.</p>
      <ToggleSwitch
        apiPath=""
        fieldName="enabled"
        label="Enable browser notifications"
        onChange={handleSwitchChange}
        checked={formDataValues.enabled}
      />
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...BROWSER_PUSH_CONFIG_FIELDS.goLiveMessage}
          required
          type={TEXTFIELD_TYPE_TEXTAREA}
          value={formDataValues.goLiveMessage}
          onChange={handleFieldChange}
        />
      </div>
      <Button
        type="primary"
        style={{
          display: enableSaveButton ? 'inline-block' : 'none',
          position: 'relative',
          marginLeft: 'auto',
          right: '0',
          marginTop: '20px',
        }}
        onClick={save}
      >
        Save
      </Button>
      <FormStatusIndicator status={submitStatus} />
    </>
  );
};
