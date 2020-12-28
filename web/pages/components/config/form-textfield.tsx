/*
- auto saves ,ajax call (submit when blur or onEnter)
- set default text
- show error state/confirm states
- show info
- label
- min/max length

- populate with curren val (from local sstate)

load page, 
get all config vals, 
save to local state/context.
read vals from there.
update vals to state, andthru api.


*/
import React, { useState, useContext } from 'react';
import { Form, Input, Tooltip } from 'antd';
import { FormItemProps } from 'antd/es/form';

import { InfoCircleOutlined } from '@ant-design/icons';

import { TEXTFIELD_DEFAULTS, TEXT_MAXLENGTH } from './defaults';

import { TextFieldProps } from '../../../types/config-section';
import { fetchData, SERVER_CONFIG_UPDATE_URL } from '../../../utils/apis';
import { ServerStatusContext } from '../../../utils/server-status-context';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';

export default function TextField(props: TextFieldProps) {
  const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const serverStatusData = useContext(ServerStatusContext);
  const { setConfigField } = serverStatusData || {};


  const {
    fieldName,
  } = props;

  const {
    apiPath = '',
    defaultValue = '', // if empty
    configPath = '',
    maxLength = TEXT_MAXLENGTH,
    placeholder = '',
    label = '',
    tip = '',
  } = TEXTFIELD_DEFAULTS[fieldName] || {};

  const postUpdateToAPI = async (postValue: any) => {
    setSubmitStatus('validating');
    const result = await fetchData(`${SERVER_CONFIG_UPDATE_URL}${apiPath}`, {
      data: { value: postValue },
      method: 'POST',
      auth: true,
    });
    if (result.success) {
      setConfigField({ fieldName, value: postValue, path: configPath });
      setSubmitStatus('success');
    } else {
      setSubmitStatus('warning');
      setSubmitStatusMessage(`There was an error: ${result.message}`);
    }
  };

  const handleEnter = (event: React.KeyboardEvent<HTMLInputElement>) => {
    const newValue = event.target.value;
    if (newValue !== '') {
      postUpdateToAPI(newValue);
    }
  }
  const handleBlur = (event: React.FocusEvent<HTMLInputElement>) => {
    const newValue = event.target.value;
    if (newValue !== '') {
      console.log("blur post..", newValue)
      postUpdateToAPI(newValue);
    } else {
      // event.target.value = value;
    }
  }

   return (
    <div className="textfield">
      <Form.Item
        label={label}
        name={fieldName}
        hasFeedback
        validateStatus={submitStatus}
        help={submitStatusMessage}
      >
       <Input
          className="field"
          allowClear
          placeholder={placeholder}
          maxLength={maxLength}
          onPressEnter={handleEnter}
          // onBlur={handleBlur}
        />
      </Form.Item>
      <div className="info">
        <Tooltip title={tip}>
          <InfoCircleOutlined />
        </Tooltip>
      </div>
    </div>
  ); 
}
