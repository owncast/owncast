import { h, Component, createRef } from '../../web_modules/preact.js';
import htm from '../../web_modules/htm.js';
const html = htm.bind(h);

import { setLocalStorage } from '../../utils/helpers.js';
import { KEY_USERNAME } from '../../utils/constants.js';

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
    if (event.keyCode === 13) { // enter
			this.handleUpdateUsername();
		} else if (event.keyCode === 27) { // esc
			this.handleHideForm();
		}
  }

  handleUpdateUsername() {
    const { username: curName, onUsernameChange } = this.props;
    let newName = this.textInput.current.value;
    newName = newName.trim();
    if (newName !== '' && newName !== curName) {
      setLocalStorage(KEY_USERNAME, newName);
      if (onUsernameChange) {
        onUsernameChange(newName);
      }
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

  render(props, state) {
    const { username } = props;
    const { displayForm } = state;

    const styles = {
      info: {
        display: displayForm ? 'none' : 'flex',
      },
      form: {
        display: displayForm ? 'flex' : 'none',
      },
    };

    return (
      html`
        <div id="user-info" class="whitespace-nowrap">
          <div id="user-info-display" style=${styles.info} title="Click to update user name" class="flex flex-row justify-end items-center cursor-pointer py-2 px-4 overflow-hidden w-full opacity-1 transition-opacity duration-200 hover:opacity-75" onClick=${this.handleDisplayForm}>
            <span id="username-display" class="text-indigo-600 text-xs font-semibold truncate overflow-hidden whitespace-no-wrap">${username}</span>
          </div>

          <div id="user-info-change" class="flex flex-row flex-no-wrap p-1 items-center justify-end" style=${styles.form}>
            <input type="text"
              id="username-change-input"
              class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-1 px-1 leading-tight text-xs focus:bg-white"
              maxlength="60"
              placeholder="Update username"
              defaultValue=${username}
              onKeydown=${this.handleKeydown}
              onFocus=${this.handleFocus}
              onBlur=${this.handleBlur}
              ref=${this.textInput}
            />
            <button id="button-update-username" onClick=${this.handleUpdateUsername}  type="button" class="bg-blue-500 hover:bg-blue-700 text-white text-xs uppercase p-1 mx-1 rounded cursor-pointer user-btn">Update</button>

            <button id="button-cancel-change" onClick=${this.handleHideForm} type="button" class="bg-gray-900 hover:bg-gray-800 py-1 px-2 mx-1 rounded cursor-pointer user-btn text-white text-xs uppercase text-opacity-50" title="cancel">X</button>
          </div>
        </div>
    `);
  }
}
