import { Spin } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useRef } from 'react';
import { LoadingOutlined } from '@ant-design/icons';

import { MessageType, NameChangeEvent } from '../../../interfaces/socket-events';
import s from './ChatContainer.module.scss';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatUserMessage } from '..';
import ChatActionMessage from '../ChatActionMessage';

interface Props {
  messages: ChatMessage[];
  loading: boolean;
}

export default function ChatContainer(props: Props) {
  const { messages, loading } = props;

  const chatContainerRef = useRef(null);
  const spinIcon = <LoadingOutlined style={{ fontSize: '32px' }} spin />;

  const getNameChangeViewForMessage = (message: NameChangeEvent) => {
    const { oldName } = message;
    const { user } = message;
    const { displayName } = user;
    const body = `<strong>${oldName}</strong> is now known as <strong>${displayName}</strong>`;
    return <ChatActionMessage body={body} />;
  };

  const getViewForMessage = message => {
    switch (message.type) {
      case MessageType.CHAT:
        return <ChatUserMessage message={message} showModeratorMenu={false} />;
      case MessageType.NAME_CHANGE:
        return getNameChangeViewForMessage(message);
      default:
        return null;
    }
  };

  return (
    <div>
      <div className={s.chatHeader}>
        <span>stream chat</span>
      </div>
      <Spin spinning={loading} indicator={spinIcon} />
      <Virtuoso
        style={{ height: '80vh' }}
        ref={chatContainerRef}
        initialTopMostItemIndex={999}
        data={messages}
        itemContent={(index, message) => getViewForMessage(message)}
        followOutput="smooth"
      />
    </div>
  );
}
