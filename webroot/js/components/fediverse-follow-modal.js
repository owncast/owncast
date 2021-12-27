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
            <div class="font-bold">Error</div>
            <span class="block sm:inline">${errorMessage}</span>
          </div>
        `
      : null;

    const loaderStyle = loading ? 'flex' : 'none';

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-4">
        <p class="text-gray-700 text-md">
          By following on the ${' '}
          <a
            class=" text-blue-500"
            href="https://en.wikipedia.org/wiki/Fediverse"
            target="_blank"
            rel="noopener noreferrer"
            >Fediverse</a
          >
          ${' '} you'll get notified when the stream goes live.
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
            placeholder="account@instance.tld"
          />
        </div>

        <p class="text-gray-600 text-xs italic">
          You'll be redirected to your Fediverse server and asked to confirm
          this action.
        </p>

        <button
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 mt-6 px-4 rounded focus:outline-none focus:shadow-outline ${buttonState}"
          type="button"
          onClick=${this.remoteFollowButtonPressed}
        >
          Follow
        </button>
        ${error}
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
    title: `Follow ${serverName} on the Fediverse`,
    url: '',
  };

  const handleClick = () => onClick(fediverseFollowAction);
  return html`
    <span id="fediverse-follow-button-container">
      <${ExternalActionButton}
        onClick=${handleClick}
        action=${fediverseFollowAction}
        label=${`Follow ${
          federationInfo.followerCount > 10
            ? ` (${federationInfo.followerCount})`
            : ''
        }`}
      />
    </span>
  `;
}
