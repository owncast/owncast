import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import Message from './message.js';
import ChatInput from './chat-input.js';
import { CALLBACKS, SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import {
  jumpToBottom,
  debounce,
  setLocalStorage,
} from '../../utils/helpers.js';
import { extraUserNamesFromMessageHistory } from '../../utils/chat.js';
import {
  URL_CHAT_HISTORY,
  MESSAGE_JUMPTOBOTTOM_BUFFER,
} from '../../utils/constants.js';

export default class Chat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      chatUserNames: [],
      messages: [],
      newMessagesReceived: false,
      webSocketConnected: true,
    };

    this.scrollableMessagesContainer = createRef();

    this.websocket = null;
    this.receivedFirstMessages = false;
    this.receivedMessageUpdate = false;
    this.hasFetchedHistory = false;

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

    if (!this.props.messagesOnly) {
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
    if (!this.props.messagesOnly) {
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
    fetch(URL_CHAT_HISTORY + `?accessToken=${accessToken}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then((data) => {
        // extra user names
        const chatUserNames = extraUserNamesFromMessageHistory(data);
        this.setState({
          messages: this.state.messages.concat(data),
          chatUserNames,
        });
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
      visible: messageVisible,
    } = message;
    const { messages: curMessages } = this.state;
    const { messagesOnly } = this.props;

    const existingIndex = curMessages.findIndex(
      (item) => item.id === messageId
    );

    // If the message already exists and this is an update event
    // then update it.
    if (messageType === 'VISIBILITY-UPDATE') {
      const updatedMessageList = [...curMessages];
      const convertedMessage = {
        ...message,
        type: 'CHAT',
      };
      // if message exists and should now hide, take it out.
      if (existingIndex >= 0 && !messageVisible) {
        this.setState({
          messages: curMessages.filter((item) => item.id !== messageId),
        });
      } else if (existingIndex === -1 && messageVisible) {
        // insert message at timestamp
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
        this.setState({
          messages: updatedMessageList,
        });
      }
    } else if (existingIndex === -1) {
      // else if message doesn't exist, add it and extra username
      const newState = {
        messages: [...curMessages, message],
      };
      const updatedChatUserNames = this.updateAuthorList(message);
      if (updatedChatUserNames.length) {
        newState.chatUserNames = [...updatedChatUserNames];
      }
      this.setState(newState);
    }

    // if window is blurred and we get a new message, add 1 to title
    if (!messagesOnly && messageType === 'CHAT' && this.windowBlurred) {
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
    const nameList = this.state.chatUserNames;

    if (
      type === SOCKET_MESSAGE_TYPES.CHAT &&
      !nameList.includes(message.user.displayName)
    ) {
      return nameList.push(message.user.displayName);
    } else if (type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
      const { oldName, newName } = message;
      const oldNameIndex = nameList.indexOf(oldName);
      return nameList.splice(oldNameIndex, 1, newName);
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
        !this.props.messagesOnly &&
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
    const { username, messagesOnly, chatInputEnabled } = props;
    const { messages, chatUserNames, webSocketConnected } = state;

    const messageList = messages
      .filter((message) => message.visible !== false)
      .map(
        (message) =>
          html`<${Message}
            message=${message}
            username=${username}
            key=${message.id}
          />`
      );

    if (messagesOnly) {
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
          />
        </div>
      </section>
    `;
  }
}
