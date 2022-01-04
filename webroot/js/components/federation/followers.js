import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { URL_FOLLOWERS } from '/js/utils/constants.js';
const html = htm.bind(h);

export default class FollowerList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      followers: [],
    };
  }

  componentDidMount() {
    try {
      this.getFollowers();
    } catch (e) {
      console.error('followers error: ', e);
    }
  }

  async getFollowers() {
    const response = await fetch(URL_FOLLOWERS);
    const followers = await response.json();

    this.setState({ followers: followers });
  }

  render() {
    const { followers } = this.state;
    if (!followers) {
      return null;
    }

    const noFollowersInfo = html`<div>
      <p class="mb-5 text-2xl">No followers yet. Be the first.</p>
      <p>Info about how to follow and what it means goes here.</p>
    </div>`;

    return html`
      <div id="followers" class="p-4 w-full">
        <div
          class="grid grid-flow-row sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
        >
          ${followers.length === 0 && noFollowersInfo}
          ${followers.map((follower) => {
            return html` <${SingleFollower} user=${follower} /> `;
          })}
        </div>
      </div>
    `;
  }
}

function SingleFollower(props) {
  const { user } = props;
  const { name, username, link, image } = user;

  var displayName = name;
  var displayUsername = username;

  if (!displayName) {
    displayName = displayUsername.split('@', 1)[0];
  }
  return html`
    <a
      href=${link}
      class="follower m-3 block bg-white flex items-center p-2 rounded-xl shadow border"
      target="_blank"
    >
      <img
        src="${image || '/img/logo.svg'}"
        class="w-16 h-16 rounded-full"
        onError=${({ currentTarget }) => {
          currentTarget.onerror = null;
          currentTarget.src = '/img/logo.svg';
        }}
      />
      <div class="p-3 truncate flex-grow">
        <p class="font-semibold text-gray-700 truncate">${displayName}</p>
        <p class="text-sm text-gray-500 truncate">${displayUsername}</p>
      </div>
    </a>
  `;
}
