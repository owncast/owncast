import { h, Component } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
const html = htm.bind(h);

export default function ModeratorActions(props) {
  // const { message } = props;

  return html`
    <div class="moderator-actions-group flex flex-row text-xs">
      <button type="button" class="ml-8" onClick=""></button>
    </div>
  `;
}

