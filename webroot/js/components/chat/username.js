import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { setLocalStorage } from '../../utils/helpers.js';
import {
  KEY_USERNAME,
  KEY_CUSTOM_USERNAME_SET,
} from '../../utils/constants.js';

import { CheckIcon, CloseIcon, EditIcon } from '../icons/index.js';

export default class UsernameForm extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      displayForm: false,
      isFocused: false,
    };

    this.textInput = createRef();

    this.handleKeydown = this.handleKeydown.bind(this);
    this.handleDisplayForm = this.handleDisplayForm.bind(this);
    this.handleHideForm = this.handleHideForm.bind(this);
    this.handleUpdateUsername = this.handleUpdateUsername.bind(this);
    this.handleFocus = this.handleFocus.bind(this);
    this.handleBlur = this.handleBlur.bind(this);
  }

  handleDisplayForm() {
    const { displayForm: curDisplay } = this.state;
    this.setState({
      displayForm: !curDisplay,
    });
  }

  handleHideForm() {
    this.setState({
      displayForm: false,
    });
  }

  handleKeydown(event) {
    if (event.keyCode === 13) {
      // enter
      this.handleUpdateUsername();
    } else if (event.keyCode === 27) {
      // esc
      this.handleHideForm();
    }
  }

  handleUpdateUsername() {
    const { username: curName, onUsernameChange } = this.props;
    let newName = this.textInput.current.value;
    newName = newName.trim();
    if (newName !== '' && newName !== curName) {
      setLocalStorage(KEY_USERNAME, newName);
      // So we know that the user has set a custom name
      setLocalStorage(KEY_CUSTOM_USERNAME_SET, true);
      if (onUsernameChange) {
        onUsernameChange(newName);
      }
      this.handleHideForm();
    } else {
      this.handleHideForm();
    }
  }

  handleFocus() {
    const { onFocus } = this.props;
    if (onFocus) {
      onFocus();
    }
  }

  handleBlur() {
    const { onBlur } = this.props;
    if (onBlur) {
      onBlur();
    }
  }

  componentDidUpdate({}, { displayForm }) {
    if (this.state.displayForm && !displayForm) {
      document.getElementById('username-change-input').select();
    }
  }

  render(props, state) {
    const { username, isModerator } = props;
    const { displayForm } = state;

    const styles = {
      info: {
        display: displayForm ? 'none' : 'flex',
      },
      form: {
        display: displayForm ? 'flex' : 'none',
      },
    };

    const moderatorFlag = html`
      <img src="/img/moderator-nobackground.svg" class="moderator-flag" />
    `;

    return html`
      <div id="user-info">
        <button
          id="user-info-display"
          style=${styles.info}
          title="Click to update user name"
          class="flex flex-row justify-end items-center align-middle cursor-pointer py-2 px-4 overflow-hidden w-full"
          onClick=${this.handleDisplayForm}
        >
          <span id="username-display">Change Username</span>
          <span><${EditIcon} /></span>
        </button>

        <div
          id="user-info-change"
          class="${displayForm
            ? 'flex'
            : 'hidden'} flex-row flex-no-wrap h-full items-center justify-end"
        >
          <input
            type="text"
            id="username-change-input"
            class="appearance-none block w-full bg-transparent rounded-md text-white p-2"
            maxlength="60"
            placeholder="Update username"
            defaultValue=${username}
            onKeydown=${this.handleKeydown}
            onBlur=${this.handleBlur}
            ref=${this.textInput}
          />
          <div class="username-buttons-wrapper flex ml-2">
            <button
              id="button-update-username"
              onClick=${this.handleUpdateUsername}
              type="button"
              class="bg-purple-500 hover:bg-purple-700 text-white uppercase p-1 ml-1 rounded-md cursor-pointer transition duration-100 user-btn"
            >
              <${CheckIcon} />
            </button>

            <button
              id="button-cancel-change"
              onClick=${this.handleHideForm}
              type="button"
              class="bg-gray-700 hover:bg-gray-500 text-white uppercase p-1 ml-1 rounded-md cursor-pointer transition duration-100 user-btn"
              title="cancel"
            >
              <${CloseIcon} />
            </button>
          </div>
        </div>
      </div>
    `;
  }
}
