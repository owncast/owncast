import React, { useState, useContext, useEffect } from 'react';
import { Collapse, Typography } from 'antd';
import { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD, TEXTFIELD_TYPE_URL } from './TextField';
import { TextFieldWithSubmit } from './TextFieldWithSubmit';
import { ServerStatusContext } from '../../utils/server-status-context';
import { AlertMessageContext } from '../../utils/alert-message-context';
import {
  TEXTFIELD_PROPS_FFMPEG,
  TEXTFIELD_PROPS_RTMP_PORT,
  TEXTFIELD_PROPS_SOCKET_HOST_OVERRIDE,
  TEXTFIELD_PROPS_ADMIN_PASSWORD,
  TEXTFIELD_PROPS_WEB_PORT,
  TEXTFIELD_PROPS_VIDEO_SERVING_ENDPOINT,
} from '../../utils/config-constants';
import { UpdateArgs } from '../../types/config-section';
import { ResetYP } from './ResetYP';

const { Panel } = Collapse;

// eslint-disable-next-line react/function-component-definition
export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { setMessage } = useContext(AlertMessageContext);

  const { serverConfig } = serverStatusData || {};

  const {
    adminPassword,
    ffmpegPath,
    rtmpServerPort,
    webServerPort,
    yp,
    socketHostOverride,
    videoServingEndpoint,
  } = serverConfig;

  useEffect(() => {
    setFormDataValues({
      adminPassword,
      ffmpegPath,
      rtmpServerPort,
      webServerPort,
      socketHostOverride,
      videoServingEndpoint,
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
      'Changing your password will log you out of the admin. You may want to refresh the page to force yourself to log back in if not prompted.',
    );
  };

  const showFfmpegChangeMessage = () => {
    if (serverStatusData.online) {
      setMessage('The updated ffmpeg path will be used when starting your next live stream.');
    }
  };

  return (
    <div className="edit-server-details-container">
      <div className="field-container field-streamkey-container">
        <div className="left-side">
          <TextFieldWithSubmit
            fieldName="adminPassword"
            {...TEXTFIELD_PROPS_ADMIN_PASSWORD}
            value={formDataValues.adminPassword}
            initialValue={adminPassword}
            type={TEXTFIELD_TYPE_PASSWORD}
            onChange={handleFieldChange}
            onSubmit={showStreamKeyChangeMessage}
          />
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

          <TextFieldWithSubmit
            fieldName="videoServingEndpoint"
            {...TEXTFIELD_PROPS_VIDEO_SERVING_ENDPOINT}
            value={formDataValues.videoServingEndpoint}
            initialValue={videoServingEndpoint || ''}
            type={TEXTFIELD_TYPE_URL}
            onChange={handleFieldChange}
          />
          {yp.enabled && <ResetYP />}
        </Panel>
      </Collapse>
    </div>
  );
}
