import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import { EmojiButton } from 'https://cdn.skypack.dev/pin/@joeattardi/emoji-button@v4.1.0-v8psdkkxts3LNdpA0m5Q/min/@joeattardi/emoji-button.js';
import ContentEditable from './content-editable.js';
import { generatePlaceholderText, getCaretPosition, convertToText, convertOnPaste } from '../../utils/chat.js';
import { getLocalStorage, setLocalStorage, classNames } from '../../utils/helpers.js';
import {
  URL_CUSTOM_EMOJIS,
  KEY_CHAT_FIRST_MESSAGE_SENT,
  CHAT_MAX_MESSAGE_LENGTH,
  CHAT_CHAR_COUNT_BUFFER,
  CHAT_OK_KEYCODES,
  CHAT_KEY_MODIFIERS,
} from '../../utils/constants.js';

export default class ChatInput extends Component {
  constructor(props, context) {
    super(props, context);
    this.formMessageInput = createRef();
    this.emojiPickerButton = createRef();

    this.messageCharCount = 0;

    this.emojiPicker = null;

    this.prepNewLine = false;
    this.modifierKeyPressed = false; // control/meta/shift/alt

    this.state = {
      inputHTML: '',
      inputText: '', // for counting
      inputCharsLeft: CHAT_MAX_MESSAGE_LENGTH,
      hasSentFirstChatMessage: getLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT),
    };

    this.handleEmojiButtonClick = this.handleEmojiButtonClick.bind(this);
    this.handleEmojiSelected = this.handleEmojiSelected.bind(this);
    this.getCustomEmojis = this.getCustomEmojis.bind(this);

    this.handleMessageInputKeydown = this.handleMessageInputKeydown.bind(this);
    this.handleMessageInputKeyup = this.handleMessageInputKeyup.bind(this);
    this.handleMessageInputBlur = this.handleMessageInputBlur.bind(this);
    this.handleSubmitChatButton = this.handleSubmitChatButton.bind(this);
    this.handlePaste = this.handlePaste.bind(this);

    this.handleContentEditableChange = this.handleContentEditableChange.bind(this);
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
          theme: 'owncast', // see chat.css
          custom: json,
          initialCategory: 'custom',
          showPreview: false,
          emojiSize: '24px',
          position: 'right-start',
          strategy: 'absolute',
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
      this.emojiPicker.togglePicker(this.emojiPickerButton.current);
    }
  }

  handleEmojiSelected(emoji) {
    const { inputHTML } = this.state;
    let content = '';
    if (emoji.url) {
      const url = location.protocol + "//" + location.host + "/" + emoji.url;
      const name = url.split('\\').pop().split('/').pop();
      content = "<img class=\"emoji\" alt=\"" + name + "\" src=\"" + url + "\"/>";
    } else {
      content = emoji.emoji;
    }

    this.setState({
      inputHTML: inputHTML + content,
    });
  }

  // autocomplete user names
  autoCompleteNames() {
    const { chatUserNames } = this.props;
    const { inputHTML } = this.state;
    const position = getCaretPosition(this.formMessageInput.current);
    const at = inputHTML.lastIndexOf('@', position - 1);
    if (at === -1) {
      return false;
    }

    let partial = inputHTML.substring(at + 1, position).trim();

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

      this.setState({
        inputHTML: inputHTML.substring(0, at + 1) + this.suggestion + ' ' + inputHTML.substring(position),
      })
    }

    return true;
  }

  handleMessageInputKeydown(event) {
    const formField = this.formMessageInput.current;

    let textValue = formField.innerText.trim(); // get this only to count chars
    const newStates = {};
    let numCharsLeft = CHAT_MAX_MESSAGE_LENGTH - textValue.length;
    const key = event && event.key;

    if (key === 'Enter') {
      if (!this.prepNewLine) {
        this.sendMessage();
        event.preventDefault();
        this.prepNewLine = false;
        return;
      }
    }
    // allow key presses such as command/shift/meta, etc even when message length is full later.
    if (CHAT_KEY_MODIFIERS.includes(key)) {
      this.modifierKeyPressed = true;
    }
    if (key === 'Control' || key === 'Shift') {
      this.prepNewLine = true;
    }
    if (key === 'Tab') {
      if (this.autoCompleteNames()) {
        event.preventDefault();

        // value could have been changed, update char count
        textValue = formField.innerText.trim();
        numCharsLeft = CHAT_MAX_MESSAGE_LENGTH - textValue.length;
      }
    }

    if (numCharsLeft <= 0 && !CHAT_OK_KEYCODES.includes(key)) {
      newStates.inputText = textValue;
      this.setState(newStates);
      if (!this.modifierKeyPressed) {
        event.preventDefault(); // prevent typing more
      }
      return;
    }
    newStates.inputText = textValue;
    this.setState(newStates);
  }

  handleMessageInputKeyup(event) {
    const formField = this.formMessageInput.current;
    const textValue = formField.innerText.trim(); // get this only to count chars

    const { key } = event;

    if (key === 'Control' || key === 'Shift') {
     this.prepNewLine = false;
    }
    if (CHAT_KEY_MODIFIERS.includes(key)) {
      this.modifierKeyPressed = false;
    }
    this.setState({
      inputCharsLeft: CHAT_MAX_MESSAGE_LENGTH - textValue.length,
    });
  }

  handleMessageInputBlur(event) {
    this.prepNewLine = false;
    this.modifierKeyPressed = false;
  }

  handlePaste(event) {
    // don't allow paste if too much text already
    if (CHAT_MAX_MESSAGE_LENGTH - this.state.inputText.length < 0) {
      event.preventDefault();
      return;
    }
    convertOnPaste(event);
    this.handleMessageInputKeydown(event);
  }

  handleSubmitChatButton(event) {
    event.preventDefault();
    this.sendMessage();
  }

  sendMessage() {
    const { handleSendMessage } = this.props;
    const { hasSentFirstChatMessage, inputHTML, inputText } = this.state;
    if (CHAT_MAX_MESSAGE_LENGTH - inputText.length < 0) {
      return;
    }
    const message = convertToText(inputHTML);
    const newStates = {
      inputHTML: '',
      inputText: '',
      inputCharsLeft: CHAT_MAX_MESSAGE_LENGTH,
    };

    handleSendMessage(message);

    if (!hasSentFirstChatMessage) {
      newStates.hasSentFirstChatMessage = true;
      setLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT, true);
    }

    // clear things out.
    this.setState(newStates);
  }

  handleContentEditableChange(event) {
    this.setState({ inputHTML: event.target.value });
  }

  render(props, state) {
    const { hasSentFirstChatMessage, inputCharsLeft, inputHTML } = state;
    const { inputEnabled } = props;
    const emojiButtonStyle = {
      display: this.emojiPicker && inputCharsLeft > 0 ? 'block' : 'none',
    };
    const extraClasses = classNames({
      'display-count': inputCharsLeft <= CHAT_CHAR_COUNT_BUFFER,
    });
    const placeholderText = generatePlaceholderText(inputEnabled, hasSentFirstChatMessage);
    return (
      html`
        <div id="message-input-container" class="relative shadow-md bg-gray-900 border-t border-gray-700 border-solid p-4 z-20 ${extraClasses}">

          <div
            id="message-input-wrap"
            class="flex flex-row justify-end appearance-none w-full bg-gray-200 border border-black-500 rounded py-2 px-2 pr-12 my-2 overflow-auto">
            <${ContentEditable}
              id="message-input"
              class="appearance-none block w-full bg-transparent text-sm text-gray-700 h-full focus:outline-none"

              placeholderText=${placeholderText}
              innerRef=${this.formMessageInput}
              html=${inputHTML}
              disabled=${!inputEnabled}
              onChange=${this.handleContentEditableChange}
              onKeyDown=${this.handleMessageInputKeydown}
              onKeyUp=${this.handleMessageInputKeyup}
              onBlur=${this.handleMessageInputBlur}

              onPaste=${this.handlePaste}
            />
          </div>
          <div id="message-form-actions" class="absolute flex flex-col w-10 justify-end items-center">
              <button
                ref=${this.emojiPickerButton}
                id="emoji-button"
                class="text-3xl leading-3 cursor-pointer text-purple-600"
                type="button"
                style=${emojiButtonStyle}
                onclick=${this.handleEmojiButtonClick}
                disabled=${!inputEnabled}
              ><img src="../../../img/smiley.png" /></button>

              <span id="message-form-warning" class="text-red-600 text-xs">${inputCharsLeft}/${CHAT_MAX_MESSAGE_LENGTH}</span>
            </div>
      </div>
    `);
  }

  }
