import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import '/js/web_modules/@justinribeiro/lite-youtube.js';

import Message from './message.js';
import ChatInput from './chat-input.js';
import { CALLBACKS, SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import { jumpToBottom, debounce } from '../../utils/helpers.js';
import { extraUserNamesFromMessageHistory } from '../../utils/chat.js';
import { URL_CHAT_HISTORY, MESSAGE_JUMPTOBOTTOM_BUFFER } from '../../utils/constants.js';

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
   this.getChatHistory();

   window.addEventListener('resize', this.handleWindowResize);

   if (!this.props.messagesOnly) {
    window.addEventListener('blur', this.handleWindowBlur);
    window.addEventListener('focus', this.handleWindowFocus);
   }

   this.messageListObserver = new MutationObserver(this.messageListCallback);
   this.messageListObserver.observe(this.scrollableMessagesContainer.current, { childList: true });
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { username, chatInputEnabled } = this.props;
    const { username: nextUserName, chatInputEnabled: nextChatEnabled } = nextProps;

    const { webSocketConnected, messages, chatUserNames, newMessagesReceived } = this.state;
    const {webSocketConnected: nextSocket, messages: nextMessages, chatUserNames: nextUserNames, newMessagesReceived: nextMessagesReceived } = nextState;

    return (
      username !== nextUserName ||
      chatInputEnabled !== nextChatEnabled ||
      webSocketConnected !== nextSocket ||
      messages.length !== nextMessages.length ||
      chatUserNames.length !== nextUserNames.length || newMessagesReceived !== nextMessagesReceived
    );
  }


  componentDidUpdate(prevProps, prevState) {
    const { username: prevName } = prevProps;
    const { username } = this.props;

    const { messages: prevMessages } = prevState;
    const { messages } = this.state;

    // if username updated, send a message
    if (prevName !== username) {
      this.sendUsernameChange(prevName, username);
    }

    // scroll to bottom of messages list when new ones come in
    if (messages.length > prevMessages.length) {
      this.setState({
        newMessagesReceived: true,
      });
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
      this.websocket.addListener(CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED, this.receivedWebsocketMessage);
      this.websocket.addListener(CALLBACKS.WEBSOCKET_CONNECTED, this.websocketConnected);
      this.websocket.addListener(CALLBACKS.WEBSOCKET_DISCONNECTED, this.websocketDisconnected);
    }
  }

  // fetch chat history
  getChatHistory() {
    fetch(URL_CHAT_HISTORY)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Network response was not ok ${response.ok}`);
      }
      return response.json();
    })
    .then(data => {
      // extra user names
      const chatUserNames = extraUserNamesFromMessageHistory(data);
      this.setState({
        messages: data,
        chatUserNames,
      });
    })
    .catch(error => {
      this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
    });
  }

  sendUsernameChange(oldName, newName) {
		const nameChange = {
			type: SOCKET_MESSAGE_TYPES.NAME_CHANGE,
			oldName,
			newName,
		};
		this.websocket.send(nameChange);
  }

  receivedWebsocketMessage(message) {
    this.addMessage(message);
  }

  handleNetworkingError(error) {
    // todo: something more useful
    console.log(error);
  }

  addMessage(message) {
    const { messages: curMessages } = this.state;
    const { messagesOnly } = this.props;

    // if incoming message has same id as existing message, don't add it
    const existing = curMessages.filter(function (item) {
      return item.id === message.id;
    })

    if (existing.length === 0 || !existing) {
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
    if (!messagesOnly && message.type === 'CHAT' && this.windowBlurred) {
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
			author: username,
			type: SOCKET_MESSAGE_TYPES.CHAT,
    };
		this.websocket.send(message);
  }

  updateAuthorList(message) {
    const { type } = message;
    const nameList = this.state.chatUserNames;

    if (
      type === SOCKET_MESSAGE_TYPES.CHAT &&
      !nameList.includes(message.author)
    ) {
      return nameList.push(message.author);
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
    const { scrollTop, scrollHeight, clientHeight } = this.scrollableMessagesContainer.current;
    const fullyScrolled = scrollHeight - clientHeight;
    const shouldScroll = scrollHeight >= clientHeight && fullyScrolled - scrollTop < MESSAGE_JUMPTOBOTTOM_BUFFER;
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
      if (this.numMessagesSinceBlur && !this.props.messagesOnly && this.windowBlurred) {
        this.updateDocumentTitle();
      }
    }
  };

  updateDocumentTitle() {
    const num = this.numMessagesSinceBlur > 10 ? '10+' : this.numMessagesSinceBlur;
    window.document.title = `${num} 💬 :: ${this.props.instanceTitle}`;
  }

  render(props, state) {
    const { username, messagesOnly, chatInputEnabled } = props;
    const { messages, chatUserNames, webSocketConnected } = state;

    const messageList = messages.map(
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

