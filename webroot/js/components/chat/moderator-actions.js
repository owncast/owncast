import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

const HIDE_MESSAGE_ICON = '🐵';
const HIDE_MESSAGE_ICON_HOVER = '🙈';
const BAN_USER_ICON = '👤';
const BAN_USER_ICON_HOVER = '🚫';

export default class ModeratorActions extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isMenuOpen: false,
    };
    this.handleOpenMenu = this.handleOpenMenu.bind(this);
    this.handleCloseMenu = this.handleCloseMenu.bind(this);
  }

  handleOpenMenu() {
    this.setState({
      isMenuOpen: true,
    });
  }

  handleCloseMenu() {
    this.setState({
      isMenuOpen: false,
    });
  }

  render() {
  const { isMenuOpen } = this.state;
  return html`
      <div class="moderator-actions-group flex flex-row text-xs p-3">
        <!-- <button type="button" class="moderator-action" onClick="" titlte="Ban user" alt="Ban this user">👤</button>

        <button type="button" class="moderator-action" onClick="" title="Hide message" alt="Hide message">🚫</button> -->

        <button type="button" class="moderator-menu-button xmoderator-action opacity-50 hover:opacity-100" onClick=${this.handleOpenMenu} title="Moderator actions" alt="Moderator actions" aria-haspopup="true" aria-controls="open-mod-actions-menu" aria-expanded=${isMenuOpen} id="open-mod-actions-button">
          <img src="/img/menu-vert.svg" alt="" />
        </button>

        ${isMenuOpen && html`<${ModeratorMenu} onDismiss=${this.handleCloseMenu} />`}
      </div>
    `;
  }
}

class ModeratorMenu extends Component {
  constructor(props) {
    super(props);
    this.menuNode = createRef();
    this.handleClickOutside = this.handleClickOutside.bind(this);
  }
  componentDidMount() {
    document.addEventListener('mousedown', this.handleClickOutside, false);
  }

  componentWillUnmount() {
    document.removeEventListener('mousedown', this.handleClickOutside, false);
  }
  handleClickOutside = e => {
    if (this.menuNode && !this.menuNode.current.contains(e.target) && this.props.onDismiss) {
      this.props.onDismiss();
    }
  };

  render() {
    return html`
      <ul
        role="menu"
        id="open-mod-actions-menu"
        aria-labelledby="open-mod-actions-button"
        class="moderator-actions-menu bg-gray-700	rounded-lg"
        ref=${this.menuNode}
      >
        <li>
          <${ModeratorMenuItem} icon=${HIDE_MESSAGE_ICON} hoverIcon=${HIDE_MESSAGE_ICON_HOVER} label="Hide message" onClick="" />
        </li>
        <li>
          <${ModeratorMenuItem} icon=${BAN_USER_ICON} hoverIcon=${BAN_USER_ICON_HOVER} label="Ban user" onClick="" />
        </li>
        <li>
          <${ModeratorMenuItem}
            icon=${html`<img src="/img/menu.svg" alt="" />`}
            label="More Info"
            onClick=""
          />
        </li>

      </ul>
    `;
  }
}

function ModeratorMenuItem({ icon, hoverIcon, label, onClick }) {

  return html`
    <button
      role="menuitem"
      type="button"
      onClick=${onClick}
      className="moderator-menu-item w-full py-2 px-4 text-white text-left whitespace-no-wrap rounded-lg hover:bg-gray-600"
    >
      ${icon && html`<span class="moderator-menu-icon menu-icon-base inline-block align-bottom mr-4">${icon}</span>`}
      <span class="moderator-menu-icon menu-icon-hover inline-block align-bottom mr-4">${hoverIcon || icon}</span>
      ${label}
    </button>
  `;
}
