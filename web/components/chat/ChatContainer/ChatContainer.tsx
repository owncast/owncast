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
import { ChatPartMessage } from '../ChatPartMessage/ChatPartMessage';
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
  desktop?: boolean;
};

let resizeWindowCallback: () => void;

function shouldCollapseMessages(message: ChatMessage, previous: ChatMessage): boolean {
  if (!message || !message.user) {
    return false;
  }

  if (previous.type !== MessageType.CHAT) {
    return false;
  }

  const {
    user: { id },
  } = message;
  if (id !== previous.user.id) {
    return false;
  }

  if (!previous.timestamp || !message.timestamp) {
    return false;
  }

  const maxTimestampDelta = 1000 * 40; // 40 seconds
  const lastTimestamp = new Date(previous.timestamp).getTime();
  const thisTimestamp = new Date(message.timestamp).getTime();
  if (thisTimestamp - lastTimestamp > maxTimestampDelta) {
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
  desktop,
  focusInput = true,
}) => {
  const [showScrollToBottomButton, setShowScrollToBottomButton] = useState(false);
  const [isAtBottom, setIsAtBottom] = useState(false);

  const chatContainerRef = useRef(null);
  const scrollToBottomDelay = useRef(null);

  const collapsedIndexes: boolean[] = [];
  let consecutiveTally: number = 1;

  function calculateCollapsedMessages() {
    // Limits the number of messages that can be collapsed in a row.
    const maxCollapsedMessageCount = 5;
    for (let i = collapsedIndexes.length; i < messages.length; i += 1) {
      const collapse: boolean =
        i > 0 &&
        consecutiveTally < maxCollapsedMessageCount &&
        shouldCollapseMessages(messages[i], messages[i - 1]);
      collapsedIndexes.push(collapse);
      consecutiveTally = 1 + (collapse ? consecutiveTally : 0);
    }
  }

  function shouldCollapse(index: number): boolean {
    if (collapsedIndexes.length <= index) {
      calculateCollapsedMessages();
    }
    return collapsedIndexes[index];
  }

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

  const getUserPartMessage = (message: ChatMessage) => {
    const {
      user: { displayName, displayColor },
    } = message;
    const isAuthorModerator = checkIsModerator(message);
    return (
      <ChatPartMessage
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
    const isAuthorModerator = checkIsModerator(message);

    return (
      <ChatUserMessage
        message={message}
        showModeratorMenu={isModerator} // Moderators have access to an additional menu
        highlightString={usernameToHighlight} // What to highlight in the message
        sentBySelf={message.user?.id === chatUserId} // The local user sent this message
        sameUserAsLast={shouldCollapse(index)}
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
      case MessageType.USER_PARTED:
        return getUserPartMessage(message as ChatMessage);
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

  const scrollChatToBottom = ref => {
    clearTimeout(scrollToBottomDelay.current);
    scrollToBottomDelay.current = setTimeout(() => {
      ref.current?.scrollTo({ top: Infinity, left: 0, behavior: 'auto' });

      setIsAtBottom(true);
    }, 150);
    setShowScrollToBottomButton(false);
  };

  // This is a hack to force a scroll to the very bottom of the chat messages
  // on initial mount of the component.
  // For https://github.com/owncast/owncast/issues/2500
  useEffect(() => {
    setTimeout(() => {
      scrollChatToBottom(chatContainerRef);
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
              scrollChatToBottom(chatContainerRef);
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
              scrollChatToBottom(chatContainerRef);
            }}
          />
        )}
      </>
    ),
    [messages, usernameToHighlight, chatUserId, isModerator, showScrollToBottomButton, isAtBottom],
  );

  const defaultChatWidth: number = 320;
  function clampChatWidth(desired) {
    return Math.max(200, Math.min(window.innerWidth * 0.666, desired));
  }

  function startDrag(dragEvent) {
    const container = document.getElementById('chat-container');
    function move(event) {
      container.style.width = `${clampChatWidth(window.innerWidth - event.x)}px`;
    }
    function endDrag() {
      window.document.removeEventListener('mousemove', move);
      window.document.removeEventListener('mouseup', endDrag);
      window.document.removeEventListener('focusout', endDrag);
    }
    window.document.addEventListener('mousemove', move);
    window.document.addEventListener('mouseup', endDrag);
    window.document.addEventListener('focusout', endDrag);
    dragEvent.preventDefault(); // Prevent selecting the page as you resize it
  }

  // Re-clamp the chat size whenever the window resizes
  function resize() {
    const container = desktop && document.getElementById('chat-container');
    if (container) {
      const currentWidth = parseFloat(container.style.width) || defaultChatWidth;
      container.style.width = `${clampChatWidth(currentWidth)}px`;
    }
  }

  if (resizeWindowCallback) window.removeEventListener('resize', resizeWindowCallback);
  if (desktop) {
    window.addEventListener('resize', resize);
    resizeWindowCallback = resize;
  } else {
    resizeWindowCallback = null;
  }

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
      <div
        id="chat-container"
        className={styles.chatContainer}
        style={desktop && { width: `${defaultChatWidth}px` }}
      >
        {MessagesTable}
        {showInput && (
          <div className={styles.chatTextField}>
            <ChatTextField enabled={chatEnabled} focusInput={focusInput} />
          </div>
        )}
        {desktop && (
          <div className={styles.resizeHandle} onMouseDown={startDrag} role="presentation" />
        )}
      </div>
    </ErrorBoundary>
  );
};

ChatContainer.defaultProps = {
  showInput: true,
  height: 'auto',
};
