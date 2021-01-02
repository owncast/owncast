/* eslint-disable react/no-array-index-key */
import React, { useContext, useState, useEffect } from 'react';
import { Typography, Tag, Input } from 'antd';

import { ServerStatusContext } from '../../../utils/server-status-context';
import { fetchData, SERVER_CONFIG_UPDATE_URL } from '../../../utils/apis';
import { TEXTFIELD_DEFAULTS, RESET_TIMEOUT, SUCCESS_STATES } from './constants';

const { Title } = Typography;

export default function EditInstanceTags() {
  const [newTagInput, setNewTagInput] = useState('');
  const [submitStatus, setSubmitStatus] = useState(null);
  const [submitDetails, setSubmitDetails] = useState('');
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setConfigField } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { tags = [] } = instanceDetails;

  const {
    apiPath,
    maxLength,
    placeholder,
    configPath,
  } = TEXTFIELD_DEFAULTS.tags || {};

  let resetTimer = null;

  useEffect(() => {
    return () => {
      clearTimeout(resetTimer);
    }  
  }, []);

  const resetStates = () => {
    setSubmitStatus(null);
    setSubmitDetails('');
    resetTimer = null;
    clearTimeout(resetTimer);
  }

  // posts all the tags at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    // const result = await fetchData(`${SERVER_CONFIG_UPDATE_URL}${apiPath}`, {
    //   data: { value: postValue },
    //   method: 'POST',
    //   auth: true,
    // });

    const result = {
      success: true,
      message: 'success yay'
    }
    if (result.success) {
      setConfigField({ fieldName: 'tags', value: postValue, path: configPath });
      setSubmitStatus('success');
      setSubmitDetails('Tags updated.');
      setNewTagInput('');
    } else {
      setSubmitStatus('error');
      setSubmitDetails(result.message);
    }
    resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
  };

  const handleInputChange = e => {
    if (submitDetails !== '') {
      setSubmitDetails('');
    }
    setNewTagInput(e.target.value);
  };

  // send to api and do stuff
  const handleSubmitNewTag = () => {
    resetStates();
    const newTag = newTagInput.trim();
    if (newTag === '') {
      setSubmitDetails('Please enter a tag');
      return;
    } 
    if (tags.some(tag => tag.toLowerCase() === newTag.toLowerCase())) {
      setSubmitDetails('This tag is already used!');
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
      <div className={`add-new-status ${submitStatus || ''}`}>
        {newStatusIcon} {newStatusMessage} {submitDetails}
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
