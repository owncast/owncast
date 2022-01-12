import React, { useState } from 'react';

import { Button, Space, Input, Modal } from 'antd';
import { STATUS_ERROR, STATUS_SUCCESS } from '../utils/input-statuses';
import { fetchData, FEDERATION_MESSAGE_SEND } from '../utils/apis';

const { TextArea } = Input;

interface ComposeFederatedPostProps {
  visible: boolean;
  handleClose: () => void;
}

export default function ComposeFederatedPost({ visible, handleClose }: ComposeFederatedPostProps) {
  const [content, setContent] = useState('');
  const [postPending, setPostPending] = useState(false);
  const [postSuccessState, setPostSuccessState] = useState(null);

  function handleEditorChange(e) {
    setContent(e.target.value);
  }

  async function sendButtonClicked() {
    setPostPending(true);

    const data = {
      value: content,
    };
    try {
      await fetchData(FEDERATION_MESSAGE_SEND, {
        data,
        method: 'POST',
        auth: true,
      });
      setPostSuccessState(STATUS_SUCCESS);
      setTimeout(handleClose, 1000);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error(e);
      setPostSuccessState(STATUS_ERROR);
    }
    setPostPending(false);
  }

  return (
    <Modal
      destroyOnClose
      width={600}
      title="Post to Followers"
      visible={visible}
      onCancel={handleClose}
      footer={[
        <Button onClick={() => handleClose()}>Cancel</Button>,
        <Button
          type="primary"
          onClick={sendButtonClicked}
          disabled={postPending || postSuccessState}
          loading={postPending}
        >
          {postSuccessState?.toUpperCase() || 'Post'}
        </Button>,
      ]}
    >
      <Space id="fediverse-post-container" direction="vertical">
        <TextArea
          placeholder="Tell the world about your streaming plans..."
          size="large"
          showCount
          maxLength={500}
          style={{ height: '150px' }}
          onChange={handleEditorChange}
        />
      </Space>
    </Modal>
  );
}
