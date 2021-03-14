import React, { useState, useContext, useEffect } from 'react';
import { Button, Tooltip, Collapse } from 'antd';
import { CopyOutlined, RedoOutlined } from '@ant-design/icons';

import { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD } from './form-textfield';
import TextFieldWithSubmit from './form-textfield-with-submit';
import ToggleSwitch from './form-toggleswitch';

import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';

import {
  TEXTFIELD_PROPS_FFMPEG,
  TEXTFIELD_PROPS_RTMP_PORT,
  TEXTFIELD_PROPS_STREAM_KEY,
  TEXTFIELD_PROPS_WEB_PORT,
  FIELD_PROPS_DISABLE_CHAT,
} from '../../utils/config-constants';

import { UpdateArgs } from '../../types/config-section';
import ResetYP from './reset-yp';

const { Panel } = Collapse;

export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { setMessage } = useContext(AlertMessageContext);

  const { serverConfig } = serverStatusData || {};

  const { chatDisabled, streamKey, ffmpegPath, rtmpServerPort, webServerPort, yp } = serverConfig;

  const [copyIsVisible, setCopyVisible] = useState(false);

  const COPY_TOOLTIP_TIMEOUT = 3000;

  useEffect(() => {
    setFormDataValues({
      streamKey,
      ffmpegPath,
      rtmpServerPort,
      webServerPort,
      chatDisabled,
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

  const showConfigurationRestartMessage = () => {
    setMessage('Updating server settings requires a restart of your Owncast server.');
  };

  const showStreamKeyChangeMessage = () => {
    setMessage(
      'Changing your stream key will log you out of the admin and block you from streaming until you change the key in your broadcasting software.',
    );
  };

  const showFfmpegChangeMessage = () => {
    if (serverStatusData.online) {
      setMessage('The updated ffmpeg path will be used when starting your next live stream.');
    }
  };

  function generateStreamKey() {
    let key = '';
    for (let i = 0; i < 3; i += 1) {
      key += Math.random().toString(36).substring(2);
    }

    handleFieldChange({ fieldName: 'streamKey', value: key });
  }

  function copyStreamKey() {
    navigator.clipboard.writeText(formDataValues.streamKey).then(() => {
      setCopyVisible(true);
      setTimeout(() => setCopyVisible(false), COPY_TOOLTIP_TIMEOUT);
    });
  }

  function handleChatDisableChange(disabled: boolean) {
    handleFieldChange({ fieldName: 'chatDisabled', value: disabled });
  }

  return (
    <div className="edit-server-details-container">
      <div className="field-container field-streamkey-container">
        <div className="left-side">
          <TextFieldWithSubmit
            fieldName="streamKey"
            {...TEXTFIELD_PROPS_STREAM_KEY}
            value={formDataValues.streamKey}
            initialValue={streamKey}
            type={TEXTFIELD_TYPE_PASSWORD}
            onChange={handleFieldChange}
            onSubmit={showStreamKeyChangeMessage}
          />
          <div className="streamkey-actions">
            <Tooltip title="Generate a stream key">
              <Button icon={<RedoOutlined />} size="small" onClick={generateStreamKey} />
            </Tooltip>

            <Tooltip
              className="copy-tooltip"
              title={copyIsVisible ? 'Copied!' : 'Copy to clipboard'}
            >
              <Button icon={<CopyOutlined />} size="small" onClick={copyStreamKey} />
            </Tooltip>
          </div>
        </div>
      </div>
      <TextFieldWithSubmit
        fieldName="ffmpegPath"
        {...TEXTFIELD_PROPS_FFMPEG}
        value={formDataValues.ffmpegPath}
        initialValue={ffmpegPath}
        onChange={handleFieldChange}
        onSubmit={showFfmpegChangeMessage}
      />
      <TextFieldWithSubmit
        fieldName="webServerPort"
        {...TEXTFIELD_PROPS_WEB_PORT}
        value={formDataValues.webServerPort}
        initialValue={webServerPort}
        type={TEXTFIELD_TYPE_NUMBER}
        onChange={handleFieldChange}
        onSubmit={showConfigurationRestartMessage}
      />
      <TextFieldWithSubmit
        fieldName="rtmpServerPort"
        {...TEXTFIELD_PROPS_RTMP_PORT}
        value={formDataValues.rtmpServerPort}
        initialValue={rtmpServerPort}
        type={TEXTFIELD_TYPE_NUMBER}
        onChange={handleFieldChange}
        onSubmit={showConfigurationRestartMessage}
      />
      <ToggleSwitch
        fieldName="chatDisabled"
        {...FIELD_PROPS_DISABLE_CHAT}
        checked={formDataValues.chatDisabled}
        onChange={handleChatDisableChange}
      />

      {yp.enabled && (
        <Collapse className="advanced-settings">
          <Panel header="Advanced Settings" key="1">
            <ResetYP />
          </Panel>
        </Collapse>
      )}
    </div>
  );
}
