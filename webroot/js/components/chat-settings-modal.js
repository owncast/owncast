import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import TabBar from './tab-bar.js';
import IndieAuthForm from './auth-indieauth.js';

const html = htm.bind(h);

export default class ChatSettingsModal extends Component {
  render() {
    const {
      accessToken,
      authenticated,
      username,
      onUsernameChange,
      indieAuthEnabled,
    } = this.props;

    const TAB_CONTENT = [
      {
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
        />`,
      },
    ];

    return html`
      <div class="bg-gray-100 bg-center bg-no-repeat p-5">
        <${TabBar} tabs=${TAB_CONTENT} ariaLabel="Chat settings" />
      </div>
    `;
  }
}
