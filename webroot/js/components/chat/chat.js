import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import Message from './message.js';
import ChatInput from './chat-input.js';
import { CALLBACKS, SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import { jumpToBottom, debounce } from '../../utils/helpers.js';
import {
  extraUserNamesFromMessageHistory,
  checkIsModerator,
} from '../../utils/chat.js';
import {
  URL_CHAT_HISTORY,
  MESSAGE_JUMPTOBOTTOM_BUFFER,
} from '../../utils/constants.js';

// Add message types that should be displayed in chat to this array.
const renderableChatStyleMessages = [
  SOCKET_MESSAGE_TYPES.NAME_CHANGE,
  SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO,
  SOCKET_MESSAGE_TYPES.USER_JOINED,
  SOCKET_MESSAGE_TYPES.SYSTEM,
  SOCKET_MESSAGE_TYPES.CHAT,
];
export default class Chat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      chatUserNames: [],
      messages: [],
      newMessagesReceived: false,
      webSocketConnected: true,
      isModerator: false,
    };

    this.scrollableMessagesContainer = createRef();

    this.websocket = null;
    this.receivedFirstMessages = false;
    this.receivedMessageUpdate = false;
    this.hasFetchedHistory = false;
    this.forceRender = false;

    this.windowBlurred = false;
    this.numMessagesSinceBlur = 0;

    this.getChatHistory = this.getChatHistory.bind(this);
    this.handleNetworkingError = this.handleNetworkingError.bind(this);
    this.handleWindowBlur = this.handleWindowBlur.bind(this);
    this.handleWindowFocus = this.handleWindowFocus.bind(this);
    this.handleWindowResize = debounce(this.handleWindowResize.bind(this), 500);
    this.messageListCallback = this.messageListCallback.bind(this);
    this.receivedWebsocketMessage = this.receivedWebsocketMessage.bind(this);
    this.scrollToBottom = this.scrollToBottom.bind(this);
    this.submitChat = this.submitChat.bind(this);
    this.websocketConnected = this.websocketConnected.bind(this);
    this.websocketDisconnected = this.websocketDisconnected.bind(this);
  }

  componentDidMount() {
    this.setupWebSocketCallbacks();

    window.addEventListener('resize', this.handleWindowResize);

    if (!this.props.readonly) {
      window.addEventListener('blur', this.handleWindowBlur);
      window.addEventListener('focus', this.handleWindowFocus);
    }

    this.messageListObserver = new MutationObserver(this.messageListCallback);
    this.messageListObserver.observe(this.scrollableMessagesContainer.current, {
      childList: true,
    });
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { username, chatInputEnabled } = this.props;
    const { username: nextUserName, chatInputEnabled: nextChatEnabled } =
      nextProps;

    const { webSocketConnected, messages, chatUserNames, newMessagesReceived } =
      this.state;

    if (this.forceRender) {
      return true;
    }

    const {
      webSocketConnected: nextSocket,
      messages: nextMessages,
      chatUserNames: nextUserNames,
      newMessagesReceived: nextMessagesReceived,
    } = nextState;

    return (
      username !== nextUserName ||
      chatInputEnabled !== nextChatEnabled ||
      webSocketConnected !== nextSocket ||
      messages.length !== nextMessages.length ||
      chatUserNames.length !== nextUserNames.length ||
      newMessagesReceived !== nextMessagesReceived
    );
  }

  componentDidUpdate(prevProps, prevState) {
    const { username: prevName } = prevProps;
    const { username, accessToken } = this.props;

    const { messages: prevMessages } = prevState;
    const { messages } = this.state;

    // scroll to bottom of messages list when new ones come in
    if (messages.length !== prevMessages.length) {
      this.setState({
        newMessagesReceived: true,
      });
    }

    // Fetch chat history
    if (!this.hasFetchedHistory && accessToken) {
      this.hasFetchedHistory = true;
      this.getChatHistory(accessToken);
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleWindowResize);
    if (!this.props.readonly) {
      window.removeEventListener('blur', this.handleWindowBlur);
      window.removeEventListener('focus', this.handleWindowFocus);
    }
    this.messageListObserver.disconnect();
  }

  setupWebSocketCallbacks() {
    this.websocket = this.props.websocket;
    if (this.websocket) {
      this.websocket.addListener(
        CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED,
        this.receivedWebsocketMessage
      );
      this.websocket.addListener(
        CALLBACKS.WEBSOCKET_CONNECTED,
        this.websocketConnected
      );
      this.websocket.addListener(
        CALLBACKS.WEBSOCKET_DISCONNECTED,
        this.websocketDisconnected
      );
    }
  }

  // fetch chat history
  getChatHistory(accessToken) {
    const { username } = this.props;
    fetch(URL_CHAT_HISTORY + `?accessToken=${accessToken}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then((data) => {
        // extra user names
        const allChatUserNames = extraUserNamesFromMessageHistory(data);
        const chatUserNames = allChatUserNames.filter(
          (name) => name != username
        );
        this.setState((previousState, currentProps) => {
          return {
            ...previousState,
            messages: data.concat(previousState.messages),
            chatUserNames,
          };
        });

        this.scrollToBottom();
      })
      .catch((error) => {
        this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
      });
  }

  receivedWebsocketMessage(message) {
    this.handleMessage(message);
  }

  handleNetworkingError(error) {
    // todo: something more useful
    console.error('chat error', error);
  }

  // handle any incoming message
  handleMessage(message) {
    const {
      id: messageId,
      type: messageType,
      timestamp: messageTimestamp,
    } = message;
    const { messages: curMessages } = this.state;
    const { username, readonly } = this.props;

    const existingIndex = curMessages.findIndex(
      (item) => item.id === messageId
    );

    // Allow non-user chat messages to be visible by default.
    const messageVisible =
      message.visible || messageType !== SOCKET_MESSAGE_TYPES.CHAT;

    // check moderator status
    if (messageType === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
      const modStatusUpdate = checkIsModerator(message);
      if (modStatusUpdate !== this.state.isModerator) {
        this.setState((previousState, currentProps) => {
          return { ...previousState, isModerator: modStatusUpdate };
        });
      }
    }

    const updatedMessageList = [...curMessages];

    // Change the visibility of messages by ID.
    if (messageType === 'VISIBILITY-UPDATE') {
      const idsToUpdate = message.ids;
      const visible = message.visible;
      updatedMessageList.forEach((item) => {
        if (idsToUpdate.includes(item.id)) {
          item.visible = visible;
        }
      });
      this.forceRender = true;
    } else if (
      renderableChatStyleMessages.includes(messageType) &&
      existingIndex === -1 &&
      messageVisible
    ) {
      // insert message at timestamp
      const convertedMessage = {
        ...message,
      };
      const insertAtIndex = curMessages.findIndex((item, index) => {
        const time = item.timestamp || messageTimestamp;
        const nextMessage =
          index < curMessages.length - 1 && curMessages[index + 1];
        const nextTime = nextMessage.timestamp || messageTimestamp;
        const messageTimestampDate = new Date(messageTimestamp);
        return (
          messageTimestampDate > new Date(time) &&
          messageTimestampDate <= new Date(nextTime)
        );
      });
      updatedMessageList.splice(insertAtIndex + 1, 0, convertedMessage);
      if (updatedMessageList.length > 300) {
        updatedMessageList = updatedMessageList.slice(
          Math.max(updatedMessageList.length - 300, 0)
        );
      }
      this.setState((previousState, currentProps) => {
        return { ...previousState, messages: updatedMessageList };
      });
    } else if (
      renderableChatStyleMessages.includes(messageType) &&
      existingIndex === -1
    ) {
      // else if message doesn't exist, add it and extra username
      const newState = {
        messages: [...curMessages, message],
      };
      const updatedAllChatUserNames = this.updateAuthorList(message);
      if (updatedAllChatUserNames.length) {
        const updatedChatUserNames = updatedAllChatUserNames.filter(
          (name) => name != username
        );
        newState.chatUserNames = [...updatedChatUserNames];
      }

      this.setState((previousState, currentProps) => {
        return { ...previousState, newState };
      });
    }

    // if window is blurred and we get a new message, add 1 to title
    if (
      !readonly &&
      messageType === SOCKET_MESSAGE_TYPES.CHAT &&
      this.windowBlurred
    ) {
      this.numMessagesSinceBlur += 1;
    }
  }

  websocketConnected() {
    this.setState({
      webSocketConnected: true,
    });
  }

  websocketDisconnected() {
    this.setState({
      webSocketConnected: false,
    });
  }

  submitChat(content) {
    if (!content) {
      return;
    }
    const { username } = this.props;
    const message = {
      body: content,
      type: SOCKET_MESSAGE_TYPES.CHAT,
    };
    this.websocket.send(message);
  }

  updateAuthorList(message) {
    const { type } = message;
    let nameList = this.state.chatUserNames;

    if (
      type === SOCKET_MESSAGE_TYPES.CHAT &&
      !nameList.includes(message.user.displayName)
    ) {
      nameList.push(message.user.displayName);
      return nameList;
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, user } = message;
      const oldNameIndex = nameList.indexOf(oldName);
      nameList.splice(oldNameIndex, 1, user.displayName);
      return nameList;
    }
    return [];
  }

  scrollToBottom() {
    jumpToBottom(this.scrollableMessagesContainer.current);
  }

  checkShouldScroll() {
    const { scrollTop, scrollHeight, clientHeight } =
      this.scrollableMessagesContainer.current;
    const fullyScrolled = scrollHeight - clientHeight;
    const shouldScroll =
      scrollHeight >= clientHeight &&
      fullyScrolled - scrollTop < MESSAGE_JUMPTOBOTTOM_BUFFER;

    return shouldScroll;
  }

  handleWindowResize() {
    this.scrollToBottom();
  }

  handleWindowBlur() {
    this.windowBlurred = true;
  }

  handleWindowFocus() {
    this.windowBlurred = false;
    this.numMessagesSinceBlur = 0;
    window.document.title = this.props.instanceTitle;
  }

  // if the messages list grows in number of child message nodes due to new messages received, scroll to bottom.
  messageListCallback(mutations) {
    const numMutations = mutations.length;

    if (numMutations) {
      const item = mutations[numMutations - 1];
      if (item.type === 'childList' && item.addedNodes.length) {
        if (this.state.newMessagesReceived) {
          if (!this.receivedFirstMessages) {
            this.scrollToBottom();
            this.receivedFirstMessages = true;
          } else if (this.checkShouldScroll()) {
            this.scrollToBottom();
          }
          this.setState({
            newMessagesReceived: false,
          });
        }
      }
      // update document title if window blurred
      if (
        this.numMessagesSinceBlur &&
        !this.props.readonly &&
        this.windowBlurred
      ) {
        this.updateDocumentTitle();
      }
    }
  }

  updateDocumentTitle() {
    const num =
      this.numMessagesSinceBlur > 10 ? '10+' : this.numMessagesSinceBlur;
    window.document.title = `${num} 💬 :: ${this.props.instanceTitle}`;
  }

  render(props, state) {
    const { username, readonly, chatInputEnabled, inputMaxBytes, accessToken } =
      props;
    const { messages, chatUserNames, webSocketConnected, isModerator } = state;

    const messageList = messages
      .filter((message) => message.visible !== false)
      .map(
        (message) =>
          html`<${Message}
            message=${message}
            username=${username}
            key=${message.id}
            isModerator=${isModerator}
            accessToken=${accessToken}
          />`
      );

    if (readonly) {
      return html`
        <div
          id="messages-container"
          ref=${this.scrollableMessagesContainer}
          class="scrollbar-hidden py-1 overflow-auto"
        >
          ${messageList}
        </div>
      `;
    }

    return html`
      <section id="chat-container-wrap" class="flex flex-col">
        <div
          id="chat-container"
          class="bg-gray-800 flex flex-col justify-end overflow-auto"
        >
          <div
            id="messages-container"
            ref=${this.scrollableMessagesContainer}
            class="scrollbar-hidden py-1 overflow-auto z-10"
          >
            ${messageList}
          </div>
          <${ChatInput}
            chatUserNames=${chatUserNames}
            inputEnabled=${webSocketConnected && chatInputEnabled}
            handleSendMessage=${this.submitChat}
            inputMaxBytes=${inputMaxBytes}
          />
        </div>
      </section>
    `;
  }
}
