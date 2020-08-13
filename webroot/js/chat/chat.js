import { h, Component } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
// Initialize htm with Preact
const html = htm.bind(h);

import { EmojiButton } from 'https://cdn.skypack.dev/@joeattardi/emoji-button';


import SOCKET_MESSAGE_TYPES from '../utils/socket-message-types.js';
import Message from './message.js';
import Websocket, { CALLBACKS } from '../websocket.js';

import { URL_CHAT_HISTORY, URL_CUSTOM_EMOJIS } from '../utils.js';

export default class Chat extends Component {
  constructor(props, context) {
    super(props, context);

    this.messageCharCount = 0;
		this.maxMessageLength = 500;
    this.maxMessageBuffer = 20;

    this.state = {
      inputEnabled: false,
      messages: [],

      chatUserNames: [],
    }

    this.emojiPicker = null;
    this.websocket = null;


    this.handleEmojiButtonClick = this.handleEmojiButtonClick.bind(this);
    this.handleEmojiSelected = this.handleEmojiSelected.bind(this);
    this.getCustomEmojis = this.getCustomEmojis.bind(this);
    this.getChatHistory = this.getChatHistory.bind(this);
    this.receivedWebsocketMessage = this.receivedWebsocketMessage.bind(this);
    this.websocketDisconnected = this.websocketDisconnected.bind(this);
  }

  componentDidMount() {
    /*
    - set up websocket
    - get emojis
    - get chat history
    */
   this.setupWebSocket();
   this.getChatHistory();
   this.getCustomEmojis();

  }

  componentDidUpdate(prevProps) {
    const { username: prevName } = prevProps;
    const { username, userAvatarImage } = this.props;
    // if username updated, send a message
    if (prevName !== username) {
      this.sendUsernameChange(prevName, username, userAvatarImage);
    }

  }

  setupWebSocket() {
    this.websocket = new Websocket();
    this.websocket.addListener(CALLBACKS.RAW_WEBSOCKET_MESSAGE_RECEIVED, this.receivedWebsocketMessage);
    this.websocket.addListener(CALLBACKS.WEBSOCKET_DISCONNECTED, this.websocketDisconnected);
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
      this.setState({
        messages: data,
      });
      // const formattedMessages = data.map(function (message) {
      //   return new Message(message);
      // })
      // this.vueApp.messages = formattedMessages.concat(this.vueApp.messages);
    })
    .catch(error => {
      this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
    });
  }

  getCustomEmojis() {
    fetch(URL_CUSTOM_EMOJIS)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then(json => {
        this.emojiPicker = new EmojiButton({
          zIndex: 100,
          theme: 'dark',
          custom: json,
          initialCategory: 'custom',
          showPreview: false,
          position: {
            top: '50%',
            right: '100'
          },
        });
        this.emojiPicker.on('emoji', emoji => {
          this.handleEmojiSelected(emoji);
        });
      })
      .catch(error => {
        this.handleNetworkingError(`Emoji Fetch: ${error}`);
      });
  }

  sendUsernameChange(oldName, newName, image) {
		const nameChange = {
			type: SOCKET_MESSAGE_TYPES.NAME_CHANGE,
			oldName: oldName,
			newName: newName,
			image: image,
		};
		this.send(nameChange);
  }

  handleEmojiButtonClick() {
    if (this.emojiPicker) {
      this.emojiPicker.togglePicker(this.picker);
    }
  }

  handleEmojiSelected(emoji) {
    if (emoji.url) {
      const url = location.protocol + "//" + location.host + "/" + emoji.url;
      const name = url.split('\\').pop().split('/').pop();
      document.querySelector('#message-body-form').innerHTML += "<img class=\"emoji\" alt=\"" + name + "\" src=\"" + url + "\"/>";
    } else {
      document.querySelector('#message-body-form').innerHTML += emoji.emoji;
    }
  }

  receivedWebsocketMessage(message) {
    this.addMessage(message);
    // if (model.type === SOCKET_MESSAGE_TYPES.CHAT) {
    //   const message = new Message(model);
    //   this.addMessage(message);
    // } else if (model.type === SOCKET_MESSAGE_TYPES.NAME_CHANGE) {
    //   this.addMessage(model);
    // }
  }

  addMessage(message) {
    const { messages: curMessages } = this.state;
    const existing = curMessages.filter(function (item) {
      return item.id === message.id;
    })
    if (existing.length === 0 || !existing) {
      this.setState({
        messages: [...curMessages, message],
      });
    }
  }
  websocketDisconnected() {
    this.websocket = null;
  }

  render(props, state) {
    const { username } = props;
    const { messages } = state;

    return (
      html`
        <section id="chat-container-wrap" class="flex">
          <div id="chat-container" class="bg-gray-800">
            <div id="messages-container">
              ${
                messages.map(message => (html`<${Message} message=${message} />`))
              }
              messages..
            </div>


            <div id="message-input-container" class="shadow-md bg-gray-900 border-t border-gray-700 border-solid">
              <form id="message-form" class="flex">

                <input type="hidden" name="inputAuthor" id="self-message-author" value=${username} />

                <textarea
                  disabled
                  id="message-body-form"
                  placeholder="Message"
                  class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-2 px-2 my-2 focus:bg-white"
                ></textarea>
                <button
                  id="emoji-button"
                  onClick=${this.handleEmojiButtonClick}
                >😏</button>

                <div id="message-form-actions" class="flex">
                  <span id="message-form-warning" class="text-red-600 text-xs"></span>
                  <button
                    id="button-submit-message"
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded"
                  > Chat
                  </button>
                </div>
              </form>
            </div>
          </div>
        </section>
    `);
  }

}






