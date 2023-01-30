// Custom component for AntDesign Button that makes an api call, then displays a confirmation icon upon
import React, { useState, useEffect, FC } from 'react';
import { Button, Tooltip } from 'antd';

import dynamic from 'next/dynamic';
import { fetchData, UPDATE_CHAT_MESSGAE_VIZ } from '../../utils/apis';
import { MessageType } from '../../types/chat';
import { isEmptyObject } from '../../utils/format';

// Lazy loaded components

const EyeOutlined = dynamic(() => import('@ant-design/icons/EyeOutlined'), {
  ssr: false,
});

const EyeInvisibleOutlined = dynamic(() => import('@ant-design/icons/EyeInvisibleOutlined'), {
  ssr: false,
});

const CheckCircleFilled = dynamic(() => import('@ant-design/icons/CheckCircleFilled'), {
  ssr: false,
});

const ExclamationCircleFilled = dynamic(() => import('@ant-design/icons/ExclamationCircleFilled'), {
  ssr: false,
});

export type MessageToggleProps = {
  isVisible: boolean;
  message: MessageType;
  setMessage: (message: MessageType) => void;
};

export const MessageVisiblityToggle: FC<MessageToggleProps> = ({
  isVisible,
  message,
  setMessage,
}) => {
  if (!message || isEmptyObject(message)) {
    return null;
  }

  let outcomeTimeout = null;
  const [outcome, setOutcome] = useState(0);

  const { id: messageId } = message || {};

  const resetOutcome = () => {
    outcomeTimeout = setTimeout(() => {
      setOutcome(0);
    }, 3000);
  };

  useEffect(() => () => {
    clearTimeout(outcomeTimeout);
  });

  const updateChatMessage = async () => {
    clearTimeout(outcomeTimeout);
    setOutcome(0);
    const result = await fetchData(UPDATE_CHAT_MESSGAE_VIZ, {
      auth: true,
      method: 'POST',
      data: {
        visible: !isVisible,
        idArray: [messageId],
      },
    });

    if (result.success && result.message === 'changed') {
      setMessage({ ...message, visible: !isVisible });
      setOutcome(1);
    } else {
      setMessage({ ...message, visible: isVisible });
      setOutcome(-1);
    }
    resetOutcome();
  };

  let outcomeIcon = <CheckCircleFilled style={{ color: 'transparent' }} />;
  if (outcome) {
    outcomeIcon =
      outcome > 0 ? (
        <CheckCircleFilled style={{ color: 'var(--ant-success)' }} />
      ) : (
        <ExclamationCircleFilled style={{ color: 'var(--ant-warning)' }} />
      );
  }

  const toolTipMessage = `Click to ${isVisible ? 'hide' : 'show'} this message`;
  return (
    <div className={`toggle-switch ${isVisible ? '' : 'hidden'}`}>
      <span className="outcome-icon">{outcomeIcon}</span>
      <Tooltip title={toolTipMessage} placement="topRight">
        <Button
          shape="circle"
          size="small"
          type="text"
          icon={isVisible ? <EyeOutlined /> : <EyeInvisibleOutlined />}
          onClick={updateChatMessage}
        />
      </Tooltip>
    </div>
  );
};
