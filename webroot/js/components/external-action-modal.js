import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';

const html = htm.bind(h);

export default function ExternalActionModal({ url, title }) {
  const loading = 'background:url(/img/loading.gif) center center no-repeat;';

  function loaded() {
    document.querySelector('#external-modal-iframe').style = '';
  }

  return html`
   <div class="modal micromodal-slide" id="external-actions-modal" aria-hidden="true">
  <div class="modal__overlay" tabindex="-1" data-micromodal-close>
    <div id="modal-container" class="modal__container" role="dialog" aria-modal="true" aria-labelledby="modal-1-title">
      <header id="modal-header" class="modal__header">
        <h2 class="modal__title">
          ${title}
        </h2>
        <button class="modal__close" aria-label="Close modal" data-micromodal-close></button>
      </header>
      <div id="modal-content-content" class="modal-content-content">
        <div id="modal-content" class="modal__content">
          <iframe
                id="external-modal-iframe"
                style=${`${loading}`}
                width="100%"
                allowpaymentrequest="true"
                allowfullscreen="false"
                sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
                src=${url}
                onload=${loaded}
              />
        </div>
        <footer id="modal-footer" class="modal__footer">
          <button class="modal__btn" data-micromodal-close aria-label="Close this dialog window">Close</button>
        </footer>
      </div>
    </div>
  </div>
</div>

  `;
}

export function ExternalActionButton({ index, action, onClick }) {
  const { title, icon, color = undefined } = action;
  const logo =
    icon &&
    html`
      <span class="external-action-icon"><img src=${icon} alt="" /></span>
    `;
  const bgcolor = color && { backgroundColor: `${color}` };
  const handleClick = () => onClick(index);
  return html`
    <button
      class="external-action-button rounded-sm flex flex-row justify-center items-center overflow-hidden bg-gray-800"
      data-index=${index}
      onClick=${handleClick}
      style=${bgcolor}
    >
      ${logo}
      <span class="external-action-label">${title}</span>
    </button>
  `;
}
