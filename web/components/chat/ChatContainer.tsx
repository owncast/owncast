import { Spin } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useRef } from 'react';
import { LoadingOutlined } from '@ant-design/icons';
import { ChatMessage } from '../../interfaces/chat-message.model';
import { ChatState } from '../../interfaces/application-state';
import ChatUserMessage from './ChatUserMessage';
import { MessageType } from '../../interfaces/socket-events';

interface Props {
  messages: ChatMessage[];
  state: ChatState;
}

export default function ChatContainer(props: Props) {
  const { messages, state } = props;
  const loading = state === ChatState.Loading;

  const chatContainerRef = useRef(null);
  const spinIcon = <LoadingOutlined style={{ fontSize: '32px' }} spin />;

  const getViewForMessage = message => {
    switch (message.type) {
      case MessageType.CHAT:
        return <ChatUserMessage message={message} showModeratorMenu={false} />;
      default:
        return null;
    }
  };

  return (
    <div>
      <h1>Chat</h1>
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
