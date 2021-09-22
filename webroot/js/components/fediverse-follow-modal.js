import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import MicroModal from '/js/web_modules/micromodal/dist/micromodal.min.js';
import { ExternalActionButton } from './external-action-modal.js';


const html = htm.bind(h);

function validateAccount(account) {
  account = account.replace(/^@+/, '');
  var regex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return regex.test(String(account).toLowerCase());
}

export default class FediverseFollowModal extends Component {
  constructor(props) {
    super(props);

    this.remoteFollowButtonPressed = this.remoteFollowButtonPressed.bind(this);

    this.state = {
      errorMessage: null,
      value: '',
      loading: false,
      valid: false,
      name: props.name,
    };
  }

  componentDidMount() {
    // initalize and display Micromodal on mount
    try {
      MicroModal.init({
        awaitCloseAnimation: false,
        awaitOpenAnimation: true, //  if using css animations to open the modal. This allows it to wait for the animation to finish before focusing on an element inside the modal.
      });
      MicroModal.show('external-actions-modal', {
        onClose: this.props.onClose,
      });
    } catch (e) {
      console.error('modal error: ', e);
    }
  }

  async remoteFollowButtonPressed() {
    this.setState({ loading: true, errorMessage: null });
    const { value } = this.state;
    const account = value.replace(/^@+/, '');
    const request = { account: account };
    const requestURL = '/api/remotefollow';
    const rawResponse = await fetch(requestURL, {
      method: 'POST',
      body: JSON.stringify(request),
    });
    const result = await rawResponse.json();

    if (!result.redirectUrl) {
      this.setState({ errorMessage: result.message, loading: false });
      return;
    }

    window.open(result.redirectUrl, '_blank');
    this.props.onClose();
  }

  onInput = (e) => {
    const { value } = e.target;
    const valid = validateAccount(value);
    this.setState({ value, valid });
  };

  render() {
    const { errorMessage, value, loading, valid, name } = this.state;
    const buttonState = valid ? '' : 'cursor-not-allowed opacity-50';
    const loadingStyle = loading ? { position: 'absolute', top: '25%', right: '50%' } : { display: 'none', top: '25%', right: '50%' };
    const error = errorMessage ? (
      html`<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
  <div class="font-bold">Error</div>
  <span class="block sm:inline">${errorMessage}</span>
</div>`
    ) : null;

    return html`
      <div class="modal micromodal-slide" id="external-actions-modal" aria-hidden="true">
        <div class="modal__overlay" tabindex="-1" data-micromodal-close>
          <div id="modal-container" class="modal__container" role="dialog" aria-modal="true"
            aria-labelledby="modal-1-title">
            <header id="modal-header"
              class="modal__header flex flex-row justify-between items-center bg-gray-300 p-3 rounded-t-lg">
              <h2 class="modal__title text-indigo-600 font-semibold">
                Follow this server.
              </h2>
              <button class="modal__close" aria-label="Close modal" data-micromodal-close></button>
            </header>

            <img src="/img/loading.gif" style=${loadingStyle} />

            <div id="modal-content-content" class="modal-content-content rounded-b-lg bg-white shadow-md px-8 pt-6 pb-8 mb-3">
              <div class="w-full">

                <p class="text-gray-700 text-lg font-semibold">
                  By following ${name} on the Fediverse you'll get notified when the stream goes live.
                </p>

                <div class="mb34">
                <label class="block text-gray-700 text-sm font-semibold mt-6" for="username">
                  Account
                </label>
                  <input onInput=${this.onInput} value="${value}"
                    class="border bg-white rounded w-full py-2 px-3 mb-2 mt-2 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
                    id="username" type="text" placeholder="account@instance.tld">
                </div>

                <p class="text-gray-600 text-xs italic">
                  You'll be redirected to your Fediverse server and asked to confirm this action.
                </p>

                <button
                  class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 mt-6 px-4 rounded focus:outline-none focus:shadow-outline ${buttonState}"
                  type="button" onClick=${this.remoteFollowButtonPressed}>
                  Follow
                </button>
                ${error}
              </div>

            </div>
          </div>
        </div>
      </div>
      </div>
    `;
  }
}




export function FediverseFollowButton({ serverName, federationInfo, onClick }) {
  const fediverseFollowAction = {
    color: 'rgba(28, 26, 59, 1)',
    description: `Follow ${serverName} at ${federationInfo.account}`,
    icon: '/img/fediverse-color.png',
    openExternally: false,
    title: `Follow ${federationInfo.followerCount > 10 ? ` (${federationInfo.followerCount})` : ''}`,
    url: "https://localhost:8080/",
  };

  const handleClick = () => onClick(fediverseFollowAction);
  return html`<span id="fediverse-follow-button-container">
    <${ExternalActionButton} onClick=${handleClick} action=${fediverseFollowAction} /></span>
    <!-- <button class="external-action-button rounded-sm flex flex-row justify-center items-center overflow-hidden bg-gray-800"
      onClick=${handleClick}>
      <span class="external-action-label">Follow on the Fediverse</span>
    </button> -->
  `;
}
