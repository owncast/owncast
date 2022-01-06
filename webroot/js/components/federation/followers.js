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
      <p class="mb-5 text-2xl">Be the first to follow this live stream.</p>
      <p class="text-md">
        By following this stream you'll get updates when it goes live, receive
        posts from the streamer, and be featured here as a follower.
      </p>
      <p class="text-md mt-5">
        Learn more about ${' '}
        <a class="underline" href="https://en.wikipedia.org/wiki/Fediverse"
          >The Fediverse</a
        >, where you can follow this server as well as so much more.
      </p>
    </div>`;

    return html`
      <div class="flex flex-wrap">
        ${followers.length === 0 && noFollowersInfo}
        ${followers.map((follower) => {
          return html` <${SingleFollower} user=${follower} /> `;
        })}
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
      class="following-list-follower block bg-white flex p-2 rounded-xl shadow border hover:no-underline m-4"
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
