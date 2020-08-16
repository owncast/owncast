import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import { EmojiButton } from 'https://cdn.skypack.dev/@joeattardi/emoji-button';

import { URL_CUSTOM_EMOJIS, getLocalStorage } from '../utils.js';
import {
  KEY_CHAT_FIRST_MESSAGE_SENT,
  generatePlaceholderText,
  getCaretPosition,
  setCaretPosition,
} from '../utils/chat.js';

export default class ChatInput extends Component {
  constructor(props, context) {
    super(props, context);
    this.formMessageInput = createRef();

    this.messageCharCount = 0;
		this.maxMessageLength = 500;
    this.maxMessageBuffer = 20;

    this.emojiPicker = null;

    this.prepNewLine = false;

    this.state = {
      inputValue: '',
      inputWarning: '',
      hasSentFirstChatMessage: getLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT),
    };

    this.handleEmojiButtonClick = this.handleEmojiButtonClick.bind(this);
    this.handleEmojiSelected = this.handleEmojiSelected.bind(this);
    this.getCustomEmojis = this.getCustomEmojis.bind(this);

    this.handleMessageInputKeydown = this.handleMessageInputKeydown.bind(this);
    this.handleMessageInputKeyup = this.handleMessageInputKeyup.bind(this);
    this.handleMessageInputBlur = this.handleMessageInputBlur.bind(this);
    this.handleMessageInput = this.handleMessageInput.bind(this);
    this.handleSubmitChatButton = this.handleSubmitChatButton.bind(this);
    this.handlePaste = this.handlePaste.bind(this);
  }

  componentDidMount() {
    this.getCustomEmojis();
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
        // this.handleNetworkingError(`Emoji Fetch: ${error}`);
      });
  }

  handleEmojiButtonClick() {
    if (this.emojiPicker) {
      this.emojiPicker.togglePicker(this.emojiPicker);
    }
  }

  handleEmojiSelected(emoji) {
    let content = '';
    if (emoji.url) {
      const url = location.protocol + "//" + location.host + "/" + emoji.url;
      const name = url.split('\\').pop().split('/').pop();
      content = "<img class=\"emoji\" alt=\"" + name + "\" src=\"" + url + "\"/>";
    } else {
      content = emoji.emoji;
    }

    this.formMessageInput.current.innerHTML += content;
  }

  // autocomplete user names
	autoCompleteNames() {
    const { chatUserNames } = this.props;
		const rawValue = this.formMessageInput.current.innerHTML;
		const position = getCaretPosition(this.formMessageInput);
		const at = rawValue.lastIndexOf('@', position - 1);

		if (at === -1) {
			return false;
		}

		const partial = rawValue.substring(at + 1, position).trim();

		if (partial === this.suggestion) {
			partial = this.partial;
		} else {
			this.partial = partial;
		}

		const possibilities = chatUserNames.filter(function (username) {
			return username.toLowerCase().startsWith(partial.toLowerCase());
		});

		if (this.completionIndex === undefined || ++this.completionIndex >= possibilities.length) {
			this.completionIndex = 0;
		}

		if (possibilities.length > 0) {
			this.suggestion = possibilities[this.completionIndex];

			// TODO: Fix the space not working.  I'm guessing because the DOM ignores spaces and it requires a nbsp or something?
			this.formMessageInput.current.innerHTML = rawValue.substring(0, at + 1) + this.suggestion + ' ' + rawValue.substring(position);
			setCaretPosition(this.formMessageInput.current, at + this.suggestion.length + 2);
		}

		return true;
  }

  handleMessageInputKeydown(event) {
    const okCodes = [37,38,39,40,16,91,18,46,8];
    const formField = this.formMessageInput.current;
    const htmlValue = formField.innerHTML.trim();
    let textValue = formField.innerText.trim();
		let numCharsLeft = this.maxMessageLength - textValue.length;
		if (event.keyCode === 13) { // enter
			if (!this.prepNewLine) {
				this.sendMessage();
				event.preventDefault();
				this.prepNewLine = false;
				return;
			}
		}
		if (event.keyCode === 16 || event.keyCode === 17) { // ctrl, shift
			this.prepNewLine = true;
		}
		if (event.keyCode === 9) { // tab
			if (this.autoCompleteNames()) {
				event.preventDefault();

				// value could have been changed, update variables
				textValue = formField.innerText.trim();
				numCharsLeft = this.maxMessageLength - textValue.length;
			}
		}

    // text count
		if (numCharsLeft <= this.maxMessageBuffer) {
			this.setState({
        inputWarning: `${numCharsLeft} chars left`,
      });
			if (numCharsLeft <= 0 && !okCodes.includes(event.keyCode)) {
				event.preventDefault(); // prevent typing more
				return;
			}
		} else {
      this.setState({
        inputWarning: '',
      });
		}
	}

	handleMessageInputKeyup(event) {
		if (event.keyCode === 16 || event.keyCode === 17) { // ctrl, shift
			this.prepNewLine = false;
		}
	}

	handleMessageInputBlur(event) {
		this.prepNewLine = false;
  }

  handleMessageInput(event) {
    console.log("========on input",event.target, this.formMessageInput.current.innerHTML, this.formMessageInput.current.innerText);
  }

  handlePaste(event) {
    event.preventDefault();
    document.execCommand('inserttext', false, event.clipboardData.getData('text/plain'));
  }

  handleSubmitChatButton(event) {
    event.preventDefault();
    this.sendMessage();
  }

  sendMessage() {
    const { handleSendMessage } = this.props;
    const { hasSentFirstChatMessage } = this.state;
    const message = this.formMessageInput.current.innerHTML.trim();
    const newStates = {
      inputWarning: '',
    };

    handleSendMessage(message);

    if (!hasSentFirstChatMessage) {
      newStates.hasSentFirstChatMessage = true;
      setLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT, true);
    }

    // clear things out.
    this.setState(newStates);
    this.formMessageInput.current.innerHTML = '';
  }

  render(props, state) {
    const { hasSentFirstChatMessage, inputWarning } = state;
    const { inputEnabled } = props;
    const emojiButtonStyle = {
      display: this.emojiPicker ? 'block' : 'none',
    };

    const placeholderText = generatePlaceholderText(inputEnabled, hasSentFirstChatMessage);

    return (
      html`
        <div id="message-input-container" class="shadow-md bg-gray-900 border-t border-gray-700 border-solid">
          <div
            class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-2 px-2 my-2 focus:bg-white"
            onkeydown=${this.handleMessageInputKeydown}
            onkeyup=${this.handleMessageInputKeyup}
            onblur=${this.handleMessageInputBlur}
            oninput=${this.handleMessageInput}
            onpaste=${this.handlePaste}
            contenteditable=${inputEnabled}
            placeholder=${placeholderText}
            ref=${this.formMessageInput}
          ></div>
          <button
            type="button"
            style=${emojiButtonStyle}
            onclick=${this.handleEmojiButtonClick}
          >😏</button>
          <div id="message-form-actions" class="flex">
            <span id="message-form-warning" class="text-red-600 text-xs">${inputWarning}</span>
            <button
              onclick=${this.handleSubmitChatButton}
              type="button"
              id="button-submit-message"
              class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded"
            > Chat
            </button>
          </div>
      </div>
    `);
  }

}
