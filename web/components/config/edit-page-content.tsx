// EDIT CUSTOM DETAILS ON YOUR PAGE
import React, { useState, useEffect, useContext } from 'react';
import { Typography, Button } from 'antd';
import dynamic from 'next/dynamic';
import MarkdownIt from 'markdown-it';

import { ServerStatusContext } from '../../utils/server-status-context';
import {
  postConfigUpdateToAPI,
  RESET_TIMEOUT,
  API_CUSTOM_CONTENT,
} from '../../utils/config-constants';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
} from '../../utils/input-statuses';
import FormStatusIndicator from './form-status-indicator';

import 'react-markdown-editor-lite/lib/index.css';

const mdParser = new MarkdownIt(/* Markdown-it options */);
const MdEditor = dynamic(() => import('react-markdown-editor-lite'), {
  ssr: false,
});

const { Title } = Typography;

export default function EditPageContent() {
  const [content, setContent] = useState('');
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);
  const [hasChanged, setHasChanged] = useState(false);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { extraPageContent: initialContent } = instanceDetails;

  let resetTimer = null;

  function handleEditorChange({ text }) {
    setContent(text);
    if (text !== initialContent && !hasChanged) {
      setHasChanged(true);
    } else if (text === initialContent && hasChanged) {
      setHasChanged(false);
    }
  }

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
      apiPath: API_CUSTOM_CONTENT,
      data: { value: content },
      onSuccess: (message: string) => {
        setFieldInConfigState({
          fieldName: 'extraPageContent',
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

  return (
    <div className="edit-page-content">
      <Title level={3} className="section-title">
        Custom Page Content
      </Title>

      <p className="description">
        Edit the content of your page by using simple{' '}
        <a href="https://www.markdownguide.org/basic-syntax/">Markdown syntax</a>.
      </p>

      <MdEditor
        style={{ height: '30em' }}
        value={content}
        renderHTML={(c: string) => mdParser.render(c)}
        onChange={handleEditorChange}
        config={{
          htmlClass: 'markdown-editor-preview-pane',
          markdownClass: 'markdown-editor-pane',
        }}
      />

      <div className="page-content-actions">
        {hasChanged ? (
          <Button type="primary" onClick={handleSave}>
            Save
          </Button>
        ) : null}
        <FormStatusIndicator status={submitStatus} />
      </div>
    </div>
  );
}
