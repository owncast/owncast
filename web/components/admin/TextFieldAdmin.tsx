import React, { FC, useEffect } from 'react';
import classNames from 'classnames';
import { Input, Form } from 'antd';
import { FieldUpdaterFunc } from '../../types/config-section';
// import InfoTip from '../info-tip';
import { StatusState } from '../../utils/input-statuses';
// import { FormStatusIndicator } from './FormStatusIndicator';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric'; // InputNumber
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea'; // Input.TextArea
export const TEXTFIELD_TYPE_URL = 'url';

export type TextFieldAdminProps = {
  fieldName: string;

  onSubmit?: () => void;
  onPressEnter?: () => void;

  className?: string;
  disabled?: boolean;
  label?: string;
  placeholder?: string;
  required?: boolean;
  status?: StatusState;
  tip?: string;
  type?: string;
  useTrim?: boolean;
  useTrimLead?: boolean;
  value?: string | number;
  onBlur?: FieldUpdaterFunc;
  onChange?: FieldUpdaterFunc;
};

export const TextFieldAdmin: FC<TextFieldAdminProps> = ({
  className,
  disabled,
  fieldName,
  label,
  onBlur,
  onChange,
  onPressEnter,
  placeholder,
  required,
  status,
  type,
  useTrim,
  value,
}) => {
  const handleChange = (e: any) => {
    // if an extra onChange handler was sent in as a prop, let's run that too.
    if (onChange) {
      const val = type === TEXTFIELD_TYPE_NUMBER ? e : e.target.value;
      onChange({ fieldName, value: useTrim ? val.trim() : val });
    }
  };

  const [form] = Form.useForm();
  // const formRef = useRef(null);

  useEffect(() => {
    // console.log('value: ', value);
    form.setFieldsValue({ adminPassword: value });
  }, [value]);

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

  const onFinish = (values: Store) => {
    console.log('Received values of form: ', values);
  };

  // display the appropriate Ant text field
  // let Field = Input as typeof Input.Password;
  // let fieldProps = {};
  // Field = Input.Password;
  // fieldProps = {
  //   visibilityToggle: true,
  // };

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
      <Form name="basic" form={form} initialValues={{ adminPassword: value }} onFinish={onFinish}>
        <Form.Item
          label={label}
          name="adminPassword"
          rules={[
            { min: 8, message: '- minimum 8 characters' },
            { max: 192, message: '- maximum 192 characters' },
            {
              pattern: /^(?=.*[a-z])/,
              message: '- at least one lowercase letter',
            },
            {
              pattern: /^(?=.*[A-Z])/,
              message: '- at least one uppercase letter',
            },
            {
              pattern: /\d/,
              message: '- at least one digit',
            },
            {
              pattern: /^(?=.*?[#?!@$%^&*-])/,
              message: '- at least one special character: !@#$%^&*',
            },
          ]}
        >
          <Input.Password
            id={fieldId}
            className={`field ${className} ${fieldId}`}
            {...(type !== TEXTFIELD_TYPE_NUMBER && { allowClear: true })}
            onChange={handleChange}
            onBlur={handleBlur}
            placeholder={placeholder}
            onPressEnter={handlePressEnter}
            disabled={disabled}
            value={value as number | (readonly string[] & number)}
          />
        </Form.Item>
      </Form>
    </div>
  );
};
export default TextFieldAdmin;

TextFieldAdmin.defaultProps = {
  className: '',
  disabled: false,

  placeholder: '',
  required: false,
  status: null,
  tip: '',
  type: TEXTFIELD_TYPE_TEXT,
  value: '',

  useTrim: false,
  useTrimLead: false,

  onSubmit: () => {},
  onBlur: () => {},
  onChange: () => {},
  onPressEnter: () => {},
};
