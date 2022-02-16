import { h, Component, createRef } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

export default class IndieAuthForm extends Component {
  constructor(props, context) {
    super(props, context);

    this.state = {};

    // this.onChange = this.onChange.bind(this);
    this.submitButtonPressed = this.submitButtonPressed.bind(this);
  }

  onInput = (e) => {
    const { value } = e.target;
    this.setState({ host: value });
  };

  async submitButtonPressed() {
    const { accessToken } = this.props;
    const { host } = this.state;

    if (!accessToken || !host) {
      return;
    }

    const url = `/api/auth/indieauth?accessToken=${accessToken}`;
    const data = { authHost: host };
    console.log(data);
    try {
      const rawResponse = await fetch(url, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      const content = await rawResponse.json();
      const redirect = content.redirect;
      window.location = redirect;
    } catch (e) {
      console.error(e);
    }
  }

  render(props, state) {
    return html`
      <div>
        <input type="url" value=${this.state.host} onInput=${this.onInput} />
        <button onclick=${this.submitButtonPressed}>Submit</button>
      </div>
    `;
  }
}
