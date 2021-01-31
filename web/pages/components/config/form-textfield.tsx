import React, { useEffect, useState, useContext } from 'react';
import { Button, Input, InputNumber } from 'antd';
import { FormItemProps } from 'antd/es/form';

import { RESET_TIMEOUT, postConfigUpdateToAPI } from './constants';

import { FieldUpdaterFunc } from '../../../types/config-section';
import { ServerStatusContext } from '../../../utils/server-status-context';
import InfoTip from '../info-tip';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea';
export const TEXTFIELD_TYPE_URL = 'url';

interface TextFieldProps {
  apiPath: string;
  fieldName: string;

  configPath?: string;
  disabled?: boolean;
  initialValue?: string;
  label?: string;
  maxLength?: number;
  placeholder?: string;
  required?: boolean;
  tip?: string;
  type?: string;
  value?: string | number;
  onSubmit?: () => void;
  onBlur?: () => void;
  onChange?: FieldUpdaterFunc;
}


export default function TextField(props: TextFieldProps) {
  const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const [hasChanged, setHasChanged] = useState(false);
  const [fieldValueForSubmit, setFieldValueForSubmit] = useState<string | number>('');

  const serverStatusData = useContext(ServerStatusContext);
  const { setFieldInConfigState } = serverStatusData || {};

  let resetTimer = null;

  const {
    apiPath,
    configPath = '',
    disabled = false,
    fieldName,
    initialValue,
    label,
    maxLength,
    onBlur,
    onChange,
    onSubmit,
    placeholder,
    required,
    tip,
    type,
    value,
  } = props;

  // Clear out any validation states and messaging
  const resetStates = () => {
    setSubmitStatus('');
    setHasChanged(false);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  useEffect(() => {
    // TODO: Add native validity checks here, somehow
    // https://developer.mozilla.org/en-US/docs/Web/API/ValidityState
    // const hasValidity = (type !== TEXTFIELD_TYPE_NUMBER && e.target.validity.valid) || type === TEXTFIELD_TYPE_NUMBER ;
    if ((required && (value === '' || value === null)) || value === initialValue) {
      setHasChanged(false);
    } else {
      // show submit button
      resetStates();
      setHasChanged(true);
      setFieldValueForSubmit(value);
    }
  }, [value]);

  // if field is required but value is empty, or equals initial value, then don't show submit/update button. otherwise clear out any result messaging and display button.
  const handleChange = (e: any) => {
    const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;

    // if an extra onChange handler was sent in as a prop, let's run that too.
    if (onChange) {
      onChange({ fieldName, value: val });
    }
  };

  // if you blur a required field with an empty value, restore its original value in state (parent's state), if an onChange from parent is available.
  const handleBlur = e => {
    if (!onChange) {
      return;
    }
    const val = e.target.value;
    if (required && val === '') {
      onChange({ fieldName, value: initialValue });
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

  const fieldId = `field-${fieldName}`;

   return (
    <div className={`textfield-container type-${type}`}>
      { required ? <span className="required-label">*</span> : null }
      <label htmlFor={fieldId} className="textfield-label">{label}</label>
      <div className="textfield">
        <Field
          id={fieldId}
          className={`field ${fieldId}`}
          {...fieldProps}
          allowClear
          placeholder={placeholder}
          maxLength={maxLength}
          onChange={handleChange}
          onBlur={handleBlur}
          disabled={disabled}
          value={value}
        />
      </div>
      <InfoTip tip={tip} />
      {submitStatus}
      {submitStatusMessage}

      { hasChanged ? <Button type="primary" size="small" className="submit-button" onClick={handleSubmit}>Update</Button> : null }

    </div>
  ); 
}

TextField.defaultProps = {
  configPath: '',
  disabled: false,
  initialValue: '',
  label: '',
  maxLength: null,
  placeholder: '',
  required: false,
  tip: '',
  type: TEXTFIELD_TYPE_TEXT,
  value: '',
  onSubmit: () => {},
  onBlur: () => {},
  onChange: () => {},
};
