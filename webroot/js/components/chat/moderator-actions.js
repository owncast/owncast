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
  const { message } = this.props;

  return html`
      <div class="moderator-actions-group flex flex-row text-xs p-3">
        <button type="button" class="moderator-menu-button" onClick=${this.handleOpenMenu} title="Moderator actions" alt="Moderator actions" aria-haspopup="true" aria-controls="open-mod-actions-menu" aria-expanded=${isMenuOpen} id="open-mod-actions-button">
          <img src="/img/menu-vert.svg" alt="" />
        </button>

        ${isMenuOpen && html`<${ModeratorMenu} message=${message} onDismiss=${this.handleCloseMenu} />`}
      </div>
    `;
  }
}

class ModeratorMenu extends Component {
  constructor(props) {
    super(props);
    this.menuNode = createRef();

    this.state = {
      displayMoreInfo: false,
    };
    this.handleClickOutside = this.handleClickOutside.bind(this);
    this.handleToggleMoreInfo = this.handleToggleMoreInfo.bind(this);
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

  handleToggleMoreInfo() {
    console.log(this.state.displayMoreInfo)
    this.setState({
      displayMoreInfo: !this.state.displayMoreInfo,
    });
  }

  render() {
    const { message } = this.props;
    const { displayMoreInfo } = this.state;
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
            onClick=${this.handleToggleMoreInfo}
          />
        </li>
        ${displayMoreInfo && html`<${ModeratorMoreInfoContainer} message=${message} />`}
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
      ${icon && html`<span className="moderator-menu-icon menu-icon-base inline-block align-bottom mr-4">${icon}</span>`}
      <span className="moderator-menu-icon menu-icon-hover inline-block align-bottom mr-4">${hoverIcon || icon}</span>
      ${label}
    </button>
  `;
}


function ModeratorMoreInfoContainer({ message }) {
  const { user, timestamp, body } = message;
  const { displayName, createdAt, previousNames,
    // mock field
    isModerator: isAuthorModerator = true,
    } = user;

  const createDate = new Date(createdAt);
  const sentDate = new Date(timestamp);
  return html`
    <div className="moderator-more-info-container text-gray-300 bg-gray-800 rounded-lg">
      <div className="moderator-more-info-message pb-2">
        <p className="text-xs text-gray-400">Sent at ${sentDate.toLocaleTimeString()}</p>
        <div className="text-sm" dangerouslySetInnerHTML=${{ __html: body }} />
      </div>

      <div className="moderator-more-info-user py-2 my-2">
        <p className="${isAuthorModerator && ' moderator-flag'}">${displayName}</p>

        <p>Created at: ${createDate.toLocaleTimeString()}</p>

        ${previousNames.length > 1 &&
          html`
            <p>Previous Names:</p>
            ${previousNames.join()}
          `
        }

      </div>
      <div className="moderator-more-info-actions pt-2">
        <${ModeratorMenuItem} icon=${HIDE_MESSAGE_ICON} hoverIcon=${HIDE_MESSAGE_ICON_HOVER} label="Hide message" onClick="" />
        <${ModeratorMenuItem} icon=${BAN_USER_ICON} hoverIcon=${BAN_USER_ICON_HOVER} label="Ban user" onClick="" />
      </div>
    </div>
  `;
}
