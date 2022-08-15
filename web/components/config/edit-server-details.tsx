import React, { useState, useContext, useEffect } from 'react';
import { Button, Tooltip, Collapse, Typography } from 'antd';
import { CopyOutlined, RedoOutlined } from '@ant-design/icons';

import {
  TEXTFIELD_TYPE_NUMBER,
  TEXTFIELD_TYPE_PASSWORD,
  TEXTFIELD_TYPE_URL,
} from './form-textfield';
import TextFieldWithSubmit from './form-textfield-with-submit';

import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';

import {
  TEXTFIELD_PROPS_FFMPEG,
  TEXTFIELD_PROPS_RTMP_PORT,
  TEXTFIELD_PROPS_SOCKET_HOST_OVERRIDE,
  TEXTFIELD_PROPS_STREAM_KEY,
  TEXTFIELD_PROPS_WEB_PORT,
} from '../../utils/config-constants';

import { UpdateArgs } from '../../types/config-section';
import ResetYP from './reset-yp';

const { Panel } = Collapse;

export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { setMessage } = useContext(AlertMessageContext);

  const { serverConfig } = serverStatusData || {};

  const { streamKey, ffmpegPath, rtmpServerPort, webServerPort, yp, socketHostOverride } =
    serverConfig;

  const [copyIsVisible, setCopyVisible] = useState(false);

  const COPY_TOOLTIP_TIMEOUT = 3000;

  useEffect(() => {
    setFormDataValues({
      streamKey,
      ffmpegPath,
      rtmpServerPort,
      webServerPort,
      socketHostOverride,
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
      <Collapse className="advanced-settings">
        <Panel header="Advanced Settings" key="1">
          <Typography.Paragraph>
            If you have a CDN in front of your entire Owncast instance, specify your origin server
            here for the websocket to connect to. Most people will never need to set this.
          </Typography.Paragraph>
          <TextFieldWithSubmit
            fieldName="socketHostOverride"
            {...TEXTFIELD_PROPS_SOCKET_HOST_OVERRIDE}
            value={formDataValues.socketHostOverride}
            initialValue={socketHostOverride || ''}
            type={TEXTFIELD_TYPE_URL}
            onChange={handleFieldChange}
          />
          {yp.enabled && <ResetYP />}
        </Panel>
      </Collapse>
    </div>
  );
}
