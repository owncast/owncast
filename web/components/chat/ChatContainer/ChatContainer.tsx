import { Virtuoso } from 'react-virtuoso';
import { useState, useMemo, useRef, CSSProperties, FC, useEffect } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
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
import { User } from '../../../interfaces/user.model';
import { ComponentError } from '../../ui/ComponentError/ComponentError';

export type ChatContainerProps = {
  messages: ChatMessage[];
  usernameToHighlight: string;
  chatUserId: string;
  isModerator: boolean;
  showInput?: boolean;
  height?: string;
  chatAvailable: boolean;
  focusInput?: boolean;
};

function shouldCollapseMessages(
  messages: ChatMessage[],
  index: number,
  collapsedMessageIds: Set<string>,
): boolean {
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

  const maxTimestampDelta = 1000 * 60; // 1 minute
  const lastTimestamp = new Date(lastMessage?.timestamp).getTime();
  const thisTimestamp = new Date(message.timestamp).getTime();
  if (thisTimestamp - lastTimestamp > maxTimestampDelta) {
    return false;
  }

  if (id !== lastMessage?.user.id) {
    return false;
  }

  // Limit the number of messages that can be collapsed in a row.
  const maxCollapsedMessageCount = 5;
  if (collapsedMessageIds.size >= maxCollapsedMessageCount) {
    return false;
  }

  return true;
}

function checkIsModerator(message: ChatMessage | ConnectedClientInfoEvent) {
  const { user } = message;

  const u = new User(user);
  return u.isModerator;
}

export const ChatContainer: FC<ChatContainerProps> = ({
  messages,
  usernameToHighlight,
  chatUserId,
  isModerator,
  showInput,
  height,
  chatAvailable: chatEnabled,
  focusInput = true,
}) => {
  const [showScrollToBottomButton, setShowScrollToBottomButton] = useState(false);
  const [isAtBottom, setIsAtBottom] = useState(false);

  const chatContainerRef = useRef(null);
  const scrollToBottomDelay = useRef(null);

  const collapsedMessageIds = new Set<string>();

  useEffect(
    () =>
      // Clear the timer when the component unmounts
      () => {
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

  const getUserChatMessageView = (index: number, message: ChatMessage) => {
    const collapsed = shouldCollapseMessages(messages, index, collapsedMessageIds);
    if (!collapsed) {
      collapsedMessageIds.clear();
    } else {
      collapsedMessageIds.add(message.id);
    }

    const isAuthorModerator = checkIsModerator(message);

    return (
      <ChatUserMessage
        message={message}
        showModeratorMenu={isModerator} // Moderators have access to an additional menu
        highlightString={usernameToHighlight} // What to highlight in the message
        sentBySelf={message.user?.id === chatUserId} // The local user sent this message
        sameUserAsLast={collapsed}
        isAuthorModerator={isAuthorModerator}
        isAuthorBot={message.user?.isBot}
        isAuthorAuthenticated={message.user?.authenticated}
        key={message.id}
      />
    );
  };
  const getViewForMessage = (
    index: number,
    message: ChatMessage | NameChangeEvent | ConnectedClientInfoEvent | FediverseEvent,
  ) => {
    switch (message.type) {
      case MessageType.CHAT:
        return getUserChatMessageView(index, message as ChatMessage);
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
    scrollToBottomDelay.current = setTimeout(() => {
      ref.current?.scrollToIndex({
        index: messages.length - 1,
        behavior,
      });
      setIsAtBottom(true);
      setShowScrollToBottomButton(false);
    }, 150);
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
          initialTopMostItemIndex={messages.length - 1}
          followOutput={() => {
            if (isAtBottom) {
              setShowScrollToBottomButton(false);
              scrollChatToBottom(chatContainerRef, 'auto');
              return 'smooth';
            }

            return false;
          }}
          alignToBottom
          atBottomThreshold={70}
          atBottomStateChange={bottom => {
            setIsAtBottom(bottom);

            if (bottom) {
              setShowScrollToBottomButton(false);
            } else {
              setShowScrollToBottomButton(true);
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
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="ChatContainer"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <div id="chat-container" className={styles.chatContainer}>
        {MessagesTable}
        {showInput && (
          <div className={styles.chatTextField}>
            <ChatTextField enabled={chatEnabled} focusInput={focusInput} />
          </div>
        )}
      </div>
    </ErrorBoundary>
  );
};

ChatContainer.defaultProps = {
  showInput: true,
  height: 'auto',
};
