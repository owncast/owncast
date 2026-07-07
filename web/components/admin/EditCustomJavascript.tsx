import React, { useState, useEffect, useContext, FC } from 'react';
import { Typography, Button } from 'antd';
import CodeMirror, { EditorView } from '@uiw/react-codemirror';
import { bbedit } from '@uiw/codemirror-theme-bbedit';
import { javascript } from '@codemirror/lang-javascript';

import { ServerStatusContext } from '../../utils/server-status-context';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  API_CUSTOM_JAVASCRIPT,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import { FormStatusIndicator } from './FormStatusIndicator';

const { Title } = Typography;

// eslint-disable-next-line import/prefer-default-export
export const EditCustomJavascript: FC = () => {
  const [content, setContent] = useState('/* Enter custom Javascript here */');
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const [hasChanged, setHasChanged] = useState(false);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { customJavascript: initialContent } = instanceDetails;

  let resetTimer = null;

  // Clear out any validation states and messaging
  const resetStates = () => {
    setSubmitStatus(null);
    setHasChanged(false);
    clearTimeout(resetTimer);
    resetTimer = null;
  };

  // posts all the tags at once as an array obj
  async function handleSave() {
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));
    await postConfigUpdateToAPI({
      apiPath: API_CUSTOM_JAVASCRIPT,
      data: { value: content },
      onSuccess: (message: string) => {
        setFieldInConfigState({
          fieldName: 'customJavascript',
          value: content,
          path: 'instanceDetails',
        });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, message));
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
      },
    });
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  }

  useEffect(() => {
    setContent(initialContent);
  }, [instanceDetails]);

  const onValueChange = React.useCallback(value => {
    setContent(value);
    if (value !== initialContent && !hasChanged) {
      setHasChanged(true);
    } else if (value === initialContent && hasChanged) {
      setHasChanged(false);
    }
  }, []);

  return (
    <div className="edit-custom-css">
      <Title level={3} className="section-title">
        Customize your page with Javascript
      </Title>

      <p className="description">
        Insert custom Javascript into your Owncast page to add your own functionality or to add 3rd
        party scripts. Read more about how to use this feature in the{' '}
        <a href="https://owncast.online/docs/website/" rel="noopener noreferrer" target="_blank">
          Web page documentation.
        </a>
        .
      </p>
      <p className="description">Please use raw Javascript, no HTML or any script tags.</p>

      <CodeMirror
        value={content}
        placeholder="/* Enter custom Javascript here */"
        theme={bbedit}
        height="200px"
        extensions={[javascript(), EditorView.lineWrapping]}
        onChange={onValueChange}
      />

      <br />
      <div className="page-content-actions">
        {hasChanged && (
          <Button type="primary" onClick={handleSave}>
            Save
          </Button>
        )}
        <FormStatusIndicator status={submitStatus} />
      </div>
    </div>
  );
};
