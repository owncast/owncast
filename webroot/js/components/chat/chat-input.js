import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { EmojiButton } from '/js/web_modules/@joeattardi/emoji-button.js';

import ContentEditable, { replaceCaret } from './content-editable.js';
import {
  generatePlaceholderText,
  getCaretPosition,
  convertToText,
  convertOnPaste,
  createEmojiMarkup,
  emojify,
} from '../../utils/chat.js';
import {
  getLocalStorage,
  setLocalStorage,
  classNames,
} from '../../utils/helpers.js';
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

    this.prepNewLine = false;
    this.modifierKeyPressed = false; // control/meta/shift/alt

    this.state = {
      inputHTML: '',
      inputText: '', // for counting
      inputCharsLeft: CHAT_MAX_MESSAGE_LENGTH,
      hasSentFirstChatMessage: getLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT),
      emojiPicker: null,
      emojiList: null,
    };

    this.handleEmojiButtonClick = this.handleEmojiButtonClick.bind(this);
    this.handleEmojiSelected = this.handleEmojiSelected.bind(this);
    this.getCustomEmojis = this.getCustomEmojis.bind(this);

    this.handleMessageInputKeydown = this.handleMessageInputKeydown.bind(this);
    this.handleMessageInputKeyup = this.handleMessageInputKeyup.bind(this);
    this.handleMessageInputBlur = this.handleMessageInputBlur.bind(this);
    this.handleSubmitChatButton = this.handleSubmitChatButton.bind(this);
    this.handlePaste = this.handlePaste.bind(this);

    this.handleContentEditableChange =
      this.handleContentEditableChange.bind(this);
  }

  componentDidMount() {
    this.getCustomEmojis();
  }

  getCustomEmojis() {
    fetch(URL_CUSTOM_EMOJIS)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then((json) => {
        const emojiList = json;
        const emojiPicker = new EmojiButton({
          zIndex: 100,
          theme: 'owncast', // see chat.css
          custom: json,
          initialCategory: 'custom',
          showPreview: false,
          autoHide: false,
          autoFocusSearch: false,
          showAnimation: false,
          emojiSize: '24px',
          position: 'right-start',
          strategy: 'absolute',
        });
        emojiPicker.on('emoji', (emoji) => {
          this.handleEmojiSelected(emoji);
        });
        emojiPicker.on('hidden', () => {
          this.formMessageInput.current.focus();
          replaceCaret(this.formMessageInput.current);
        });
        this.setState({ emojiList, emojiPicker });
      })
      .catch((error) => {
        // this.handleNetworkingError(`Emoji Fetch: ${error}`);
      });
  }

  handleEmojiButtonClick() {
    const { emojiPicker } = this.state;
    if (emojiPicker) {
      emojiPicker.togglePicker(this.emojiPickerButton.current);
    }
  }

  handleEmojiSelected(emoji) {
    const { inputHTML } = this.state;
    let content = '';
    if (emoji.url) {
      content = createEmojiMarkup(emoji, false);
    } else {
      content = emoji.emoji;
    }

    this.setState({
      inputHTML: inputHTML + content,
    });
    // a hacky way add focus back into input field
    setTimeout(() => {
      const input = this.formMessageInput.current;
      input.focus();
      replaceCaret(input);
    }, 100);
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

    if (
      this.completionIndex === undefined ||
      ++this.completionIndex >= possibilities.length
    ) {
      this.completionIndex = 0;
    }

    if (possibilities.length > 0) {
      this.suggestion = possibilities[this.completionIndex];

      this.setState({
        inputHTML:
          inputHTML.substring(0, at + 1) +
          this.suggestion +
          ' ' +
          inputHTML.substring(position),
      });
    }

    return true;
  }

  // replace :emoji: with the emoji <img>
  injectEmoji() {
    const { inputHTML, emojiList } = this.state;
    const textValue = convertToText(inputHTML);
    const processedHTML = emojify(inputHTML, emojiList);

    if (textValue != convertToText(processedHTML)) {
      this.setState({
        inputHTML: processedHTML,
      });
      return true;
    }
    return false;
  }

  handleMessageInputKeydown(event) {
    const formField = this.formMessageInput.current;
    let textValue = formField.textContent; // get this only to count chars
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
        textValue = formField.textContent;
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
    let textValue = formField.textContent; // get this only to count chars

    const { key } = event;

    if (key === 'Control' || key === 'Shift') {
      this.prepNewLine = false;
    }
    if (CHAT_KEY_MODIFIERS.includes(key)) {
      this.modifierKeyPressed = false;
    }
    if (key === ':' || key === ';') {
      if (this.injectEmoji()) {
        // value could have been changed, update char count
        textValue = formField.textContent;
      }
    }
    this.setState({
      inputText: textValue,
      inputCharsLeft: CHAT_MAX_MESSAGE_LENGTH - textValue.length,
    });
  }

  handleMessageInputBlur() {
    this.prepNewLine = false;
    this.modifierKeyPressed = false;
  }

  handlePaste(event) {
    // don't allow paste if too much text already
    if (CHAT_MAX_MESSAGE_LENGTH - this.state.inputText.length < 0) {
      event.preventDefault();
      return;
    }
    convertOnPaste(event, this.state.emojiList);
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
    const { hasSentFirstChatMessage, inputCharsLeft, inputHTML, emojiPicker } =
      state;
    const { inputEnabled } = props;
    const emojiButtonStyle = {
      display: emojiPicker && inputCharsLeft > 0 ? 'block' : 'none',
    };
    const extraClasses = classNames({
      'display-count': inputCharsLeft <= CHAT_CHAR_COUNT_BUFFER,
    });
    const placeholderText = generatePlaceholderText(
      inputEnabled,
      hasSentFirstChatMessage
    );
    return html`
      <div
        id="message-input-container"
        class="relative shadow-md bg-gray-900 border-t border-gray-700 border-solid p-4 z-20 ${extraClasses}"
      >
        <div
          id="message-input-wrap"
          class="flex flex-row justify-end appearance-none w-full bg-gray-200 border border-black-500 rounded py-2 px-2 pr-20 my-2 overflow-auto"
        >
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
        <div
          id="message-form-actions"
          class="absolute flex flex-col justify-end items-end mr-4"
        >
          <span class="flex flex-row justify-center">
            <button
              ref=${this.emojiPickerButton}
              id="emoji-button"
              class="text-3xl leading-3 cursor-pointer text-purple-600"
              type="button"
              style=${emojiButtonStyle}
              onclick=${this.handleEmojiButtonClick}
              aria-label="Select an emoji"
              disabled=${!inputEnabled}
            >
              <img src="../../../img/smiley.png" />
            </button>

            <button
              id="send-message-button"
              class="text-sm text-white rounded bg-gray-600 hidden p-1 ml-1 -mr-2"
              type="button"
              onclick=${this.handleSubmitChatButton}
              disabled=${inputHTML === '' || inputCharsLeft < 0}
              aria-label="Send message"
            >
              Send
            </button>
          </span>

          <span id="message-form-warning" class="text-red-600 text-xs"
            >${inputCharsLeft}/${CHAT_MAX_MESSAGE_LENGTH}</span
          >
        </div>
      </div>
    `;
  }
}
