/* eslint-disable react/no-array-index-key */
import React, { useContext, useState, useEffect } from 'react';
import { Typography, Tag } from 'antd';

import { ServerStatusContext } from '../../utils/server-status-context';
import {
  FIELD_PROPS_TAGS,
  RESET_TIMEOUT,
  postConfigUpdateToAPI,
} from '../../utils/config-constants';
import TextField from './form-textfield';
import { UpdateArgs } from '../../types/config-section';
import {
  createInputStatus,
  StatusState,
  STATUS_ERROR,
  STATUS_PROCESSING,
  STATUS_SUCCESS,
  STATUS_WARNING,
} from '../../utils/input-statuses';

const { Title } = Typography;

export default function EditInstanceTags() {
  const [newTagInput, setNewTagInput] = useState<string>('');
  const [submitStatus, setSubmitStatus] = useState<StatusState>(null);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig, setFieldInConfigState } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { tags = [] } = instanceDetails;

  const { apiPath, maxLength, placeholder, configPath } = FIELD_PROPS_TAGS;

  let resetTimer = null;

  useEffect(() => {
    return () => {
      clearTimeout(resetTimer);
    };
  }, []);

  const resetStates = () => {
    setSubmitStatus(null);
    resetTimer = null;
    clearTimeout(resetTimer);
  };

  // posts all the tags at once as an array obj
  const postUpdateToAPI = async (postValue: any) => {
    setSubmitStatus(createInputStatus(STATUS_PROCESSING));

    await postConfigUpdateToAPI({
      apiPath,
      data: { value: postValue },
      onSuccess: () => {
        setFieldInConfigState({ fieldName: 'tags', value: postValue, path: configPath });
        setSubmitStatus(createInputStatus(STATUS_SUCCESS, 'Tags updated.'));
        setNewTagInput('');
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
      onError: (message: string) => {
        setSubmitStatus(createInputStatus(STATUS_ERROR, message));
        resetTimer = setTimeout(resetStates, RESET_TIMEOUT);
      },
    });
  };

  const handleInputChange = ({ value }: UpdateArgs) => {
    if (!submitStatus) {
      setSubmitStatus(null);
    }
    setNewTagInput(value);
  };

  // send to api and do stuff
  const handleSubmitNewTag = () => {
    resetStates();
    const newTag = newTagInput.trim();
    if (newTag === '') {
      setSubmitStatus(createInputStatus(STATUS_WARNING, 'Please enter a tag'));
      return;
    }
    if (tags.some(tag => tag.toLowerCase() === newTag.toLowerCase())) {
      setSubmitStatus(createInputStatus(STATUS_WARNING, 'This tag is already used!'));
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
  };

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
            <Tag closable onClose={handleClose} key={`tag-${tag}-${index}`}>
              {tag}
            </Tag>
          );
        })}
      </div>

      <div className="add-new-tag-section">
        <TextField
          fieldName="tag-input"
          value={newTagInput}
          className="new-tag-input"
          onChange={handleInputChange}
          onPressEnter={handleSubmitNewTag}
          maxLength={maxLength}
          placeholder={placeholder}
          status={submitStatus}
        />
      </div>
    </div>
  );
}
