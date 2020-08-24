import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
const html = htm.bind(h);

import { generateAvatar, setLocalStorage } from '../../utils/helpers.js';
import { KEY_USERNAME, KEY_AVATAR } from '../../utils/constants.js';

export default class UsernameForm extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      displayForm: false,
    };

    this.textInput = createRef();

    this.handleKeydown = this.handleKeydown.bind(this);
    this.handleDisplayForm = this.handleDisplayForm.bind(this);
    this.handleHideForm = this.handleHideForm.bind(this);
    this.handleUpdateUsername = this.handleUpdateUsername.bind(this);
  }

  handleDisplayForm() {
    this.setState({
      displayForm: true,
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
    const { username: curName, handleUsernameChange } = this.props;
    let newName = this.textInput.current.value;
    newName = newName.trim();
    if (newName !== '' && newName !== curName) {
      const newAvatar = generateAvatar(`${newName}${Date.now()}`);
      setLocalStorage(KEY_USERNAME, newName);
      setLocalStorage(KEY_AVATAR, newAvatar);
      if (handleUsernameChange) {
        handleUsernameChange(newName, newAvatar);
      }
      this.handleHideForm();
    }

  }

  render(props, state) {
    const { username, userAvatarImage } = props;
    const { displayForm } = state;

    const narrowSpace = document.body.clientWidth < 640;
    const formDisplayStyle = narrowSpace ? 'inline-block' : 'flex';
    const styles = {
      info: {
        display: displayForm || narrowSpace ? 'none' : 'flex',
      },
      form: {
        display: displayForm ? formDisplayStyle : 'none',
      },
    };

    return (
      html`
        <div id="user-info">
          <div id="user-info-display" style=${styles.info} title="Click to update user name" class="flex flex-row justify-end items-center cursor-pointer py-2 px-4 overflow-hidden w-full opacity-1 transition-opacity duration-200 hover:opacity-75" onClick=${this.handleDisplayForm}>
            <img
              src=${userAvatarImage}
              alt=""
              id="username-avatar"
              class="rounded-full bg-black bg-opacity-50 border border-solid border-gray-700 mr-2 h-8 w-8"
            />
            <span id="username-display" class="text-indigo-600 text-xs font-semibold truncate overflow-hidden whitespace-no-wrap">${username}</span>
          </div>

          <div id="user-info-change" class="flex flex-no-wrap p-1 items-center justify-end" style=${styles.form}>
            <input type="text"
              id="username-change-input"
              class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-1 px-1 leading-tight text-xs focus:bg-white"
              maxlength="100"
              placeholder="Update username"
              value=${username}
              onKeydown=${this.handleKeydown}
              ref=${this.textInput}
            />
            <button id="button-update-username" onClick=${this.handleUpdateUsername}  type="button" class="bg-blue-500 hover:bg-blue-700 text-white text-xs uppercase p-1 mx-1 rounded cursor-pointer user-btn">Update</button>

            <button id="button-cancel-change" onClick=${this.handleHideForm} type="button" class="bg-gray-900 hover:bg-gray-800 py-1 px-2 mx-1 rounded cursor-pointer user-btn text-white text-xs uppercase text-opacity-50" title="cancel">X</button>
          </div>
        </div>
    `);
  }
}
