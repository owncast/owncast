import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import TabBar from './tab-bar.js';
import IndieAuthForm from './auth-indieauth.js';
import FediverseAuth from './auth-fediverse.js';

const html = htm.bind(h);

export default class ChatSettingsModal extends Component {
  render() {
    const {
      accessToken,
      authenticated,
      federationEnabled,
      username,
      indieAuthEnabled,
    } = this.props;

    const TAB_CONTENT = [];
    if (indieAuthEnabled) {
      TAB_CONTENT.push({
        label: html`<span style=${{ display: 'flex', alignItems: 'center' }}
          ><img
            style=${{
              display: 'inline',
              height: '0.8em',
              marginRight: '5px',
            }}
            src="/img/indieauth.png"
          />
          IndieAuth</span
        >`,
        content: html`<${IndieAuthForm}}
          accessToken=${accessToken}
          authenticated=${authenticated}
          username=${username}
        />`,
      });
    }

    if (federationEnabled) {
      TAB_CONTENT.push({
        label: html`<span style=${{ display: 'flex', alignItems: 'center' }}
          ><img
            style=${{
              display: 'inline',
              height: '0.8em',
              marginRight: '5px',
            }}
            src="/img/fediverse-black.png"
          />
          FediAuth</span
        >`,
        content: html`<${FediverseAuth}}
          authenticated=${authenticated}
          accessToken=${accessToken}
          authenticated=${authenticated}
          username=${username}
        />`,
      });
    }

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-5">
        <${TabBar} tabs=${TAB_CONTENT} ariaLabel="Chat settings" />
        <p class="mt-4">
          <b>Note:</b> This is for authentication purposes only, and no personal
          information will be accessed or stored.
        </p>
      </div>
    `;
  }
}
