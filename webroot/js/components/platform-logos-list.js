import {h} from '/js/web_modules/preact.js';
import htm from '/js/web_modules/htm.js';
import {SOCIAL_PLATFORMS} from '../utils/platforms.js';
import {classNames} from '../utils/helpers.js';

const html = htm.bind(h);

function SocialIcon(props) {
  const {platform, icon, url} = props;
  const platformInfo = SOCIAL_PLATFORMS[platform.toLowerCase()];
  const inList = !!platformInfo;
  const iconSupplied = !!icon
  const name = inList ? platformInfo.name : platform;

  const finalIcon = iconSupplied ? icon : (inList ? platformInfo.icon : '/img/socialhandles/default.svg')

  const style = `background-image: url(${finalIcon});`

  const itemClass = classNames({
    "user-social-item": true,
    "flex": true,
    "justify-start": true,
    "items-center": true,
    "-mr-1": true,
    "use-default": !inList,
  });
  const labelClass = classNames({
    "platform-label": true,
    "visually-hidden": !!finalIcon,
    "text-indigo-800": true,
    "text-xs": true,
    "uppercase": true,
    "max-w-xs": true,
    "inline-block": true,
  });

  return (
    html`
      <a class=${itemClass} target="_blank" href=${url}>
        <span class="platform-icon rounded-lg bg-no-repeat"
              style=${style} title="Find me on ${name}"></span>
        <span class=${labelClass}>Find me on ${name}</span>
      </a>
    `);
}

/*
function SocialIcon(props) {
  const { platform, url } = props;
  const platformInfo = SOCIAL_PLATFORMS[platform.toLowerCase()];
  const inList = !!platformInfo;
  const imgRow = inList ? platformInfo.imgPos[0] : 0;
  const imgCol = inList ? platformInfo.imgPos[1] : 0;

  const name = inList ? platformInfo.name : platform;

  const style = `--imgRow: -${imgRow}; --imgCol: -${imgCol};`;
  const itemClass = classNames({
    "user-social-item": true,
    "flex": true,
    "justify-start": true,
    "items-center": true,
    "-mr-1": true,
    "use-default": !inList,
  });
  const labelClass = classNames({
    "platform-label": true,
    "visually-hidden": inList,
    "text-indigo-800": true,
    "text-xs": true,
    "uppercase": true,
    "max-w-xs": true,
    "inline-block": true,
  });

  return (
    html`
      <a class=${itemClass} target="_blank" href=${url}>
        <span class="platform-icon rounded-lg bg-no-repeat" style=${style}></span>
        <span class=${labelClass}>Find me on ${name}</span>
      </a>
  `);
}*/

export default function (props) {
  const {handles} = props;
  if (handles == null) {
    return null;
  }

  const list = handles.map((item, index) => html`
    <li key="social${index}">
      <${SocialIcon} platform=${item.platform} icon=${item.icon}
                     url=${item.url} color=${item.color}/>
    </li>
  `);

  return html`
    <ul
      id="social-list"
      class="social-list flex flex-row items-center justify-start flex-wrap">
      <span
        class="follow-label text-xs font-bold mr-2 uppercase">Follow me:</span>
      ${list}
    </ul>`;
}
