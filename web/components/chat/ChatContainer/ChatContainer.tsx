import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useRef, CSSProperties, FC, useEffect } from 'react';
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
import { ChatSystemMessage } from '../ChatSystemMessage/ChatSystemMessage';
import { ChatJoinMessage } from '../ChatJoinMessage/ChatJoinMessage';
import { ScrollToBotBtn } from './ScrollToBotBtn';
import { ChatActionMessage } from '../ChatActionMessage/ChatActionMessage';
import { ChatSocialMessage } from '../ChatSocialMessage/ChatSocialMessage';
import { ChatNameChangeMessage } from '../ChatNameChangeMessage/ChatNameChangeMessage';

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
  const [showScrollToBottomButton, setShowScrollToBottomButton] = useState(false);
  const [isAtBottom, setIsAtBottom] = useState(false);

  const chatContainerRef = useRef(null);
  const showScrollToBottomButtonDelay = useRef(null);
  const scrollToBottomDelay = useRef(null);

  const setShowScrolltoBottomButtonWithDelay = (show: boolean) => {
    showScrollToBottomButtonDelay.current = setTimeout(() => {
      setShowScrollToBottomButton(show);
    }, 1500);
  };

  useEffect(
    () =>
      // Clear the timer when the component unmounts
      () => {
        clearTimeout(showScrollToBottomButtonDelay.current);
        clearTimeout(scrollToBottomDelay.current);
      },
    [],
  );

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
            isAuthorModerator={(message as ChatMessage).user?.scopes?.includes('MODERATOR')}
            isAuthorBot={(message as ChatMessage).user?.scopes?.includes('BOT')}
            isAuthorAuthenticated={(message as ChatMessage).user?.authenticated}
            key={message.id}
          />
        );
      case MessageType.NAME_CHANGE:
        return <ChatNameChangeMessage message={message as NameChangeEvent} />;
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

  const scrollChatToBottom = (ref, behavior = 'smooth') => {
    clearTimeout(scrollToBottomDelay.current);
    clearTimeout(showScrollToBottomButtonDelay.current);
    scrollToBottomDelay.current = setTimeout(() => {
      ref.current?.scrollToIndex({
        index: messages.length - 1,
        behavior,
      });
      setIsAtBottom(true);
      setShowScrollToBottomButton(false);
    }, 100);
  };

  // This is a hack to force a scroll to the very bottom of the chat messages
  // on initial mount of the component.
  // For https://github.com/owncast/owncast/issues/2500
  useEffect(() => {
    setTimeout(() => {
      scrollChatToBottom(chatContainerRef, 'auto');
      setShowScrolltoBottomButtonWithDelay(false);
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
          initialTopMostItemIndex={messages.length - 1}
          followOutput={() => {
            clearTimeout(showScrollToBottomButtonDelay.current);

            if (isAtBottom) {
              setShowScrollToBottomButton(false);
              scrollChatToBottom(chatContainerRef, 'auto');
              return 'smooth';
            }
            setShowScrolltoBottomButtonWithDelay(true);

            return false;
          }}
          alignToBottom
          atBottomThreshold={70}
          atBottomStateChange={bottom => {
            setIsAtBottom(bottom);

            if (bottom) {
              setShowScrollToBottomButton(false);
            } else {
              setShowScrolltoBottomButtonWithDelay(true);
            }
          }}
        />
        {showScrollToBottomButton && (
          <ScrollToBotBtn
            onClick={() => {
              scrollChatToBottom(chatContainerRef, 'auto');
            }}
          />
        )}
      </>
    ),
    [messages, usernameToHighlight, chatUserId, isModerator, showScrollToBottomButton, isAtBottom],
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
