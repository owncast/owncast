import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { ExternalActionButton } from './external-action-modal.js';

const html = htm.bind(h);

export default class AuthModal extends Component {
  constructor(props) {
    super(props);

    this.submitButtonPressed = this.submitButtonPressed.bind(this);

    this.state = {
      errorMessage: null,
      loading: false,
      valid: false,
    };
  }

  async submitButtonPressed() {
    const { accessToken } = this.props;
    const { host, valid } = this.state;

    if (!valid) {
      return;
    }

    const url = `/api/auth/indieauth?accessToken=${accessToken}`;
    const data = { authHost: host };

    this.setState({ loading: true });

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
      if (content.message) {
        this.setState({ errorMessage: content.message, loading: false });
        return;
      } else if (!content.redirect) {
        this.setState({
          errorMessage: 'Auth provider did not return a redirect URL.',
          loading: false,
        });
        return;
      }

      const redirect = content.redirect;
      window.location = redirect;
    } catch (e) {
      console.error(e);
      this.setState({ errorMessage: e, loading: false });
    }
  }

  onInput = (e) => {
    const { value } = e.target;
    let valid = validateURL(value);
    this.setState({ host: value, valid });
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
      : `You are already authenticated, however you can add other external sites or accounts to your chat account or log in as a different user.`;

    const error = errorMessage
      ? html` <div
          class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
          role="alert"
        >
          <div class="font-bold mb-2">There was an error.</div>
          <div class="block mt-2">
            Server error:
            <div>${errorMessage}</div>
          </div>
        </div>`
      : null;

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-4">
        <p class="text-gray-700 text-md">${message}</p>

        ${error}

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
          <p class="text-gray-700 text-lg">Authenticating.</p>
          <p class="text-gray-600 text-lg">Please wait...</p>
        </div>
      </div>
    `;
  }
}

function validateURL(url) {
  if (!url) {
    return false;
  }

  try {
    const u = new URL(url);
    if (!u) {
      return false;
    }

    if (u.protocol !== 'https:') {
      return false;
    }
  } catch (e) {
    return false;
  }

  return true;
}
