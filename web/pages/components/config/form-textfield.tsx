
import React, { useState, useContext } from 'react';
import { Button, Form, Input, InputNumber } from 'antd';
import { FormItemProps } from 'antd/es/form';

import { TEXTFIELD_DEFAULTS, TEXT_MAXLENGTH, RESET_TIMEOUT, postConfigUpdateToAPI } from './constants';

import { TextFieldProps } from '../../../types/config-section';
import { ServerStatusContext } from '../../../utils/server-status-context';
import InfoTip from '../info-tip';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea';
export const TEXTFIELD_TYPE_URL = 'url';


export default function TextField(props: TextFieldProps) {
  const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const [hasChanged, setHasChanged] = useState(false);
  const [fieldValueForSubmit, setFieldValueForSubmit] = useState('');

  let resetTimer = null;

  const serverStatusData = useContext(ServerStatusContext);
  const { setFieldInConfigState } = serverStatusData || {};
  
  const {
    configPath = '',
    disabled = false,
    fieldName,
    handleResetValue = () => {},
    initialValues = {},
    placeholder,
    onSubmit,
    onBlur,
    onChange,
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
  };

  // if field is required but value is empty, or equals initial value, then don't show submit/update button. otherwise clear out any result messaging and display button.
  const handleChange = (e: any) => {
    const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;

    // https://developer.mozilla.org/en-US/docs/Web/API/ValidityState
    const hasValidity = (type !== TEXTFIELD_TYPE_NUMBER && e.target.validity.valid) || type === TEXTFIELD_TYPE_NUMBER ;

    if ((required && (val === '' || val === null)) || val === initialValue || !hasValidity) {
      setHasChanged(false);
    } else {
      // show submit button
      resetStates();
      setHasChanged(true);
      setFieldValueForSubmit(val);
    }
    // if an extra onChange handler was sent in as a prop, let's run that too.
    if (onChange) {
      onChange();
    }
  };

  // if you blur a required field with an empty value, restore its original value
  const handleBlur = e => {
    const val = e.target.value;
    if (required && val === '') {
      handleResetValue(fieldName);
    }

    // if an extra onBlur handler was sent in as a prop, let's run that too.
    if (onBlur) {
      onBlur();
    }
  };

  // how to get current value of input
  const handleSubmit = async () => {
    if ((required && fieldValueForSubmit !== '') || fieldValueForSubmit !== initialValue) {
      setSubmitStatus('validating');

      await postConfigUpdateToAPI({
        apiPath,
        data: { value: fieldValueForSubmit },
        onSuccess: () => {
          setFieldInConfigState({ fieldName, value: fieldValueForSubmit, path: configPath });
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
    fieldProps = {
      type: 'number',
      min: 1,
      max: (10**maxLength) - 1,
      onKeyDown: (e: React.KeyboardEvent) => {
        if (e.target.value.length > maxLength - 1 )
        e.preventDefault()
        return false;
      }
    };
  } else if (type === TEXTFIELD_TYPE_URL) {
    fieldProps = {
      type: 'url',
    };
  }

   return (
    <div className={`textfield-container type-${type}`}>
      <div className="textfield">
       <InfoTip tip={tip} />
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
