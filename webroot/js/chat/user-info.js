import { h, Component, createRef } from 'https://unpkg.com/preact?module';
import htm from 'https://unpkg.com/htm?module';
// Initialize htm with Preact
const html = htm.bind(h);

import { generateAvatar, setLocalStorage } from '../utils.js';
import { KEY_USERNAME, KEY_AVATAR } from '../utils/chat.js';


export default class UserInfo extends Component {
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
    const { username, userAvatarImage, handleChatToggle } = props;
    const { displayForm } = state;

    const narrowSpace = document.body.clientWidth < 640;
    const styles = {
      info: {
        display: displayForm || narrowSpace ? 'none' : 'flex',
      },
      form: {
        display: displayForm ? 'flex' : 'none',
      },
    };
    if (narrowSpace) {
      styles.form.display = 'inline-block';
    }
    return (
      html`
        <div id="user-options-container" class="flex">
          <div id="user-info">
            <div id="user-info-display" style=${styles.info} title="Click to update user name" class="flex" onClick=${this.handleDisplayForm}>
              <img
                src=${userAvatarImage}
                alt=""
                class="rounded-full bg-black bg-opacity-50 border border-solid border-gray-700"
              />
              <span class="text-indigo-600">${username}</span>
            </div>

            <div id="user-info-change" style=${styles.form}>
              <input type="text"
                class="appearance-none block w-full bg-gray-200 text-gray-700 border border-black-500 rounded py-1 px-1 leading-tight focus:bg-white"
                maxlength="100"
                placeholder="Update username"
                value=${username}
                onKeydown=${this.handleKeydown}
                ref=${this.textInput}
              >
              <button onClick=${this.handleUpdateUsername} class="bg-blue-500 hover:bg-blue-700 text-white py-1 px-1 rounded user-btn">Update</button>
              <button onClick=${this.handleHideForm} class="bg-gray-900 hover:bg-gray-800 py-1 px-2 rounded user-btn text-white text-opacity-50" title="cancel">X</button>
            </div>
          </div>
          <button type="button" onClick=${handleChatToggle} class="flex bg-gray-800 hover:bg-gray-700">💬</button>
        </div>
    `);
  }
}
