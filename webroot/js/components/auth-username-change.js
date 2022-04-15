import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

import { setLocalStorage } from '../utils/helpers.js';
import { KEY_USERNAME, KEY_CUSTOM_USERNAME_SET } from '../utils/constants.js';

export default class AuthUsernameChange extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      displayForm: true,
      isFocused: true,
    };

    this.textInput = createRef();

    this.handleKeydown = this.handleKeydown.bind(this);
    this.handleUpdateUsername = this.handleUpdateUsername.bind(this);
  }

  handleKeydown(event) {
    if (event.keyCode === 13) {
      // enter
      this.handleUpdateUsername();
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
    }
  }

  componentDidUpdate({}, { displayForm }) {
    if (this.state.displayForm && !displayForm) {
      document.getElementById('username-change-input').select();
    }
  }

  render(props, state) {
    const { username } = props;
    const { displayForm } = state;

    const styles = {
      form: {
        display: displayForm ? 'block' : 'none',
      },
    };

    const moderatorFlag = html`
      <img src="/img/moderator-nobackground.svg" class="moderator-flag" />
    `;
    const userIcon = html`
      <img src="/img/user-icon.svg" class="user-icon-flag" />
    `;

    return html`
      <div id="user-info" class="whitespace-nowrap">
        <div id="user-info-change" style=${styles.form}>
          This is the name you will be known by in the chat.
          <label
            class="block text-gray-700 text-sm font-semibold mt-6"
            for="username-change-input"
          >
            Chat name
          </label>
          <input
            type="text"
            id="username-change-input"
            class="border bg-white rounded w-full py-2 px-3 mb-2 mt-2 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            maxlength="60"
            placeholder="Update username"
            defaultValue=${username}
            onKeydown=${this.handleKeydown}
            ref=${this.textInput}
          />
          <button
            id="button-update-username"
            onClick=${this.handleUpdateUsername}
            type="button"
            class="bg-indigo-500 hover:bg-indigo-600 text-white font-bold py-2 mt-6 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            Update name
          </button>
        </div>
      </div>
    `;
  }
}
