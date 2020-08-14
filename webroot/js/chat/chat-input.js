import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import { EmojiButton } from 'https://cdn.skypack.dev/@joeattardi/emoji-button';

import { URL_CUSTOM_EMOJIS } from '../utils.js';
import { generatePlaceholderText } from '../utils/chat.js';

export default class ChatInput extends Component {
  constructor(props, context) {
    super(props, context);
    this.formMessageInput = createRef();

    this.messageCharCount = 0;
		this.maxMessageLength = 500;
    this.maxMessageBuffer = 20;

    this.emojiPicker = null;

    this.prepNewLine = false;

    this.handleEmojiButtonClick = this.handleEmojiButtonClick.bind(this);
    this.handleEmojiSelected = this.handleEmojiSelected.bind(this);
    this.getCustomEmojis = this.getCustomEmojis.bind(this);

    this.handleMessageInputKeydown = this.handleMessageInputKeydown.bind(this);
    this.handleMessageInputKeyup = this.handleMessageInputKeyup.bind(this);
    this.handleMessageInputBlur = this.handleMessageInputBlur.bind(this);
    this.handleMessageInput = this.handleMessageInput.bind(this);
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
    if (emoji.url) {
      const url = location.protocol + "//" + location.host + "/" + emoji.url;
      const name = url.split('\\').pop().split('/').pop();
      document.querySelector('#message-body-form').innerHTML += "<img class=\"emoji\" alt=\"" + name + "\" src=\"" + url + "\"/>";
    } else {
      document.querySelector('#message-body-form').innerHTML += emoji.emoji;
    }
  }

  // autocomplete user names
	autoCompleteNames() {
    const { chatUserNames } = this.props;
		const rawValue = this.formMessageInput.innerHTML;
		const position = getCaretPosition(this.formMessageInput);
		const at = rawValue.lastIndexOf('@', position - 1);

		if (at === -1) {
			return false;
		}

		var partial = rawValue.substring(at + 1, position).trim();

		if (partial === this.suggestion) {
			partial = this.partial;
		} else {
			this.partial = partial;
		}

		const possibilities = chatUsernames.filter(function (username) {
			return username.toLowerCase().startsWith(partial.toLowerCase());
		});

		if (this.completionIndex === undefined || ++this.completionIndex >= possibilities.length) {
			this.completionIndex = 0;
		}

		if (possibilities.length > 0) {
			this.suggestion = possibilities[this.completionIndex];

			// TODO: Fix the space not working.  I'm guessing because the DOM ignores spaces and it requires a nbsp or something?
			this.formMessageInput.innerHTML = rawValue.substring(0, at + 1) + this.suggestion + ' ' + rawValue.substring(position);
			setCaretPosition(this.formMessageInput, at + this.suggestion.length + 2);
		}

		return true;
  }

  handleMessageInputKeydown(event) {
    console.log("========this.formMessageInput", this.formMessageInput)
		const okCodes = [37,38,39,40,16,91,18,46,8];
		const value = this.formMessageInput.innerHTML.trim();
		const numCharsLeft = this.maxMessageLength - value.length;
		if (event.keyCode === 13) { // enter
			if (!this.prepNewLine) {
				this.submitChat(value);
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
				value = this.formMessageInput.innerHTML.trim();
				numCharsLeft = this.maxMessageLength - value.length;
			}
		}

		// if (numCharsLeft <= this.maxMessageBuffer) {
		// 	this.tagMessageFormWarning.innerText = `${numCharsLeft} chars left`;
		// 	if (numCharsLeft <= 0 && !okCodes.includes(event.keyCode)) {
		// 		event.preventDefault();
		// 		return;
		// 	}
		// } else {
		// 	this.tagMessageFormWarning.innerText = '';
		// }
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
    // event.target.value
  }

  // setChatPlaceholderText() {
	// 	// NOTE: This is a fake placeholder that is being styled via CSS.
	// 	// You can't just set the .placeholder property because it's not a form element.
	// 	const hasSentFirstChatMessage = getLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT);
	// 	const placeholderText = hasSentFirstChatMessage ? CHAT_PLACEHOLDER_TEXT : CHAT_INITIAL_PLACEHOLDER_TEXT;
	// 	this.formMessageInput.setAttribute("placeholder", placeholderText);
  // }

  render(props, state) {
    const { contenteditable, hasSentFirstChatMessage } = props;
    const emojiButtonStyle = {
      display: this.emojiPicker ? 'block' : 'none',
    };

    const placeholderText = generatePlaceholderText(contenteditable, hasSentFirstChatMessage);

    return (
      html`
        <div>
          <div
            class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-2 px-2 my-2 focus:bg-white"
            onkeydown=${this.handleMessageInputKeydown}
            onkeyup=${this.handleMessageInputKeyup}
            onblur=${this.handleMessageInputBlur}
            oninput=${this.handleMessageInput}
            contenteditable=${contenteditable}
            placeholder=${placeholderText}
            ref=${this.formMessageInput}

          ></div>
          <button
            type="button"
            style=${emojiButtonStyle}
            onclick=${this.handleEmojiButtonClick}
          >😏</button>
      </div>
    `);
  }

}
