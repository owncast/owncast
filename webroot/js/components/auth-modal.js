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

export default class AuthModal extends Component {
  constructor(props) {
    super(props);

    this.submitButtonPressed = this.submitButtonPressed.bind(this);

    this.state = {
      errorMessage: null,
      loading: false,
      valid: true,
    };
  }

  async submitButtonPressed() {
    const { accessToken } = this.props;
    const { host } = this.state;

    const url = `/api/auth/indieauth?accessToken=${accessToken}`;
    const data = { authHost: host };
    console.log(data);
    try {
      const rawResponse = await fetch(url, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      const content = await rawResponse.json();
      const redirect = content.redirect;
      window.location = redirect;
    } catch (e) {
      console.error(e);
    }
  }

  onInput = (e) => {
    const { value } = e.target;
    this.setState({ host: value });
  };

  render() {
    const { errorMessage, host, valid, loading } = this.state;
    const { authenticated } = this.props;
    const buttonState = valid ? '' : 'cursor-not-allowed opacity-50';

    const loaderStyle = loading ? 'flex' : 'none';

    const message = !authenticated
      ? `While you can chat completely anonymously you can also add
    authentication so you can rejoin with the same chat persona from any
    device or browser.`
      : `You are already authenticated, however, if you'd like, you can add other external sites or accounts to your chat account.`;

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-4">
        <p class="text-gray-700 text-md">${message}</p>

        ${authenticated}

        <div class="mb34">
          <label
            class="block text-gray-700 text-sm font-semibold mt-6"
            for="username"
          >
            IndieAuth with your own site
          </label>
          <input
            onInput=${this.onInput}
            type="url"
            value=${host}
            class="border bg-white rounded w-full py-2 px-3 mb-2 mt-2 text-indigo-700 leading-tight focus:outline-none focus:shadow-outline"
            id="username"
            type="text"
            placeholder="https://yoursite.com"
          />
          <button
            class="bg-indigo-500 hover:bg-indigo-600 text-white font-bold py-2 mt-6 px-4 rounded focus:outline-none focus:shadow-outline ${buttonState}"
            type="button"
            onClick=${this.submitButtonPressed}
          >
            Authenticate with IndieAuth
          </button>
          <p>
            Learn more about
            <a class="underline" href="https://indieauth.net/">IndieAuth</a>.
          </p>
        </div>

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
