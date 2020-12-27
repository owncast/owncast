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
import React from 'react';
import { Form, Input } from 'antd';


interface TextFieldProps {
  onSubmit: (value: string) => void;
  label: string;
  defaultValue: string;
  value: string;
  helpInfo: string;
  maxLength: number;
  type: string;
}

// // do i need this?
// export const initialProps: TextFieldProps = {
// }

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; //Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';

export default function TextField(props: TextFieldProps) {
  const {
    label,
    defaultValue,
    value,
    onSubmit,
    helpInfo,
    maxLength,
    type,
  } = props;
  
  return (
    <div className="textfield">
      <Form.Item
        label={label}
        hasFeedback
        validateStatus="error"
        help="Should be combination of numbers & alphabets"
      >
       <Input placeholder="Owncast" value={value} />
      </Form.Item>
    </div>
  ); 
}
