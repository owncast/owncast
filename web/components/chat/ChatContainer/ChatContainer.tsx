import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useRef, CSSProperties, FC, useEffect } from 'react';
import dynamic from 'next/dynamic';
import {
  ConnectedClientInfoEvent,
  FediverseEvent,
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
import { ScrollToBotBtn } from './ScrollToBotBtn';
import { ChatActionMessage } from '../ChatActionMessage/ChatActionMessage';
import { ChatSocialMessage } from '../ChatSocialMessage/ChatSocialMessage';

// Lazy loaded components

const EditFilled = dynamic(() => import('@ant-design/icons/EditFilled'), {
  ssr: false,
});
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
  if (!message || !message.user) {
    return false;
  }

  const {
    user: { id },
  } = message;
  const lastMessage = messages[index - 1];
  if (lastMessage?.type !== MessageType.CHAT) {
    return false;
  }

  if (!lastMessage?.timestamp || !message.timestamp) {
    return false;
  }

  const maxTimestampDelta = 1000 * 60 * 2; // 2 minutes
  const lastTimestamp = new Date(lastMessage?.timestamp).getTime();
  const thisTimestamp = new Date(message.timestamp).getTime();
  if (thisTimestamp - lastTimestamp > maxTimestampDelta) {
    return false;
  }

  return id === lastMessage?.user.id;
}

function checkIsModerator(message: ChatMessage | ConnectedClientInfoEvent) {
  const {
    user: { scopes },
  } = message;

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
  const chatContainerRef = useRef(null);

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

  const getFediverseMessage = (message: FediverseEvent) => <ChatSocialMessage message={message} />;

  const getUserJoinedMessage = (message: ChatMessage) => {
    const {
      user: { displayName, displayColor },
    } = message;
    const isAuthorModerator = checkIsModerator(message);
    return (
      <ChatJoinMessage
        displayName={displayName}
        userColor={displayColor}
        isAuthorModerator={isAuthorModerator}
      />
    );
  };

  const getActionMessage = (message: ChatMessage) => {
    const { body } = message;
    return <ChatActionMessage body={body} />;
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
    message: ChatMessage | NameChangeEvent | ConnectedClientInfoEvent | FediverseEvent,
  ) => {
    switch (message.type) {
      case MessageType.CHAT:
        return (
          <ChatUserMessage
            message={message as ChatMessage}
            showModeratorMenu={isModerator} // Moderators have access to an additional menu
            highlightString={usernameToHighlight} // What to highlight in the message
            sentBySelf={(message as ChatMessage).user?.id === chatUserId} // The local user sent this message
            sameUserAsLast={shouldCollapseMessages(messages, index)}
            isAuthorModerator={(message as ChatMessage).user.scopes?.includes('MODERATOR')}
            isAuthorBot={(message as ChatMessage).user.scopes?.includes('BOT')}
            isAuthorAuthenticated={(message as ChatMessage).user?.authenticated}
            key={message.id}
          />
        );
      case MessageType.NAME_CHANGE:
        return getNameChangeViewForMessage(message as NameChangeEvent);
      case MessageType.CONNECTED_USER_INFO:
        return getConnectedInfoMessage(message as ConnectedClientInfoEvent);
      case MessageType.USER_JOINED:
        return getUserJoinedMessage(message as ChatMessage);
      case MessageType.CHAT_ACTION:
        return getActionMessage(message as ChatMessage);
      case MessageType.SYSTEM:
        return (
          <ChatSystemMessage
            message={message as ChatMessage}
            highlightString={usernameToHighlight} // What to highlight in the message
            key={message.id}
          />
        );
      case MessageType.FEDIVERSE_ENGAGEMENT_FOLLOW:
        return getFediverseMessage(message as FediverseEvent);
      case MessageType.FEDIVERSE_ENGAGEMENT_LIKE:
        return getFediverseMessage(message as FediverseEvent);
      case MessageType.FEDIVERSE_ENGAGEMENT_REPOST:
        return getFediverseMessage(message as FediverseEvent);

      default:
        return null;
    }
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const scrollChatToBottom = (ref, behavior = 'smooth') => {
    setTimeout(() => {
      ref.current?.scrollToIndex({
        index: messages.length - 1,
        behavior,
      });
    }, 100);

    setAtBottom(true);
  };

  // This is a hack to force a scroll to the very bottom of the chat messages
  // on initial mount of the component.
  // For https://github.com/owncast/owncast/issues/2500
  useEffect(() => {
    setTimeout(() => {
      scrollChatToBottom(chatContainerRef, 'auto');
    }, 500);
  }, []);

  const MessagesTable = useMemo(
    () => (
      <>
        <Virtuoso
          id="virtuoso"
          style={{ height }}
          className={styles.virtuoso}
          ref={chatContainerRef}
          data={messages}
          itemContent={(index, message) => getViewForMessage(index, message)}
          followOutput={(isAtBottom: boolean) => {
            if (isAtBottom) {
              scrollChatToBottom(chatContainerRef, 'smooth');
            }
            return false;
          }}
          alignToBottom
          atBottomThreshold={70}
          atBottomStateChange={bottom => {
            setAtBottom(bottom);
          }}
        />
        {!atBottom && <ScrollToBotBtn chatContainerRef={chatContainerRef} messages={messages} />}
      </>
    ),
    [messages, usernameToHighlight, chatUserId, isModerator, atBottom],
  );

  return (
    <div id="chat-container" className={styles.chatContainer}>
      {MessagesTable}
      {showInput && (
        <div className={styles.chatTextField}>
          <ChatTextField />
        </div>
      )}
    </div>
  );
};

ChatContainer.defaultProps = {
  showInput: true,
  height: 'auto',
};
