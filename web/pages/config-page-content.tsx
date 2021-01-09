import React, { useState, useEffect, useContext } from 'react';
import { Typography, Button } from "antd";
import { FormItemProps } from 'antd/lib/form';
import dynamic from 'next/dynamic';
import MarkdownIt from 'markdown-it';

import { ServerStatusContext } from '../utils/server-status-context';
import { TEXTFIELD_DEFAULTS, postConfigUpdateToAPI, RESET_TIMEOUT, SUCCESS_STATES} from './components/config/constants';

import 'react-markdown-editor-lite/lib/index.css';

const { Title } = Typography;

const mdParser = new MarkdownIt(/* Markdown-it options */);

const MdEditor = dynamic(() => import('react-markdown-editor-lite'), {
    ssr: false,
});

export default function PageContentEditor() {
    const [content, setContent] = useState('');
    const [submitStatus, setSubmitStatus] = useState<FormItemProps['validateStatus']>('');
    const [submitStatusMessage, setSubmitStatusMessage] = useState('');
    const [hasChanged, setHasChanged] = useState(false);
  
    const serverStatusData = useContext(ServerStatusContext);
    const { serverConfig, setFieldInConfigState } = serverStatusData || {};
  
    const { instanceDetails } = serverConfig;
    const { extraPageContent: initialContent } = instanceDetails;
  
    const { apiPath } = TEXTFIELD_DEFAULTS.instanceDetails.extraPageContent;

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
      setSubmitStatus('');
      setHasChanged(false);
      clearTimeout(resetTimer);
      resetTimer = null;
    };

    // posts all the tags at once as an array obj
    async function handleSave() {
      setSubmitStatus('validating');
      await postConfigUpdateToAPI({
        apiPath,
        data: { value: content },
        onSuccess: () => {
          setFieldInConfigState({ fieldName: 'extraPageContent', value: content, path: apiPath });
          setSubmitStatus('success');
        },
        onError: (message: string) => {
          setSubmitStatus('error');
          setSubmitStatusMessage(`There was an error: ${message}`);
        },
      });
      resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
    }
  
    useEffect(() => {
      setContent(initialContent);
    }, [instanceDetails]);

    const {
      icon: newStatusIcon = null,
      message: newStatusMessage = '',
    } = SUCCESS_STATES[submitStatus] || {};


    return (
      <div className="config-page-content-form">
        <Title level={2}>Edit custom content</Title>

        <p>Add some content about your site with the Markdown editor below. This content shows up at the bottom half of your Owncast page.</p>

        <MdEditor
          style={{ height: "30em" }}
          value={content}
          renderHTML={(c: string) => mdParser.render(c)}
          onChange={handleEditorChange}
          config={{
            htmlClass: 'markdown-editor-preview-pane',
            markdownClass: 'markdown-editor-pane',
          }}
        />
        <div className="page-content-actions">
          { hasChanged ? <Button type="primary" onClick={handleSave}>Save</Button> : null }
          <div className={`status-message ${submitStatus || ''}`}>
            {newStatusIcon} {newStatusMessage} {submitStatusMessage}
          </div>
        </div>
      </div>
    );
}
