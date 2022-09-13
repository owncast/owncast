import { Button } from 'antd';
import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useRef, CSSProperties, FC } from 'react';
import { EditFilled, VerticalAlignBottomOutlined } from '@ant-design/icons';
import {
  ConnectedClientInfoEvent,
  MessageType,
  NameChangeEvent,
} from '../../../interfaces/socket-events';
import styles from './ChatContainer.module.scss';
import { ChatMessage } from '../../../interfaces/chat-message.model';
import { ChatUserMessage } from '../ChatUserMessage/ChatUserMessage';
import { ChatTextField } from '../ChatTextField/ChatTextField';
import { ChatModeratorNotification } from '../ChatModeratorNotification/ChatModeratorNotification';
// import ChatActionMessage from '../ChatAction/ChatActionMessage';
import { ChatSystemMessage } from '../ChatSystemMessage/ChatSystemMessage';
import { ChatJoinMessage } from '../ChatJoinMessage/ChatJoinMessage';

export type ChatContainerProps = {
  messages: ChatMessage[];
  usernameToHighlight: string;
  chatUserId: string;
  isModerator: boolean;
  showInput?: boolean;
  height?: string;
};

function shouldCollapseMessages(messages: ChatMessage[], index: number): boolean {
  if (messages.length < 2) {
    return false;
  }

  const message = messages[index];
  const {
    user: { id },
  } = message;
  const lastMessage = messages[index - 1];
  if (lastMessage?.type !== MessageType.CHAT) {
    return false;
  }

  if (!lastMessage.timestamp || !message.timestamp) {
    return false;
  }

  const maxTimestampDelta = 1000 * 60 * 2; // 2 minutes
  const lastTimestamp = new Date(lastMessage.timestamp).getTime();
  const thisTimestamp = new Date(message.timestamp).getTime();
  if (thisTimestamp - lastTimestamp > maxTimestampDelta) {
    return false;
  }

  return id === lastMessage?.user.id;
}

function checkIsModerator(message) {
  const { user } = message;
  const { scopes } = user;

  if (!scopes || scopes.length === 0) {
    return false;
  }

  return scopes.includes('MODERATOR');
}

export const ChatContainer: FC<ChatContainerProps> = ({
  messages,
  usernameToHighlight,
  chatUserId,
  isModerator,
  showInput,
  height,
}) => {
  const [atBottom, setAtBottom] = useState(false);
  // const [showButton, setShowButton] = useState(false);
  const chatContainerRef = useRef(null);
  // const spinIcon = <LoadingOutlined style={{ fontSize: '32px' }} spin />;

  const getNameChangeViewForMessage = (message: NameChangeEvent) => {
    const { oldName, user } = message;
    const { displayName, displayColor } = user;
    const color = `var(--theme-color-users-${displayColor})`;

    return (
      <div className={styles.nameChangeView}>
        <div style={{ marginRight: 5, height: 'max-content', margin: 'auto 5px auto 0' }}>
          <EditFilled />
        </div>
        <div className={styles.nameChangeText}>
          <span style={{ color }}>{oldName}</span>
          <span className={styles.plain}> is now known as </span>
          <span style={{ color }}>{displayName}</span>
        </div>
      </div>
    );
  };

  const getUserJoinedMessage = (message: ChatMessage) => {
    const { user } = message;
    const { displayName, displayColor } = user;
    const isAuthorModerator = checkIsModerator(message);
    return (
      <ChatJoinMessage
        displayName={displayName}
        userColor={displayColor}
        isAuthorModerator={isAuthorModerator}
      />
    );
  };

  const getConnectedInfoMessage = (message: ConnectedClientInfoEvent) => {
    const modStatusUpdate = checkIsModerator(message);
    if (!modStatusUpdate) {
      // Important note: We can't return null or an element with zero width
      // or zero height. So to work around this we return a very small 1x1 div.
      const st: CSSProperties = { width: '1px', height: '1px' };
      return <div style={st} />;
    }

    // Alert the user that they are a moderator.
    return <ChatModeratorNotification />;
  };

  const getViewForMessage = (
    index: number,
    message: ChatMessage | NameChangeEvent | ConnectedClientInfoEvent,
  ) => {
    switch (message.type) {
      case MessageType.CHAT:
        return (
          <ChatUserMessage
            message={message as ChatMessage}
            showModeratorMenu={isModerator} // Moderators have access to an additional menu
            highlightString={usernameToHighlight} // What to highlight in the message
            sentBySelf={message.user?.id === chatUserId} // The local user sent this message
            sameUserAsLast={shouldCollapseMessages(messages, index)}
            isAuthorModerator={(message as ChatMessage).user.scopes?.includes('MODERATOR')}
            isAuthorAuthenticated={message.user?.authenticated}
            key={message.id}
          />
        );
      case MessageType.NAME_CHANGE:
        return getNameChangeViewForMessage(message as NameChangeEvent);
      case MessageType.CONNECTED_USER_INFO:
        return getConnectedInfoMessage(message);
      case MessageType.USER_JOINED:
        return getUserJoinedMessage(message as ChatMessage);
      case MessageType.SYSTEM:
        return (
          <ChatSystemMessage
            message={message as ChatMessage}
            highlightString={usernameToHighlight} // What to highlight in the message
            key={message.id}
          />
        );

      default:
        return null;
    }
  };

  const MessagesTable = useMemo(
    () => (
      <>
        <Virtuoso
          style={{ height }}
          className={styles.virtuoso}
          ref={chatContainerRef}
          initialTopMostItemIndex={messages.length - 1} // Force alignment to bottom
          data={messages}
          itemContent={(index, message) => getViewForMessage(index, message)}
          followOutput="auto"
          alignToBottom
          atBottomStateChange={bottom => setAtBottom(bottom)}
        />
        {!atBottom && (
          <div className={styles.toBottomWrap}>
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
    <div className={styles.chatContainer}>
      {MessagesTable}
      {showInput && <ChatTextField />}
    </div>
  );
};

ChatContainer.defaultProps = {
  showInput: true,
  height: 'auto',
};
