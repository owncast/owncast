import React, { useState, useContext } from 'react';
import { Switch } from 'antd';
import { FormItemProps } from 'antd/es/form';

import { RESET_TIMEOUT, SUCCESS_STATES, postConfigUpdateToAPI } from './constants';

import { ToggleSwitchProps } from '../../../types/config-section';
import { ServerStatusContext } from '../../../utils/server-status-context';
import InfoTip from '../info-tip';


export default function ToggleSwitch(props: ToggleSwitchProps) {
  const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');

  let resetTimer = null;

  const serverStatusData = useContext(ServerStatusContext);
  const { setFieldInConfigState } = serverStatusData || {};
  
  const {
    apiPath,
    checked,
    configPath = '',
    disabled = false,
    fieldName,
    label,
    tip,
  } = props;

  const resetStates = () => {
    setSubmitStatus('');
    clearTimeout(resetTimer);
    resetTimer = null;
  }

  const handleChange = async isChecked => {
    setSubmitStatus('validating');
    await postConfigUpdateToAPI({
      apiPath,
      data: { value: isChecked },
      onSuccess: () => {
        setFieldInConfigState({ fieldName, value: isChecked, path: configPath });
        setSubmitStatus('success');
      },
      onError: (message: string) => {
        setSubmitStatus('error');
        setSubmitStatusMessage(`There was an error: ${message}`);
      },
    });
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  }

  const {
    icon: newStatusIcon = null,
    message: newStatusMessage = '',
  } = SUCCESS_STATES[submitStatus] || {};

  return (
    <div className="toggleswitch-container">
      <div className="toggleswitch">
        <Switch
          className={`switch field-${fieldName}`}
          loading={submitStatus === 'validating'}
          onChange={handleChange}
          defaultChecked={checked}
          checkedChildren="ON"  
          unCheckedChildren="OFF"
          disabled={disabled}
        />
        <span className="label">{label} <InfoTip tip={tip} /></span>
        {submitStatus}
      </div>
      <div className={`status-message ${submitStatus || ''}`}>
        {newStatusIcon} {newStatusMessage} {submitStatusMessage}
      </div>
    </div>
  ); 
}
