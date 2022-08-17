import { h } from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import { classNames } from '../utils/helpers.js';

const html = htm.bind(h);

function SocialIcon(props) {
  const { platform, icon, url } = props;
  const iconSupplied = !!icon;
  const name = platform;

  const finalIcon = iconSupplied ? icon : '/img/platformlogos/default.svg';

  const style = `background-image: url(${finalIcon});`;

  const itemClass = classNames({
    'user-social-item': true,
    flex: true,
    'justify-start': true,
    'items-center': true,
    'm-1': true,
  });
  const labelClass = classNames({
    'platform-label': true,
    'visually-hidden': !!finalIcon,
    'text-indigo-800': true,
    'text-xs': true,
    uppercase: true,
    'max-w-xs': true,
    'inline-block': true,
  });

  return html`
    <a class=${itemClass} target="_blank" rel="me" href=${url}>
      <span
        class="platform-icon bg-no-repeat"
        style=${style}
        title="Find me on ${name}"
      ></span>
      <span class=${labelClass}>Find me on ${name}</span>
    </a>
  `;
}

export default function (props) {
  const { handles } = props;
  if (handles == null) {
    return null;
  }

  const list = handles.map(
    (item, index) => html`
      <li key="social${index}">
        <${SocialIcon}
          platform=${item.platform}
          icon=${item.icon}
          url=${item.url}
        />
      </li>
    `
  );

  return html` <ul id="social-list" class="social-list m-2 text-center">
    <p
      class="follow-icon-list flex flex-row items-center justify-center flex-wrap"
    >
      ${list}
    </p>
  </ul>`;
}
