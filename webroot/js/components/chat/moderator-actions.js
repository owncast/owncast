import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { textColorForHue } from '../../utils/user-colors.js';
import { URL_BAN_USER, URL_HIDE_MESSAGE } from '../../utils/constants.js';

const html = htm.bind(h);

const HIDE_MESSAGE_ICON = `/img/hide-message-grey.svg`;
const HIDE_MESSAGE_ICON_HOVER = '/img/hide-message.svg';
const BAN_USER_ICON = '/img/ban-user-grey.svg';
const BAN_USER_ICON_HOVER = '/img/ban-user.svg';

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
    const { message, accessToken } = this.props;
    const { id } = message;
    const { user } = message;

    return html`
      <div class="moderator-actions-group flex flex-row text-xs p-3">
        <button
          type="button"
          class="moderator-menu-button"
          onClick=${this.handleOpenMenu}
          title="Moderator actions"
          alt="Moderator actions"
          aria-haspopup="true"
          aria-controls="open-mod-actions-menu"
          aria-expanded=${isMenuOpen}
          id="open-mod-actions-button"
        >
          <img src="/img/menu-vert.svg" alt="" />
        </button>

        ${isMenuOpen &&
        html`<${ModeratorMenu}
          message=${message}
          onDismiss=${this.handleCloseMenu}
          accessToken=${accessToken}
          id=${id}
          userId=${user.id}
        />`}
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
    this.handleBanUser = this.handleBanUser.bind(this);
    this.handleHideMessage = this.handleHideMessage.bind(this);
  }

  componentDidMount() {
    document.addEventListener('mousedown', this.handleClickOutside, false);
  }

  componentWillUnmount() {
    document.removeEventListener('mousedown', this.handleClickOutside, false);
  }

  handleClickOutside = (e) => {
    if (
      this.menuNode &&
      !this.menuNode.current.contains(e.target) &&
      this.props.onDismiss
    ) {
      this.props.onDismiss();
    }
  };

  handleToggleMoreInfo() {
    this.setState({
      displayMoreInfo: !this.state.displayMoreInfo,
    });
  }

  async handleHideMessage() {
    if (!confirm('Are you sure you want to remove this message from chat?')) {
      this.props.onDismiss();
      return;
    }

    const { accessToken, id } = this.props;
    const url = new URL(location.origin + URL_HIDE_MESSAGE);
    url.searchParams.append('accessToken', accessToken);
    const hideMessageUrl = url.toString();

    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ idArray: [id] }),
    };

    try {
      await fetch(hideMessageUrl, options);
    } catch (e) {
      console.error(e);
    }

    this.props.onDismiss();
  }

  async handleBanUser() {
    if (!confirm('Are you sure you want to remove this user from chat?')) {
      this.props.onDismiss();
      return;
    }

    const { accessToken, userId } = this.props;
    const url = new URL(location.origin + URL_BAN_USER);
    url.searchParams.append('accessToken', accessToken);
    const hideMessageUrl = url.toString();

    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ userId: userId }),
    };

    try {
      await fetch(hideMessageUrl, options);
    } catch (e) {
      console.error(e);
    }

    this.props.onDismiss();
  }

  render() {
    const { message } = this.props;
    const { displayMoreInfo } = this.state;
    return html`
      <ul
        role="menu"
        id="open-mod-actions-menu"
        aria-labelledby="open-mod-actions-button"
        class="moderator-actions-menu bg-gray-700	rounded-lg shadow-md"
        ref=${this.menuNode}
      >
        <li>
          <${ModeratorMenuItem}
            icon=${HIDE_MESSAGE_ICON}
            hoverIcon=${HIDE_MESSAGE_ICON_HOVER}
            label="Hide message"
            onClick="${this.handleHideMessage}"
          />
        </li>
        <li>
          <${ModeratorMenuItem}
            icon=${BAN_USER_ICON}
            hoverIcon=${BAN_USER_ICON_HOVER}
            label="Ban user"
            onClick="${this.handleBanUser}"
          />
        </li>
        <li>
          <${ModeratorMenuItem}
            icon="/img/menu.svg"
            label="More Info"
            onClick=${this.handleToggleMoreInfo}
          />
        </li>
        ${displayMoreInfo &&
        html`<${ModeratorMoreInfoContainer}
          message=${message}
          handleBanUser=${this.handleBanUser}
          handleHideMessage=${this.handleHideMessage}
        />`}
      </ul>
    `;
  }
}

// 3 dots button
function ModeratorMenuItem({ icon, hoverIcon, label, onClick }) {
  return html`
    <button
      role="menuitem"
      type="button"
      onClick=${onClick}
      className="moderator-menu-item w-full py-2 px-4 text-white text-left whitespace-no-wrap rounded-lg hover:bg-gray-600"
    >
      ${icon &&
      html`<span
        className="moderator-menu-icon menu-icon-base inline-block align-bottom mr-4"
        ><img src="${icon}"
      /></span>`}
      <span
        className="moderator-menu-icon menu-icon-hover inline-block align-bottom mr-4"
        ><img src="${hoverIcon || icon}"
      /></span>
      ${label}
    </button>
  `;
}

// more details panel that display message, prev usernames, actions
function ModeratorMoreInfoContainer({
  message,
  handleHideMessage,
  handleBanUser,
}) {
  const { user, timestamp, body } = message;
  const { displayName, createdAt, previousNames, displayColor } = user;
  const isAuthorModerator = user.scopes && user.scopes.includes('MODERATOR');

  const authorTextColor = { color: textColorForHue(displayColor) };
  const createDate = new Date(createdAt);
  const sentDate = new Date(timestamp);
  return html`
    <div
      className="moderator-more-info-container text-gray-300 bg-gray-800 rounded-lg p-4 border border-white text-base absolute"
    >
      <div
        className="moderator-more-info-message scrollbar-hidden bg-gray-700 rounded-md pb-2"
      >
        <p className="text-xs text-gray-500">
          Sent at ${sentDate.toLocaleTimeString()}
        </p>
        <div className="text-sm" dangerouslySetInnerHTML=${{ __html: body }} />
      </div>
      <div className="moderator-more-info-user py-2 my-2">
        <p className="text-xs text-gray-500">Sent by:</p>
        <p
          className="font-bold ${isAuthorModerator && ' moderator-flag'}"
          style=${authorTextColor}
        >
          ${displayName}
        </p>

        <p className="text-xs text-gray-500 mt-2">
          First joined: ${createDate.toLocaleString()}
        </p>

        ${previousNames.length > 1 &&
        html`
          <p className="text-xs text-gray-500 my-1">
            Previously known as: ${' '}
            <span className="text-white text-gray-400"
              >${previousNames.join(', ')}</span
            >
          </p>
        `}
      </div>
      <div
        className="moderator-more-info-actions pt-2 flex flex-row border-t border-gray-700 shadow-md"
      >
        <${handleHideMessage && ModeratorMenuItem}
          icon=${HIDE_MESSAGE_ICON}
          hoverIcon=${HIDE_MESSAGE_ICON_HOVER}
          label="Hide message"
          onClick="${handleHideMessage}"
        />
        <${handleBanUser && ModeratorMenuItem}
          icon=${BAN_USER_ICON}
          hoverIcon=${BAN_USER_ICON_HOVER}
          label="Ban user"
          onClick="${handleBanUser}"
        />
      </div>
    </div>
  `;
}
