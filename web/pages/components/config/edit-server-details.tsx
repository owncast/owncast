import React, { useState, useContext, useEffect } from 'react';
import { Button, Tooltip } from 'antd';
import { CopyOutlined, RedoOutlined } from '@ant-design/icons';

import { TEXTFIELD_TYPE_NUMBER, TEXTFIELD_TYPE_PASSWORD } from './form-textfield-nosubmit';
import TextFieldWithSubmit from './form-textfield-with-submit';

import { ServerStatusContext } from '../../../utils/server-status-context';
import { TEXTFIELD_PROPS_FFMPEG, TEXTFIELD_PROPS_RTMP_PORT, TEXTFIELD_PROPS_STREAM_KEY, TEXTFIELD_PROPS_WEB_PORT,  } from './constants';

import configStyles from '../../../styles/config-pages.module.scss';
import { UpdateArgs } from '../../../types/config-section';

export default function EditInstanceDetails() {
  const [formDataValues, setFormDataValues] = useState(null);
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { streamKey, ffmpegPath, rtmpServerPort, webServerPort } = serverConfig;

  const [copyIsVisible, setCopyVisible]  = useState(false);

  const COPY_TOOLTIP_TIMEOUT = 3000;

  useEffect(() => {
    setFormDataValues({
      streamKey, ffmpegPath, rtmpServerPort, webServerPort
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
  }

  function generateStreamKey () {
    let key = '';
    for (let i = 0; i < 3; i+=1) {
      key += Math.random().toString(36).substring(2);
    }

    handleFieldChange({ fieldName: 'streamKey', value: key });
  }

  function copyStreamKey () {
    navigator.clipboard.writeText(formDataValues.streamKey)
      .then(() => {
        setCopyVisible(true);
        setTimeout(() => setCopyVisible(false), COPY_TOOLTIP_TIMEOUT);
      });
  }

  return (  
    <div className={configStyles.publicDetailsContainer}>
      <div className={configStyles.textFieldsSection}>
        <TextFieldWithSubmit
          fieldName="streamKey"
          {...TEXTFIELD_PROPS_STREAM_KEY}
          value={formDataValues.streamKey}
          initialValue={streamKey}
          type={TEXTFIELD_TYPE_PASSWORD}
          onChange={handleFieldChange}
        />
        <div>
          <span style={{fontSize: '0.75em', color: '#ff7777', marginRight: '0.5em'}}>
            Save this key somewhere safe,
            you will need it to stream or login to the admin dashboard!
          </span>
          <Tooltip className="copy-tooltip"
                   title="Copied!"
                   trigger=""
                   visible={copyIsVisible}>
            <Button type="primary"
                    icon={<CopyOutlined />}
                    size="small"
                    onClick={copyStreamKey}
            />
          </Tooltip>
          <Button type="primary"
                  icon={<RedoOutlined />}
                  size="small"
                  onClick={generateStreamKey}
          />
        </div>
        <TextFieldWithSubmit
          fieldName="ffmpegPath"
          {...TEXTFIELD_PROPS_FFMPEG}
          value={formDataValues.ffmpegPath}
          initialValue={ffmpegPath}
          onChange={handleFieldChange}
        />
        <TextFieldWithSubmit
          fieldName="webServerPort"
          {...TEXTFIELD_PROPS_WEB_PORT}
          value={formDataValues.webServerPort}
          initialValue={webServerPort}
          type={TEXTFIELD_TYPE_NUMBER}
          onChange={handleFieldChange}
        />
        <TextFieldWithSubmit
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


