import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import MicroModal from '/js/web_modules/micromodal/dist/micromodal.min.js';

const html = htm.bind(h);

export default class ExternalActionModal extends Component {
  constructor(props) {
    super(props);
    this.state = {
      iframeLoaded: false,
    };

    this.setIframeLoaded = this.setIframeLoaded.bind(this);
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

  setIframeLoaded() {
    this.setState({
      iframeLoaded: true,
    });
  }

  render() {
    const { action, useIframe = true, customContent = null } = this.props;
    const { url, title, description } = action;
    const { iframeLoaded } = this.state;
    const iframeStyle = iframeLoaded
      ? null
      : { backgroundImage: 'url(/img/loading.gif)' };

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
              <h2
                id="external-action-modal-header"
                class="modal__title text-indigo-600 font-semibold"
              >
                ${title || description}
              </h2>
              <button
                class="modal__close"
                aria-label="Close modal"
                data-micromodal-close
              ></button>
            </header>
            <div
              id="modal-content-content"
              class="modal-content-content rounded-b-md"
            >
              ${useIframe
                ? html`
                    <div
                      id="modal-content"
                      class="modal__content text-gray-600 overflow-y-auto overflow-x-hidden"
                    >
                      <iframe
                        id="external-modal-iframe"
                        style=${{ iframeStyle }}
                        class="bg-gray-100 bg-center bg-no-repeat"
                        width="100%"
                        allowpaymentrequest="true"
                        allowfullscreen="false"
                        sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
                        src=${url}
                        onload=${this.setIframeLoaded}
                      />
                    </div>
                  `
                : customContent}
            </div>
          </div>
        </div>
      </div>
    `;
  }
}

export function ExternalActionButton({ action, onClick, label = '' }) {
  const { title, icon, color = undefined, description } = action;
  const logo =
    icon &&
    html`
      <span class="external-action-icon"><img src=${icon} alt="" /></span>
    `;
  const bgcolor = color && { backgroundColor: `${color}` };
  const handleClick = () => onClick(action);
  return html`
    <button
      class="external-action-button rounded-sm flex flex-row justify-center items-center overflow-hidden m-1 px-3 py-1 text-base text-white bg-gray-800 rounded"
      onClick=${handleClick}
      style=${bgcolor}
      aria-label=${description}
      title=${description || title}
    >
      ${logo}
      <span class="external-action-label">${label || title}</span>
    </button>
  `;
}
