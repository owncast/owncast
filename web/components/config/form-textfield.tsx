import React from 'react';
import classNames from 'classnames';
import { Input, InputNumber } from 'antd';
import { FieldUpdaterFunc } from '../../types/config-section';
// import InfoTip from '../info-tip';
import { StatusState } from '../../utils/input-statuses';
import FormStatusIndicator from './form-status-indicator';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric'; // InputNumber
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea'; // Input.TextArea
export const TEXTFIELD_TYPE_URL = 'url';

export interface TextFieldProps {
  fieldName: string;

  onSubmit?: () => void;
  onPressEnter?: () => void;

  className?: string;
  disabled?: boolean;
  label?: string;
  maxLength?: number;
  pattern?: string;
  placeholder?: string;
  required?: boolean;
  status?: StatusState;
  tip?: string;
  type?: string;
  useTrim?: boolean;
  value?: string | number;
  onBlur?: FieldUpdaterFunc;
  onChange?: FieldUpdaterFunc;
}

export default function TextField(props: TextFieldProps) {
  const {
    className,
    disabled,
    fieldName,
    label,
    maxLength,
    onBlur,
    onChange,
    onPressEnter,
    pattern,
    placeholder,
    required,
    status,
    tip,
    type,
    useTrim,
    value,
  } = props;

  const handleChange = (e: any) => {
    // if an extra onChange handler was sent in as a prop, let's run that too.
    if (onChange) {
      const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;
      onChange({ fieldName, value: useTrim ? val.trim() : val });
    }
  };

  // if you blur a required field with an empty value, restore its original value in state (parent's state), if an onChange from parent is available.
  const handleBlur = (e: any) => {
    const val = e.target.value;
    if (onBlur) {
      onBlur({ value: val });
    }
  };

  const handlePressEnter = () => {
    if (onPressEnter) {
      onPressEnter();
    }
  };

  // display the appropriate Ant text field
  let Field = Input as
    | typeof Input
    | typeof InputNumber
    | typeof Input.TextArea
    | typeof Input.Password;
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
      max: 10 ** maxLength - 1,
    };
  } else if (type === TEXTFIELD_TYPE_URL) {
    fieldProps = {
      type: 'url',
      pattern,
    };
  }

  const fieldId = `field-${fieldName}`;

  const { type: statusType } = status || {};

  const containerClass = classNames({
    'formfield-container': true,
    'textfield-container': true,
    [`type-${type}`]: true,
    required,
    [`status-${statusType}`]: status,
  });

  return (
    <div className={containerClass}>
      {label ? (
        <div className="label-side">
          <label htmlFor={fieldId} className="formfield-label">
            {label}
          </label>
        </div>
      ) : null}

      <div className="input-side">
        <div className="input-group">
          <Field
            id={fieldId}
            className={`field ${className} ${fieldId}`}
            {...fieldProps}
            {...(type !== TEXTFIELD_TYPE_NUMBER && { allowClear: true })}
            placeholder={placeholder}
            maxLength={maxLength}
            onChange={handleChange}
            onBlur={handleBlur}
            onPressEnter={handlePressEnter}
            disabled={disabled}
            value={value as number | (readonly string[] & number)}
          />
        </div>
        <FormStatusIndicator status={status} />
        <p className="field-tip">{tip}</p>
      </div>
    </div>
  );
}

TextField.defaultProps = {
  className: '',
  disabled: false,
  label: '',
  maxLength: 255,

  placeholder: '',
  required: false,
  status: null,
  tip: '',
  type: TEXTFIELD_TYPE_TEXT,
  value: '',
  onSubmit: () => {},
  onBlur: () => {},
  onChange: () => {},
  onPressEnter: () => {},
};
