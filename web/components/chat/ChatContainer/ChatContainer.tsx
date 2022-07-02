import { Button, Spin } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useRef } from 'react';
import { EditFilled, LoadingOutlined, VerticalAlignBottomOutlined } from '@ant-design/icons';
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

  const [atBottom, setAtBottom] = useState(false);
  // const [showButton, setShowButton] = useState(false);
  const chatContainerRef = useRef(null);
  const spinIcon = <LoadingOutlined style={{ fontSize: '32px' }} spin />;

  const getNameChangeViewForMessage = (message: NameChangeEvent) => {
    const { oldName, user } = message;
    const { displayName, displayColor } = user;
    const color = `var(--theme-user-colors-${displayColor})`;

    return (
      <div className={s.nameChangeView}>
        <div style={{ marginRight: 5, height: 'max-content', margin: 'auto 5px auto 0' }}>
          <EditFilled />
        </div>
        <div className={s.nameChangeText}>
          <span style={{ color }}>{oldName}</span>
          <span className={s.plain}> is now known as </span>
          <span style={{ color }}>{displayName}</span>
        </div>
      </div>
    );
  };

  const getViewForMessage = (index: number, message: ChatMessage | NameChangeEvent) => {
    switch (message.type) {
      case MessageType.CHAT:
        return (
          <ChatUserMessage
            message={message as ChatMessage}
            showModeratorMenu={isModerator} // Moderators have access to an additional menu
            highlightString={usernameToHighlight} // What to highlight in the message
            sentBySelf={message.user?.id === chatUserId} // The local user sent this message
            sameUserAsLast={isSameUserAsLast(messages, index)}
            key={message.id}
          />
        );
      case MessageType.NAME_CHANGE:
        return getNameChangeViewForMessage(message as NameChangeEvent);
      default:
        return null;
    }
  };

  const MessagesTable = useMemo(
    () => (
      <>
        <Virtuoso
          style={{ height: '75vh' }}
          ref={chatContainerRef}
          initialTopMostItemIndex={messages.length - 1} // Force alignment to bottom
          data={messages}
          itemContent={(index, message) => getViewForMessage(index, message)}
          followOutput="auto"
          alignToBottom
          atBottomStateChange={bottom => setAtBottom(bottom)}
        />
        {!atBottom && (
          <div className={s.toBottomWrap}>
            <Button
              type="default"
              icon={<VerticalAlignBottomOutlined />}
              onClick={() =>
                chatContainerRef.current.scrollToIndex({
                  index: messages.length - 1,
                  behavior: 'smooth',
                })
              }
            >
              Go to last message
            </Button>
          </div>
        )}
      </>
    ),
    [messages, usernameToHighlight, chatUserId, isModerator, atBottom],
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

function isSameUserAsLast(messages: ChatMessage[], index: number) {
  const message = messages[index];
  const {
    user: { id },
  } = message;
  const lastMessage = messages[index - 1];

  return id === lastMessage?.user.id;
}
