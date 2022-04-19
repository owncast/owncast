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

const MAX_RENDER_BACKLOG = 300;

// Add message types that should be displayed in chat to this array.
const renderableChatStyleMessages = [
  SOCKET_MESSAGE_TYPES.NAME_CHANGE,
  SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO,
  SOCKET_MESSAGE_TYPES.USER_JOINED,
  SOCKET_MESSAGE_TYPES.CHAT_ACTION,
  SOCKET_MESSAGE_TYPES.SYSTEM,
  SOCKET_MESSAGE_TYPES.CHAT,
  SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_FOLLOW,
  SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_LIKE,
  SOCKET_MESSAGE_TYPES.FEDIVERSE_ENGAGEMENT_REPOST,
];
export default class Chat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      chatUserNames: [],
      // Ordered array of messages sorted by timestamp.
      sortedMessages: [],

      newMessagesReceived: false,
      webSocketConnected: true,
      isModerator: false,
    };

    this.scrollableMessagesContainer = createRef();

    this.websocket = null;
    this.receivedFirstMessages = false;
    this.receivedMessageUpdate = false;
    this.hasFetchedHistory = false;

    // Unordered dictionary of messages keyed by ID.
    this.messages = {};

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

    const {
      webSocketConnected,
      chatUserNames,
      newMessagesReceived,
      sortedMessages,
    } = this.state;

    const {
      webSocketConnected: nextSocket,
      chatUserNames: nextUserNames,
      newMessagesReceived: nextMessagesReceived,
    } = nextState;

    // If there are an updated number of sorted message then a render pass
    // needs to take place to render these new messages.
    if (
      Object.keys(sortedMessages).length !==
      Object.keys(nextState.sortedMessages).length
    ) {
      return true;
    }

    if (newMessagesReceived) {
      return true;
    }

    return (
      username !== nextUserName ||
      chatInputEnabled !== nextChatEnabled ||
      webSocketConnected !== nextSocket ||
      chatUserNames.length !== nextUserNames.length ||
      newMessagesReceived !== nextMessagesReceived
    );
  }

  componentDidUpdate(prevProps, prevState) {
    const { accessToken } = this.props;

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
  async getChatHistory(accessToken) {
    const { username } = this.props;
    try {
      const response = await fetch(
        URL_CHAT_HISTORY + `?accessToken=${accessToken}`
      );
      const data = await response.json();

      // Backlog of usernames from history
      const allChatUserNames = extraUserNamesFromMessageHistory(data);
      const chatUserNames = allChatUserNames.filter((name) => name != username);

      this.addNewRenderableMessages(data);

      this.setState((previousState) => {
        return {
          ...previousState,
          chatUserNames,
        };
      });
    } catch (error) {
      this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
    }

    jumpToBottom(this.scrollableMessagesContainer.current, 'instant');
  }

  receivedWebsocketMessage(message) {
    this.handleMessage(message);
  }

  handleNetworkingError(error) {
    // todo: something more useful
    console.error('chat error', error);
  }

  // Give a list of message IDs and the visibility state they should change to.
  updateMessagesVisibility(idsToUpdate, visible) {
    let messageList = { ...this.messages };

    // Iterate through each ID and mark the associated ID in our messages
    // dictionary with the new visibility.
    for (const id of idsToUpdate) {
      const message = messageList[id];
      if (message) {
        message.visible = visible;
        messageList[id] = message;
      }
    }

    const updatedMessagesList = {
      ...this.messages,
      ...messageList,
    };

    this.messages = updatedMessagesList;

    this.resortAndRenderMessages();
  }

  handleChangeModeratorStatus(isModerator) {
    if (isModerator !== this.state.isModerator) {
      this.setState((previousState) => {
        return { ...previousState, isModerator: isModerator };
      });
    }
  }

  handleWindowFocusNotificationCount(readonly, messageType) {
    // if window is blurred and we get a new message, add 1 to title
    if (
      !readonly &&
      messageType === SOCKET_MESSAGE_TYPES.CHAT &&
      this.windowBlurred
    ) {
      this.numMessagesSinceBlur += 1;
    }
  }

  addNewRenderableMessages(messagesArray) {
    // Convert the array of chat history messages into an object
    // to be merged with the existing chat messages.
    const newMessages = messagesArray.reduce(
      (o, message) => ({ ...o, [message.id]: message }),
      {}
    );

    // Keep our unsorted collection of messages keyed by ID.
    const updatedMessagesList = {
      ...newMessages,
      ...this.messages,
    };
    this.messages = updatedMessagesList;

    this.resortAndRenderMessages();
  }

  resortAndRenderMessages() {
    // Convert the unordered dictionary of messages to an ordered array.
    // NOTE: This sorts the entire collection of messages on every new message
    // because the order a message comes in cannot be trusted that it's the order
    // it was sent, you need to sort by timestamp. I don't know if there
    // is a  performance problem waiting to occur here for larger chat feeds.
    var sortedMessages = Object.values(this.messages)
      // Filter out messages set to not be visible
      .filter((message) => message.visible !== false)
      .sort((a, b) => {
        return Date.parse(a.timestamp) - Date.parse(b.timestamp);
      });

    // Cap this list to 300 items to improve browser performance.
    if (sortedMessages.length >= MAX_RENDER_BACKLOG) {
      sortedMessages = sortedMessages.slice(
        sortedMessages.length - MAX_RENDER_BACKLOG
      );
    }

    this.setState((previousState) => {
      return {
        ...previousState,
        newMessagesReceived: true,
        sortedMessages,
      };
    });
  }

  // handle any incoming message
  handleMessage(message) {
    const { type: messageType } = message;
    const { readonly, username } = this.props;

    // Allow non-user chat messages to be visible by default.
    const messageVisible =
      message.visible || messageType !== SOCKET_MESSAGE_TYPES.CHAT;

    // Show moderator status
    if (messageType === SOCKET_MESSAGE_TYPES.CONNECTED_USER_INFO) {
      const modStatusUpdate = checkIsModerator(message);
      this.handleChangeModeratorStatus(modStatusUpdate);
    }

    // Change the visibility of messages by ID.
    if (messageType === SOCKET_MESSAGE_TYPES.VISIBILITY_UPDATE) {
      const idsToUpdate = message.ids;
      const visible = message.visible;
      this.updateMessagesVisibility(idsToUpdate, visible);
    } else if (
      renderableChatStyleMessages.includes(messageType) &&
      messageVisible
    ) {
      // Add new message to the chat feed.
      this.addNewRenderableMessages([message]);

      // Update the usernames list, filtering out our own name.
      const updatedAllChatUserNames = this.updateAuthorList(message);
      if (updatedAllChatUserNames.length) {
        const updatedChatUserNames = updatedAllChatUserNames.filter(
          (name) => name != username
        );
        this.setState((previousState) => {
          return {
            ...previousState,
            chatUserNames: [...updatedChatUserNames],
          };
        });
      }
    }

    // Update the window title if needed.
    this.handleWindowFocusNotificationCount(readonly, messageType);
  }

  websocketConnected() {
    this.setState((previousState) => {
      return {
        ...previousState,
        webSocketConnected: true,
      };
    });
  }

  websocketDisconnected() {
    this.setState((previousState) => {
      return {
        ...previousState,
        webSocketConnected: false,
      };
    });
  }

  submitChat(content) {
    if (!content) {
      return;
    }
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

          this.setState((previousState) => {
            return {
              ...previousState,
              newMessagesReceived: false,
            };
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
    const { sortedMessages, chatUserNames, webSocketConnected, isModerator } =
      state;

    const messageList = sortedMessages.map(
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
