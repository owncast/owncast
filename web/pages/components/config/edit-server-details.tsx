import React, { useState, useContext, useEffect } from 'react';
import TextField, { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD } from './form-textfield';

import { ServerStatusContext } from '../../../utils/server-status-context';
import { TEXTFIELD_PROPS_FFMPEG, TEXTFIELD_PROPS_RTMP_PORT, TEXTFIELD_PROPS_STREAM_KEY, TEXTFIELD_PROPS_WEB_PORT,  } from './constants';

import configStyles from '../../../styles/config-pages.module.scss';

export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { streamKey, ffmpegPath, rtmpServerPort, webServerPort } = serverConfig;

  useEffect(() => {
    setFormDataValues({
      streamKey, ffmpegPath, rtmpServerPort, webServerPort
    });
  }, [serverConfig]);

  if (!formDataValues) {
    return null;
  }

  const handleFieldChange = (fieldName: string, value: string) => {
    setFormDataValues({
      ...formDataValues,
      [fieldName]: value,
    });
  }

  return (  
    <div className={configStyles.publicDetailsContainer}>
      <div className={configStyles.textFieldsSection}>
        <TextField
          fieldName="streamKey"
          {...TEXTFIELD_PROPS_STREAM_KEY}
          value={formDataValues.streamKey}
          initialValue={streamKey}
          type={TEXTFIELD_TYPE_PASSWORD}
          onChange={handleFieldChange}
        />
        
        <TextField
          fieldName="ffmpegPath"
          {...TEXTFIELD_PROPS_FFMPEG}
          value={formDataValues.ffmpegPath}
          initialValue={ffmpegPath}
          onChange={handleFieldChange}
        />
        <TextField
          fieldName="webServerPort"
          {...TEXTFIELD_PROPS_WEB_PORT}
          value={formDataValues.webServerPort}
          initialValue={webServerPort}
          type={TEXTFIELD_TYPE_NUMBER}
          onChange={handleFieldChange}
        />
        <TextField
          fieldName="rtmpServerPort"
          {...TEXTFIELD_PROPS_RTMP_PORT}
          value={formDataValues.rtmpServerPort}
          initialValue={rtmpServerPort}
          type={TEXTFIELD_TYPE_NUMBER}
          onChange={handleFieldChange}
        />
      </div>
    </div>      
  ); 
}


