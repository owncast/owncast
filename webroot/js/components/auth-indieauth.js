import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

export default class IndieAuthForm extends Component {
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

      if (content.redirect) {
        const redirect = content.redirect;
        window.location = redirect;
      }
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
    const { errorMessage, loading, host, valid } = this.state;
    const { authenticated, username } = this.props;
    const buttonState = valid ? '' : 'cursor-not-allowed opacity-50';
    const loaderStyle = loading ? 'flex' : 'none';

    const message = !authenticated
      ? html`Use your own domain to authenticate ${' '}
          <span class="font-bold">${username}</span> or login as a previously
          ${' '} authenticated chat user using IndieAuth.`
      : html`<span
          ><b>You are already authenticated</b>. However, you can add other
          domains or log in as a different user.</span
        >`;

    let errorMessageText = errorMessage;
    if (!!errorMessageText) {
      if (errorMessageText.includes('url does not support indieauth')) {
        errorMessageText =
          'The provided URL is either invalid or does not support IndieAuth.';
      }
    }

    const error = errorMessage
      ? html` <div
          class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative"
          role="alert"
        >
          <div class="font-bold mb-2">There was an error.</div>
          <div class="block mt-2">
            <div>${errorMessageText}</div>
          </div>
        </div>`
      : null;

    return html` <div>
      <p class="text-gray-700">${message}</p>

      <p>${error}</p>

      <div class="mb34">
        <label
          class="block text-gray-700 text-sm font-semibold mt-6"
          for="username"
        >
          Your domain
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
          Authenticate with your domain
        </button>
      </div>

      <p class="mt-4">
        <details>
          <summary class="cursor-pointer">
            Learn more about using IndieAuth to authenticate with chat.
          </summary>
          <div class="inline">
            <p class="mt-4">
              IndieAuth allows for a completely independent and decentralized
              way of identifying yourself using your own domain.
            </p>

            <p class="mt-4">
              If you run an Owncast instance, you can use that domain here.
              Otherwise, ${' '}
              <a class="underline" href="https://indieauth.net/#providers"
                >learn more about how you can support IndieAuth</a
              >.
            </p>
          </div>
        </details>
      </p>

      <div
        id="follow-loading-spinner-container"
        style="display: ${loaderStyle}"
      >
        <img id="follow-loading-spinner" src="/img/loading.gif" />
        <p class="text-gray-700 text-lg">Authenticating.</p>
        <p class="text-gray-600 text-lg">Please wait...</p>
      </div>
    </div>`;
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
