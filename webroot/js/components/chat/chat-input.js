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
  trimNbsp,
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
      inputCharsLeft: props.inputMaxBytes,
      hasSentFirstChatMessage: getLocalStorage(KEY_CHAT_FIRST_MESSAGE_SENT),
      emojiPicker: null,
      emojiList: null,
      emojiNames: null,
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
        const emojiNames = emojiList.map((emoji) => emoji.name);
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
        this.setState({ emojiNames, emojiList, emojiPicker });
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
    const { inputHTML, inputCharsLeft } = this.state;
    // if we're already at char limit, don't do anything
    if (inputCharsLeft < 0) {
      return;
    }
    let content = '';
    if (emoji.url) {
      content = createEmojiMarkup(emoji, false);
    } else {
      content = emoji.emoji;
    }

    const position = getCaretPosition(this.formMessageInput.current);
    const newHTML =
      inputHTML.substring(0, position) +
      content +
      inputHTML.substring(position);

    const charsLeft = this.calculateCurrentBytesLeft(newHTML);
    this.setState({
      inputHTML: newHTML,
      inputCharsLeft: charsLeft,
    });
    // a hacky way add focus back into input field
    setTimeout(() => {
      const input = this.formMessageInput.current;
      input.focus();
      replaceCaret(input);
    }, 100);
  }

  // autocomplete text from the given "list". "token" marks the start of word lookup.
  autoComplete(token, list) {
    const { inputHTML } = this.state;
    const position = getCaretPosition(this.formMessageInput.current);
    const at = inputHTML.lastIndexOf(token, position - 1);
    if (at === -1) {
      return false;
    }

    let partial = inputHTML.substring(at + 1, position).trim();

    if (this.partial === undefined) {
      this.partial = [];
    }

    if (partial === this.suggestion) {
      partial = this.partial[token];
    } else {
      this.partial[token] = partial;
    }

    const possibilities = list.filter(function (item) {
      return item.toLowerCase().startsWith(partial.toLowerCase());
    });

    if (this.completionIndex === undefined) {
      this.completionIndex = [];
    }

    if (
      this.completionIndex[token] === undefined ||
      ++this.completionIndex[token] >= possibilities.length
    ) {
      this.completionIndex[token] = 0;
    }

    if (possibilities.length > 0) {
      this.suggestion = possibilities[this.completionIndex[token]];

      const newHTML =
        inputHTML.substring(0, at + 1) +
        this.suggestion +
        ' ' +
        inputHTML.substring(position);

      this.setState({
        inputHTML: newHTML,
        inputCharsLeft: this.calculateCurrentBytesLeft(newHTML),
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
      const { chatUserNames } = this.props;
      const { emojiNames } = this.state;
      if (this.autoComplete('@', chatUserNames)) {
        event.preventDefault();
      }
      if (this.autoComplete(':', emojiNames)) {
        event.preventDefault();
      }
    }

    // if new input pushes the potential chars over, don't do anything
    const formField = this.formMessageInput.current;
    const tempCharsLeft = this.calculateCurrentBytesLeft(formField.innerHTML);
    if (tempCharsLeft <= 0 && !CHAT_OK_KEYCODES.includes(key)) {
      if (!this.modifierKeyPressed) {
        event.preventDefault(); // prevent typing more
      }
      return;
    }
  }

  handleMessageInputKeyup(event) {
    const { key } = event;
    if (key === 'Control' || key === 'Shift') {
      this.prepNewLine = false;
    }
    if (CHAT_KEY_MODIFIERS.includes(key)) {
      this.modifierKeyPressed = false;
    }

    if (key === ':' || key === ';') {
      this.injectEmoji();
    }
  }

  handleMessageInputBlur() {
    this.prepNewLine = false;
    this.modifierKeyPressed = false;
  }

  handlePaste(event) {
    // don't allow paste if too much text already
    if (this.state.inputCharsLeft < 0) {
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
    const { handleSendMessage, inputMaxBytes } = this.props;
    const { hasSentFirstChatMessage, inputHTML, inputCharsLeft } = this.state;
    if (inputCharsLeft < 0) {
      return;
    }
    const message = convertToText(inputHTML);
    const newStates = {
      inputHTML: '',
      inputCharsLeft: inputMaxBytes,
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
    const value = event.target.value;
    this.setState({
      inputHTML: value,
      inputCharsLeft: this.calculateCurrentBytesLeft(value),
    });
  }

  calculateCurrentBytesLeft(inputContent) {
    const { inputMaxBytes } = this.props;
    const curBytes = new Blob([trimNbsp(inputContent)]).size;
    return inputMaxBytes - curBytes;
  }

  render(props, state) {
    const { hasSentFirstChatMessage, inputCharsLeft, inputHTML, emojiPicker } =
      state;
    const { inputEnabled, inputMaxBytes } = props;
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
            aria-role="textbox"
            class="appearance-none block w-full bg-transparent text-sm text-gray-700 h-full focus:outline-none"
            aria-placeholder=${placeholderText}
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
            >${inputCharsLeft} bytes</span
          >
        </div>
      </div>
    `;
  }
}
