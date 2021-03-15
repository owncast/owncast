import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

export default function ExternalActionModal({ url, title, onClose }) {
  return html`
    <div class="fixed z-10 inset-0 overflow-y-auto">
      <div
        class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
      >
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-900 bg-opacity-50"></div>
        </div>

        <!-- This element is to trick the browser into centering the modal contents. -->
        <span
          class="hidden sm:inline-block sm:align-middle sm:h-screen"
          aria-hidden="true"
          >&#8203;</span
        >
        <div
          class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle w-2/5"
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-headline"
        >
          <div class="bg-white ">
            <div
              class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse"
            >
              <button
                onclick=${onClose}
                type="button"
                class="mt-3 rounded-md border border-gray-300 shadow-sm  bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Close
              </button>
              <div>${title}</div>
            </div>

            <!-- TODO: Show a loading spinner while the iframe loads -->
            <iframe
              style="height: 70vh"
              width="100%"
              allowpaymentrequest="true"
              allowfullscreen="false"
              sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
              src=${url}
            />
          </div>
        </div>
      </div>
    </div>
  `;
}


export function ExternalActionButton({index, action, onClick}) {
  const { title, icon, color = undefined } = action;
  const logo = icon && html`
    <span class="external-action-icon"><img src=${icon} alt="" /></span>
  `;
  const bgcolor = color && { backgroundColor: `${color}` };
  return html`
    <button
      class="external-action-button rounded-sm flex flex-row justify-center items-center overflow-hidden bg-gray-800"
      data-index=${index}
      onClick=${onClick}
      style=${bgcolor}
    >
      ${logo}
      <span class="external-action-label">${title}</span>
    </button>
  `;
}
