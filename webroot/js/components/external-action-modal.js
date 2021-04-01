import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';

const html = htm.bind(h);

export default function ExternalActionModal({ url, title, onClose }) {
  const loading = 'background:url(/img/loading.gif) center center no-repeat;';

  function loaded() {
    document.querySelector('#external-modal-iframe').style = '';
  }

  return html`
    <div class="fixed inset-0 overflow-y-auto" style="z-index: 9999">
      <div
        class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
      >
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div
            onClick=${() => onClose()}
            class="absolute inset-0 bg-gray-900 bg-opacity-75"
          ></div>
        </div>

        <!-- This element is to trick the browser into centering the modal contents. -->
        <span
          class="hidden sm:inline-block sm:align-middle sm:h-screen"
          aria-hidden="true"
          >&#8203;</span
        >
        <div
          class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:align-middle w-screen md:max-w-2xl lg:max-w-2xl"
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-headline"
        >
          <div class="bg-white ">
            <div
              class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex justify-between items-center"
            >
              <h3 class="font-bold hidden md:block">${title}</h3>
              <span class="" onclick=${onClose}>
                <svg
                  class="h-12 w-12 fill-current text-grey hover:text-grey-darkest"
                  role="button"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 20 20"
                >
                  <title>Close</title>
                  <path
                    d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"
                  />
                </svg>
              </span>
            </div>

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
