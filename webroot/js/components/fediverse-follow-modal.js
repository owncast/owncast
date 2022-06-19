import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { ExternalActionButton } from './external-action-modal.js';

const html = htm.bind(h);

function validateAccount(account) {
  account = account.replace(/^@+/, '');
  var regex =
    /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
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
    };
  }

  async remoteFollowButtonPressed() {
    if (!this.state.valid) {
      return;
    }

    this.setState({ loading: true, errorMessage: null });
    const { value } = this.state;
    const { onClose } = this.props;

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
    onClose();
  }

  navigateToFediverseJoinPage() {
    window.open('https://owncast.online/join-fediverse', '_blank');
  }

  onInput = (e) => {
    const { value } = e.target;
    const valid = validateAccount(value);
    this.setState({ value, valid });
  };

  render() {
    const { name, federationInfo = {}, logo } = this.props;
    const { account } = federationInfo;
    const { errorMessage, value, valid, loading } = this.state;
    const buttonState = valid ? '' : 'cursor-not-allowed opacity-50';

    const error = errorMessage
      ? html`
          <div
            class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
            role="alert"
          >
            <div class="font-bold mb-2">
              There was an error following this Owncast server.
            </div>
            <span class="block">
              Please verify you entered the correct user account. It's also
              possible your server may not support remote following, so you may
              want to manually follow ${' '}
              <span class="font-semibold">${account}</span> using your service's
              own interface.
            </span>
            <div class="block mt-2">
              Server error: <span class="">${errorMessage}</span>
            </div>
          </div>
        `
      : null;

    const loaderStyle = loading ? 'flex' : 'none';

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-4">
        <p class="text-gray-700 text-md">
          By following this stream on the Fediverse you'll receive updates when
          it goes live, get posts from the streamer, and be featured as a
          follower.
        </p>

        <div
          class="p-4 my-2 rounded bg-gray-300 border border-indigo-400 border-solid flex items-center justify-start"
        >
          <img src=${logo} style=${{ height: '3em', width: '3em' }} />
          <p class="ml-4">
            <span class="font-bold">${name}</span>
            <br />
            <span class="">${account}</span>
          </p>
        </div>

        ${error}

        <div class="mb34">
          <label
            class="block text-gray-700 text-sm font-semibold mt-6"
            for="username"
          >
            Enter your username@server to follow:
          </label>
          <input
            onInput=${this.onInput}
            value="${value}"
            class="border bg-white rounded w-full py-2 px-3 mb-2 mt-2 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            id="username"
            type="text"
            placeholder="Fediverse account@instance.tld"
          />
        </div>

        <p class="text-gray-600 text-xs italic">
          You'll be redirected to your Fediverse server and asked to confirm
          this action. ${'  '}
          <a
            class=" text-blue-500"
            href="https://owncast.online/join-fediverse"
            target="_blank"
            rel="noopener noreferrer"
            >Join the Fediverse if you haven't.</a
          >
        </p>

        <button
          class="bg-indigo-500 hover:bg-indigo-600 text-white font-bold py-2 mt-6 px-4 rounded focus:outline-none focus:shadow-outline ${buttonState}"
          type="button"
          onClick=${this.remoteFollowButtonPressed}
        >
          Follow
        </button>
        <button
          class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 ml-4 mt-6 px-4 rounded focus:outline-none focus:shadow-outline"
          type="button"
          onClick=${this.navigateToFediverseJoinPage}
        >
          Join the Fediverse
        </button>
        <div
          id="follow-loading-spinner-container"
          style="display: ${loaderStyle}"
        >
          <img id="follow-loading-spinner" src="/img/loading.gif" />
          <p class="text-gray-700 text-lg">Contacting your server.</p>
          <p class="text-gray-600 text-lg">Please wait...</p>
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
    title: `Follow ${serverName}`,
    url: '',
  };

  const handleClick = () => onClick(fediverseFollowAction);
  return html`
    <span id="fediverse-follow-button-container">
      <${ExternalActionButton}
        onClick=${handleClick}
        action=${fediverseFollowAction}
        label="Follow"
      />
    </span>
  `;
}
