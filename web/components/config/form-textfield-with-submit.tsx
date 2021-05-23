import React, { useEffect, useState, useContext } from 'react';
import { Button } from 'antd';
import classNames from 'classnames';
import { RESET_TIMEOUT, postConfigUpdateToAPI } from '../../utils/config-constants';

import { ServerStatusContext } from '../../utils/server-status-context';
import TextField, { TextFieldProps } from './form-textfield';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { UpdateArgs } from '../../types/config-section';
import FormStatusIndicator from './form-status-indicator';

export const TEXTFIELD_TYPE_TEXT = 'default';
export const TEXTFIELD_TYPE_PASSWORD = 'password'; // Input.Password
export const TEXTFIELD_TYPE_NUMBER = 'numeric';
export const TEXTFIELD_TYPE_TEXTAREA = 'textarea';
export const TEXTFIELD_TYPE_URL = 'url';

interface TextFieldWithSubmitProps extends TextFieldProps {
  apiPath: string;
  configPath?: string;
  initialValue?: string;
}

export default function TextFieldWithSubmit(props: TextFieldWithSubmitProps) {
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  const [hasChanged, setHasChanged] = useState(false);

  const serverStatusData = useContext(ServerStatusContext);
  const { setFieldInConfigState } = serverStatusData || {};

  let resetTimer = null;

  const {
    apiPath,
    configPath = '',
    initialValue,
    useTrim,
    ...textFieldProps // rest of props
  } = props;

  const { fieldName, required, tip, status, value, onChange, onSubmit } = textFieldProps;

  // Clear out any validation states and messaging
  const resetStates = () => {
    setSubmitStatus(null);
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
    }
  }, [value]);

  // if field is required but value is empty, or equals initial value, then don't show submit/update button. otherwise clear out any result messaging and display button.
  const handleChange = ({ fieldName: changedFieldName, value: changedValue }: UpdateArgs) => {
    if (onChange) {
      onChange({ fieldName: changedFieldName, value: useTrim ? changedValue.trim() : changedValue });
    }
  };

  // if you blur a required field with an empty value, restore its original value in state (parent's state), if an onChange from parent is available.
  const handleBlur = ({ value: changedValue }: UpdateArgs) => {
    if (onChange && required && changedValue === '') {
      onChange({ fieldName, value: initialValue });
    }
  };

  // how to get current value of input
  const handleSubmit = async () => {
    if ((required && value !== '') || value !== initialValue) {
      setSubmitStatus(createInputStatus(STATUS_PROCESSING));

      await postConfigUpdateToAPI({
        apiPath,
        data: { value },
        onSuccess: () => {
          setFieldInConfigState({ fieldName, value, path: configPath });
          setSubmitStatus(createInputStatus(STATUS_SUCCESS));
        },
        onError: (message: string) => {
          setSubmitStatus(createInputStatus(STATUS_ERROR, `There was an error: ${message}`));
        },
      });
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);

      // if an extra onSubmit handler was sent in as a prop, let's run that too.
      if (onSubmit) {
        onSubmit();
      }
    }
  };

  const textfieldContainerClass = classNames({
    'textfield-with-submit-container': true,
    submittable: hasChanged,
  });
  return (
    <div className={textfieldContainerClass}>
      <div className="textfield-component">
        <TextField
          {...textFieldProps}
          onSubmit={null}
          onBlur={handleBlur}
          onChange={handleChange}
        />
      </div>
      <div className="formfield-container lower-container">
        <p className="label-spacer" />
        <div className="lower-content">
          <div className="field-tip">{tip}</div>
          <FormStatusIndicator status={status || submitStatus} />
          <div className="update-button-container">
            <Button
              type="primary"
              size="small"
              className="submit-button"
              onClick={handleSubmit}
              disabled={!hasChanged}
            >
              Update
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

TextFieldWithSubmit.defaultProps = {
  configPath: '',
  initialValue: '',
};
