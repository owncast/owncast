import { Spin } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useCallback, useEffect, useRef } from 'react';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { ChatState } from '../../interfaces/application-state';
import ChatUserMessage from './ChatUserMessage';

interface Props {
  messages: ChatMessage[];
  state: ChatState;
}

export default function ChatContainer(props: Props) {
  const { messages, state } = props;
  const loading = state === ChatState.Loading;

  const chatContainerRef = useRef(null);

  return (
    <div>
      <Spin spinning={loading} />

      <Virtuoso
        style={{ height: '400px' }}
        ref={chatContainerRef}
        initialTopMostItemIndex={999}
        data={messages}
        itemContent={(index, message) => (
          <ChatUserMessage message={message} showModeratorMenu={false} />
        )}
        followOutput="smooth"
      />
    </div>
  );
}
