import { h, Component, render } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
// Initialize htm with Preact
const html = htm.bind(h);

import SOCKET_MESSAGE_TYPES from '../utils/socketMessageTypes.js';


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

  }

  componentDidMount() {

  }

  componentDidUpdate(prevProps) {
    const { username: prevName } = prevProps;
    const { username, userAvatarImage } = this.props;
    // if username updated, send a message
    if (prevName !== username) {
      this.sendUsernameChange(prevName, username, userAvatarImage);
    }

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

  render() {
    const { username, userAvatarImage } = this.state;
    return (
      html`
        <section id="chat-container-wrap" class="flex">
          <div id="chat-container" class="bg-gray-800">
            <div id="messages-container">
              messages...
              <!-- <div v-for="message in messages" v-cloak>
                  <div class="message flex" v-if="message.type === 'CHAT'">
                    <div class="message-avatar rounded-full flex items-center justify-center" v-bind:style="{ backgroundColor: message.userColor() }">
                      <img
                        v-bind:src="message.image"
                      />
                    </div>
                    <div class="message-content">
                      <p class="message-author text-white font-bold">{{ message.author }}</p>
                      <p class="message-text text-gray-400 font-thin " v-html="message.formatText()"></p>
                    </div>
                </div>
                <div class="message flex" v-else-if="message.type === 'NAME_CHANGE'">
                    <img
                    class="mr-2"
                    width="30px"
                      v-bind:src="message.image"
                    />
                    <div class="text-white text-center">
                      <span class="font-bold">{{ message.oldName }}</span> is now known as <span class="font-bold">{{ message.newName }}</span>.
                    </div>
                </div>
              </div> -->
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






