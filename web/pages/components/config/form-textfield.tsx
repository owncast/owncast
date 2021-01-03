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

import { TEXTFIELD_DEFAULTS, TEXT_MAXLENGTH, RESET_TIMEOUT, postConfigUpdateToAPI } from './constants';

import { TextFieldProps } from '../../../types/config-section';
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
    configPath = '',
    disabled = false,
    fieldName,
    handleResetValue,
    initialValues = {},
    onSubmit,
    type,
  } = props;

  // Keep track of what the initial value is
  // Note: we're not using `initialValue` as a prop, because we expect this component to be controlled by a parent Ant <Form> which is doing a form.setFieldsValue() upstream.
  const initialValue = initialValues[fieldName] || '';
  
  // Get other static info we know about this field.
  const defaultDetails =  TEXTFIELD_DEFAULTS[configPath] || TEXTFIELD_DEFAULTS;
  const {
    apiPath = '',
    maxLength = TEXT_MAXLENGTH,
    placeholder = '',
    label = '',
    tip = '',
    required = false,
  } = defaultDetails[fieldName] || {};

  // Clear out any validation states and messaging
  const resetStates = () => {
    setSubmitStatus('');
    setHasChanged(false);
    clearTimeout(resetTimer);
    resetTimer = null;
  }

  // if field is required but value is empty, or equals initial value, then don't show submit/update button. otherwise clear out any result messaging and display button.
  const handleChange = (e: any) => {
    const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;
    if ((required && (val === '' || val === null)) || val === initialValue) {
      setHasChanged(false);
    } else {
      resetStates();
      setHasChanged(true);
      setFieldValueForSubmit(val);
    }
  };

  // if you blur a required field with an empty value, restore its original value
  const handleBlur = e => {
    const val = e.target.value;
    if (required && val === '') {
      handleResetValue(fieldName);
    }
  };

  // how to get current value of input
  const handleSubmit = async () => {
    if ((required && fieldValueForSubmit !== '') || fieldValueForSubmit !== initialValue) {
      // postUpdateToAPI(fieldValueForSubmit);
      setSubmitStatus('validating');

      await postConfigUpdateToAPI({
        apiPath,
        data: { value: fieldValueForSubmit },
        onSuccess: () => {
          setConfigField({ fieldName, value: fieldValueForSubmit, path: configPath });
          setSubmitStatus('success');
        },
        onError: (message: string) => {
          setSubmitStatus('error');
          setSubmitStatusMessage(`There was an error: ${message}`);
        },
      });
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);

      // if an extra onSubmit handler was sent in as a prop, let's run that too.
      if (onSubmit) {
        onSubmit();
      }
    }
  }

  // display the appropriate Ant text field
  let Field = Input as typeof Input | typeof InputNumber | typeof Input.TextArea | typeof Input.Password;
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
          required={required}
        >
        <Field
            className={`field field-${fieldName}`}
            allowClear
            placeholder={placeholder}
            maxLength={maxLength}
            onChange={handleChange}
            onBlur={handleBlur}
            disabled={disabled}
            {...fieldProps}
          />
        </Form.Item>
   
      </div>

      { hasChanged ? <Button type="primary" size="small" className="submit-button" onClick={handleSubmit}>Update</Button> : null }

    </div>
  ); 
}
