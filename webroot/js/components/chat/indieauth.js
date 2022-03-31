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

  render(props, state) {
    return html`
      <div>
        <input type="url" value=${this.state.host} onInput=${this.onInput} />
        <button onclick=${this.submitButtonPressed}>Submit</button>
      </div>
    `;
  }
}
