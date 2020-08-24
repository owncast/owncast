import { h, Component } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import Message from './message.js';
import ChatInput from './chat-input.js';
import { CALLBACKS, SOCKET_MESSAGE_TYPES } from '../../utils/websocket.js';
import { setVHvar, hasTouchScreen } from '../../utils/helpers.js';
import { extraUserNamesFromMessageHistory } from '../../utils/chat.js';
import { URL_CHAT_HISTORY } from '../../utils/constants.js';

export default class Chat extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      inputEnabled: true,
      messages: [],
      chatUserNames: [],
    };

    this.websocket = null;

    this.getChatHistory = this.getChatHistory.bind(this);
    this.receivedWebsocketMessage = this.receivedWebsocketMessage.bind(this);
    this.websocketDisconnected = this.websocketDisconnected.bind(this);

    // this.handleSubmitChatButton = this.handleSubmitChatButton.bind(this);
    this.submitChat = this.submitChat.bind(this);
  }

  componentDidMount() {
   this.setupWebSocketCallbacks();
   this.getChatHistory();

   if (hasTouchScreen()) {
    setVHvar();
    window.addEventListener("orientationchange", setVHvar);
    // this.tagAppContainer.classList.add('touch-screen');
    }
  }

  componentDidUpdate(prevProps) {
    const { username: prevName } = prevProps;
    const { username, userAvatarImage } = this.props;
    // if username updated, send a message
    if (prevName !== username) {
      this.sendUsernameChange(prevName, username, userAvatarImage);
    }
  }

  setupWebSocketCallbacks() {
    this.websocket = this.props.websocket;
    if (this.websocket) {
      this.websocket.addListener(CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED, this.receivedWebsocketMessage);
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
      // this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
    });
  }

  sendUsernameChange(oldName, newName, image) {
		const nameChange = {
			type: SOCKET_MESSAGE_TYPES.NAME_CHANGE,
			oldName: oldName,
			newName: newName,
			image: image,
		};
		this.websocket.send(nameChange);
  }

  receivedWebsocketMessage(message) {
    this.addMessage(message);
  }

  // if incoming message has same id as existing message, don't add it
  addMessage(message) {
    const { messages: curMessages } = this.state;
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

    // todo - jump to bottom
    // jumpToBottom(this.scrollableMessagesContainer);
  }
  websocketDisconnected() {
    // this.websocket = null;
    this.disableChat();
  }

  submitChat(content) {
		if (!content) {
			return;
    }
    const { username, userAvatarImage } = this.props;
    const message = {
      body: content,
			author: username,
			image: userAvatarImage,
			type: SOCKET_MESSAGE_TYPES.CHAT,
    };
		this.websocket.send(message);
  }

  disableChat() {
    this.setState({
      inputEnabled: false,
    });
	}

	enableChat() {
    this.setState({
      inputEnabled: true,
    });
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


  render(props, state) {
    const { username, messagesOnly, chatEnabled } = props;
    const { messages, inputEnabled, chatUserNames } = state;

    const messageList = messages.map((message) => (html`<${Message} message=${message} username=${username} key=${message.id} />`));

    if (messagesOnly) {
      return (
        html`
          <div id="messages-container" class="py-1 overflow-auto">
            ${messageList}
          </div>
      `);
    }

    return (
      html`
        <section id="chat-container-wrap" class="flex flex-col">
          <div id="chat-container" class="bg-gray-800 flex flex-col justify-end overflow-auto">
            <div id="messages-container" class="py-1 overflow-auto">
              ${messageList}
            </div>
            <${ChatInput}
              chatUserNames=${chatUserNames}
              inputEnabled=${chatEnabled && inputEnabled}
              handleSendMessage=${this.submitChat}
            />
          </div>
        </section>
    `);
  }

}

