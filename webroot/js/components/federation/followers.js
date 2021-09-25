import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
// import SingleFollower from './single-follower.js';

const html = htm.bind(h);

export default class FollowerList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      followers: []
    }
  }

  componentDidMount() {
    try {
      this.getFollowers();
    } catch (e) {
      console.error('followers error: ', e);
    }
  }

  async getFollowers() {
    const response = await fetch("https://randomuser.me/api/?results=50");
    const result = await response.json();
    const followers = result.results.map(user => {
      return {
        name: `${user.name.first} ${user.name.last}`,
        username: user.email,
        image: user.picture.thumbnail,
        link: 'https://botsin.space/@owncast'
      };
    });
    this.setState({ followers: followers });
  }

  render() {
    const { followers } = this.state;
    return html`
      <div id="followers" class="p-4 w-full">
        <h3 class="text-3xl font-semibold mb-4">Followers</h3>
        <div class="grid grid-flow-row sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 sm:grid-cols-2">
          ${
            followers.map(follower => {
              return html`
                <${SingleFollower} user=${follower} />
              `;
            })
          }
        </div>
      </div>
    `;
  };
}


function SingleFollower(props) {
  const { user } = props;
  const { name, username, image, link } = user;

  return html`
    <a href=${link} class="follower m-3 block bg-white flex  items-center p-2 rounded-xl shadow border" target="_blank">
      <img src="${image}" alt="My profile" class="w-16 h-16 rounded-full" />
      <div class="p-3 truncate flex-grow">
        <p class="font-semibold text-gray-700 truncate">${name}</p>
        <p class="text-sm text-gray-500 truncate">${username}</p>
      </div>
    </a>
  `;
}
