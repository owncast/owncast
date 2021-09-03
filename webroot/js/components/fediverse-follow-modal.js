import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import MicroModal from '/js/web_modules/micromodal/dist/micromodal.min.js';

const html = htm.bind(h);

export default class FediverseFollowModal extends Component {
  constructor(props) {
    super(props);

    this.remoteFollowButtonPressed = this.remoteFollowButtonPressed.bind(this);

    this.state = {
      errorMessage: null,
      value: '',
    };
  }

  componentDidMount() {
    // initalize and display Micromodal on mount
    try {
      MicroModal.init({
        awaitCloseAnimation: false,
        awaitOpenAnimation: true, //  if using css animations to open the modal. This allows it to wait for the animation to finish before focusing on an element inside the modal.
      });
      MicroModal.show('external-actions-modal', {
        onClose: this.props.onClose,
      });
    } catch (e) {
      console.error('modal error: ', e);
    }
  }

  async remoteFollowButtonPressed() {
    const { value } = this.state;
    const request = { account: value };
    const requestURL = '/api/remotefollow';
    const rawResponse = await fetch(requestURL, {
      method: 'POST',
      body: JSON.stringify(request),
    });
    const result = await rawResponse.json();
    console.log(result);

    if (!result.redirectUrl) {
      // error should be displayed
      return;
    }

    window.open(result.redirectUrl, '_blank');
    this.props.onClose();
  }

  onInput = (e) => {
    const { value } = e.target;
    this.setState({ value });
  };

  render() {
    const { errorMessage, value } = this.state;

    return html`
      <div
        class="modal micromodal-slide"
        id="external-actions-modal"
        aria-hidden="true"
      >
        <div class="modal__overlay" tabindex="-1" data-micromodal-close>
          <div
            id="modal-container"
            class="modal__container rounded-md"
            role="dialog"
            aria-modal="true"
            aria-labelledby="modal-1-title"
          >
            <header
              id="modal-header"
              class="modal__header flex flex-row justify-between items-center bg-gray-300 p-3 rounded-t-md"
            >
              <h2 class="modal__title text-indigo-600 font-semibold">
                Follow this server.
              </h2>
              <button
                class="modal__close"
                aria-label="Close modal"
                data-micromodal-close
              ></button>
            </header>
            <div id="modal-content-content" class="modal-content-content">
              <div
                id="modal-content"
                class="modal__content text-gray-600 rounded-b-md overflow-y-auto overflow-x-hidden"
              >
                Input form asking you to type in your fediverse account goes
                here.

                <input type="text" value="${value}" onInput="${this.onInput}" />

                <button onClick=${this.remoteFollowButtonPressed}>
                  Follow
                </button>

                <div>${errorMessage}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    `;
  }
}

export function FediverseFollowButton({ action, onClick }) {
  const handleClick = () => onClick(action);
  return html`
    <button
      class="external-action-button rounded-sm flex flex-row justify-center items-center overflow-hidden bg-gray-800"
      onClick=${handleClick}
    >
      <span class="external-action-label">Follow on the Fediverse</span>
    </button>
  `;
}
