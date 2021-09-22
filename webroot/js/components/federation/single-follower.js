import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';

const html = htm.bind(h);

export default function (props) {
    const { user } = props;
    const {name, username, image} = user;

    return html`
    <div class="bg-white w-full flex items-center p-2 rounded-xl shadow border">
        <div class="relative flex items-center space-x-4">
            <img src="${image}" alt="My profile" class="w-16 h-16 rounded-full" />
        </div>
        <div class="flex-grow p-3">
            <div class="font-semibold text-gray-700">
                ${name}
            </div>
            <div class="text-sm text-gray-500">
            ${username}
            </div>
        </div>
    </div>`;
}