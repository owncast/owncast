import { Button, Typography } from 'antd';
import React, { useState, useContext, useEffect } from 'react';
import { ServerStatusContext } from '../../../utils/server-status-context';
import { TextField } from '../TextField';
import { FormStatusIndicator } from '../FormStatusIndicator';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  DISCORD_CONFIG_FIELDS,
} from '../../../utils/config-constants';
import { ToggleSwitch } from '../ToggleSwitch';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_SUCCESS,
} from '../../../utils/input-statuses';
import { UpdateArgs } from '../../../types/config-section';

const { Title } = Typography;

export const DiscordNotify = () => {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  const { notifications } = serverConfig || {};
  const { discord } = notifications || {};

  const { enabled, webhook, goLiveMessage } = discord || {};

  const [formDataValues, setFormDataValues] = useState<any>({});
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  const [enableSaveButton, setEnableSaveButton] = useState<boolean>(false);

  useEffect(() => {
    setFormDataValues({
      enabled,
      webhook,
      goLiveMessage,
    });
  }, [notifications, discord]);

  const canSave = (): boolean => {
    if (webhook === '' || goLiveMessage === '') {
      return false;
    }

    return true;
  };

  // update individual values in state
  const handleFieldChange = ({ fieldName, value }: UpdateArgs) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });

    setEnableSaveButton(canSave());
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
      apiPath: '/notifications/discord',
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({
          fieldName: 'discord',
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

  // toggle switch.
  const handleSwitchChange = (switchEnabled: boolean) => {
    // setShouldDisplayForm(storageEnabled);
    handleFieldChange({ fieldName: 'enabled', value: switchEnabled });
  };

  return (
    <>
      <Title>Discord</Title>
      <p className="description reduced-margins">
        Let your Discord channel know each time you go live.
      </p>
      <p className="description reduced-margins">
        <a
          href="https://support.discord.com/hc/en-us/articles/228383668"
          target="_blank"
          rel="noreferrer"
        >
          Create a webhook
        </a>{' '}
        under <i>Edit Channel / Integrations</i> on your Discord channel and provide it below.
      </p>

      <ToggleSwitch
        apiPath=""
        fieldName="discordEnabled"
        label="Enable Discord"
        checked={formDataValues.enabled}
        onChange={handleSwitchChange}
      />
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...DISCORD_CONFIG_FIELDS.webhookUrl}
          required
          value={formDataValues.webhook}
          onChange={handleFieldChange}
        />
      </div>
      <div style={{ display: formDataValues.enabled ? 'block' : 'none' }}>
        <TextField
          {...DISCORD_CONFIG_FIELDS.goLiveMessage}
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
};
