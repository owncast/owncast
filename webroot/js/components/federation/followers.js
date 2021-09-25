import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import SingleFollower from './single-follower.js';

const html = htm.bind(h);

function followerClicked(link) {
  window.open(link, '_blank');
}

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

    const userViews = followers.map(follower => {
      const {link} = follower;
      return html`
        <div class="w-1/4 m-3 cursor-pointer" onClick=${() => followerClicked(link)}>
          <${SingleFollower} user=${follower} />
        </div>
      `;
    });

    return html`
      <div class="flex flex-wrap -mb-4">
        ${userViews}
        <div class="p-12 text-center" style=${{ height: '300px', width: '100%'}}>
          things in stuff yo
        </div>
      </div>
    `;
  };
}
