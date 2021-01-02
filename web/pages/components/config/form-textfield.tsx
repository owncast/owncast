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
import { Button, Form, Input, InputNumber, Tooltip } from 'antd';
import { FormItemProps } from 'antd/es/form';

import { InfoCircleOutlined } from '@ant-design/icons';

import { TEXTFIELD_DEFAULTS, TEXT_MAXLENGTH, RESET_TIMEOUT } from './constants';

import { TextFieldProps } from '../../../types/config-section';
import { fetchData, SERVER_CONFIG_UPDATE_URL } from '../../../utils/apis';
import { ServerStatusContext } from '../../../utils/server-status-context';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea';


export default function TextField(props: TextFieldProps) {
  const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const [hasChanged, setHasChanged] = useState(false);
  const [fieldValueForSubmit, setFieldValueForSubmit] = useState('');

  let resetTimer = null;

  const serverStatusData = useContext(ServerStatusContext);
  const { setConfigField } = serverStatusData || {};
  
  const {
    fieldName,
    type,
    initialValues = {},
    handleResetValue,
  } = props;

  const initialValue = initialValues[fieldName] || '';
  
  const {
    apiPath = '',
    configPath = '',
    maxLength = TEXT_MAXLENGTH,
    // placeholder = '',
    label = '',
    tip = '',
  } = TEXTFIELD_DEFAULTS[fieldName] || {};

  const resetStates = () => {
    setSubmitStatus('');
    setHasChanged(false);
    clearTimeout(resetTimer);
    resetTimer = null;
  }

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
      setSubmitStatus('error');
      setSubmitStatusMessage(`There was an error: ${result.message}`);
    }
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  };

  const handleChange = e => {
    const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;
    if (val === '' || val === initialValue) {
      setHasChanged(false);
    } else {
      resetStates();
      setHasChanged(true);
      setFieldValueForSubmit(val);
    }
  };

  const handleBlur = e => {
    const val = e.target.value;
    if (val === '') {
      handleResetValue(fieldName);
    }
  };

  // how to get current value of input
  const handleSubmit = () => {
    if (fieldValueForSubmit !== '' && fieldValueForSubmit !== initialValue) {
      postUpdateToAPI(fieldValueForSubmit);
    }
  }

  let Field = Input;
  let fieldProps = {};
  if (type === TEXTFIELD_TYPE_TEXTAREA) {
    Field = Input.TextArea;
    fieldProps = {
      autoSize: true,
    };
  } else if (type === TEXTFIELD_TYPE_PASSWORD) {
    Field = Input.Password;
    fieldProps = {
      visibilityToggle: true,
    };
  } else if (type === TEXTFIELD_TYPE_NUMBER) {
    Field = InputNumber;
  }

   return (
    <div className="textfield-container">
      <div className="textfield">
       <span className="info">
          <Tooltip title={tip}>
            <InfoCircleOutlined />
          </Tooltip>
        </span>
        <Form.Item
          label={label}
          name={fieldName}
          hasFeedback
          validateStatus={submitStatus}
          help={submitStatusMessage}
        >
        <Field
            className="field"
            allowClear
            placeholder={initialValue}
            maxLength={maxLength}
            onChange={handleChange}
            onBlur={handleBlur}
            {...fieldProps}
          />
        </Form.Item>
   
      </div>

      { hasChanged ? <Button type="primary" size="small" className="submit-button" onClick={handleSubmit}>Update</Button> : null }

    </div>
  ); 
}
