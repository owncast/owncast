import { Spin } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useMemo, useRef } from 'react';
import { LoadingOutlined } from '@ant-design/icons';

import { MessageType, NameChangeEvent } from '../../../interfaces/socket-events';
import s from './ChatContainer.module.scss';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatUserMessage } from '..';

interface Props {
  messages: ChatMessage[];
  loading: boolean;
  usernameToHighlight: string;
  chatUserId: string;
  isModerator: boolean;
}

export default function ChatContainer(props: Props) {
  const { messages, loading, usernameToHighlight, chatUserId, isModerator } = props;

  const chatContainerRef = useRef(null);
  const spinIcon = <LoadingOutlined style={{ fontSize: '32px' }} spin />;

  const getNameChangeViewForMessage = (message: NameChangeEvent) => {
    const { oldName, user } = message;
    const { displayName, displayColor } = user;
    const color = `var(--theme-user-colors-${displayColor})`;

    return (
      <div>
        <span style={{ color }}>{oldName}</span> is now known as {displayName}
      </div>
    );
  };

  const getViewForMessage = message => {
    switch (message.type) {
      case MessageType.CHAT:
        return (
          <ChatUserMessage
            message={message}
            showModeratorMenu={isModerator} // Moderators have access to an additional menu
            highlightString={usernameToHighlight} // What to highlight in the message
            renderAsPersonallySent={message.user?.id === chatUserId} // The local user sent this message
          />
        );
      case MessageType.NAME_CHANGE:
        return getNameChangeViewForMessage(message);
      default:
        return null;
    }
  };

  const MessagesTable = useMemo(
    () => (
      <Virtuoso
        style={{ height: '80vh' }}
        ref={chatContainerRef}
        initialTopMostItemIndex={999}
        data={messages}
        itemContent={(index, message) => getViewForMessage(message)}
        followOutput="smooth"
      />
    ),
    [messages, usernameToHighlight, chatUserId, isModerator],
  );

  return (
    <div>
      <div className={s.chatHeader}>
        <span>stream chat</span>
      </div>
      <Spin spinning={loading} indicator={spinIcon}>
        {MessagesTable}
      </Spin>
    </div>
  );
}
