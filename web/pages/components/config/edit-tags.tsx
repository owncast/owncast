/* eslint-disable react/no-array-index-key */
import React, { useContext, useState, useEffect } from 'react';
import { Typography, Tag, Input } from 'antd';

import { ServerStatusContext } from '../../../utils/server-status-context';
import { FIELD_PROPS_TAGS, RESET_TIMEOUT, SUCCESS_STATES, postConfigUpdateToAPI } from './constants';

const { Title } = Typography;

export default function EditInstanceTags() {
  const [newTagInput, setNewTagInput] = useState('');
  const [submitStatus, setSubmitStatus] = useState(null);
  const [submitStatusMessage, setSubmitStatusMessage] = useState('');
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { tags = [] } = instanceDetails;

  const {
    apiPath,
    maxLength,
    placeholder,
    configPath,
  } = FIELD_PROPS_TAGS;

  let resetTimer = null;

  useEffect(() => {
    return () => {
      clearTimeout(resetTimer);
    }  
  }, []);

  const resetStates = () => {
    setSubmitStatus(null);
    setSubmitStatusMessage('');
    resetTimer = null;
    clearTimeout(resetTimer);
  }

  // posts all the tags at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    await postConfigUpdateToAPI({
      apiPath,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'tags', value: postValue, path: configPath });
        setSubmitStatus('success');
        setSubmitStatusMessage('Tags updated.');
        setNewTagInput('');
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus('error');
        setSubmitStatusMessage(message);
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const handleInputChange = e => {
    if (submitStatusMessage !== '') {
      setSubmitStatusMessage('');
    }
    setNewTagInput(e.target.value);
  };

  // send to api and do stuff
  const handleSubmitNewTag = () => {
    resetStates();
    const newTag = newTagInput.trim();
    if (newTag === '') {
      setSubmitStatusMessage('Please enter a tag');
      return;
    } 
    if (tags.some(tag => tag.toLowerCase() === newTag.toLowerCase())) {
      setSubmitStatusMessage('This tag is already used!');
      return;
    }

    const updatedTags = [...tags, newTag];
    postUpdateToAPI(updatedTags);
  };

  const handleDeleteTag = index => {
    resetStates();
    const updatedTags = [...tags];
    updatedTags.splice(index, 1);
    postUpdateToAPI(updatedTags);
  }

  const {
    icon: newStatusIcon = null,
    message: newStatusMessage = '',
  } = SUCCESS_STATES[submitStatus] || {};

  return (
    <div className="tag-editor-container">

      <Title level={3}>Add Tags</Title>
      <p>This is a great way to categorize your Owncast server on the Directory!</p>

      <div className="tag-current-tags">
        {tags.map((tag, index) => {
          const handleClose = () => {
            handleDeleteTag(index);
          };
          return (
            <Tag closable onClose={handleClose} key={`tag-${tag}-${index}`}>{tag}</Tag>
          );
        })}
      </div>
      <div className={`status-message ${submitStatus || ''}`}>
        {newStatusIcon} {newStatusMessage} {submitStatusMessage}
      </div>
      <div className="add-new-tag-section">
        <Input
          type="text"
          className="new-tag-input"
          value={newTagInput}
          onChange={handleInputChange}
          onPressEnter={handleSubmitNewTag}
          maxLength={maxLength}
          placeholder={placeholder}
          allowClear
        />
      </div>
    </div>
  );
}
