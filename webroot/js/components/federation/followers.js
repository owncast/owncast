import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { URL_FOLLOWERS } from '/js/utils/constants.js';
const html = htm.bind(h);
import { paginateArray } from '../../utils/helpers.js';
export default class FollowerList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      followers: [],
      followersPage: 0,
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

    this.setState({
      followers: followers,
    });
  }

  changeFollowersPage(page) {
    this.setState({ followersPage: page });
  }

  render() {
    const FOLLOWER_PAGE_SIZE = 16;
    const { followersPage } = this.state;

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

    const paginatedFollowers = paginateArray(
      followers,
      followersPage + 1,
      FOLLOWER_PAGE_SIZE
    );

    const paginationControls =
      paginatedFollowers.totalPages > 1 &&
      Array(paginatedFollowers.totalPages)
        .fill()
        .map((x, n) => {
          const activePageClass =
            n === followersPage &&
            'bg-indigo-600 rounded-full shadow-md focus:shadow-md text-white';
          return html` <li class="page-item active">
            <a
              class="page-link relative block cursor-pointer hover:no-underline py-1.5 px-3 border-0 rounded-full hover:text-gray-800 hover:bg-gray-200 outline-none transition-all duration-300 ${activePageClass}"
              onClick=${() => this.changeFollowersPage(n)}
            >
              ${n + 1}
            </a>
          </li>`;
        });

    return html`
      <div>
        <div class="flex flex-wrap">
          ${followers.length === 0 && noFollowersInfo}
          ${paginatedFollowers.items.map((follower) => {
            return html` <${SingleFollower} user=${follower} /> `;
          })}
        </div>
        <div class="flex">
          <nav aria-label="Page navigation example">
            <ul class="flex list-style-none">
              ${paginationControls}
            </ul>
          </nav>
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
      class="following-list-follower block bg-white flex p-2 rounded-xl shadow border hover:no-underline mb-3 mr-3"
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
