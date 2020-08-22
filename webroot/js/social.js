import { html } from "https://unpkg.com/htm/preact/index.mjs?module";
import { SOCIAL_PLATFORMS } from './utils/social.js';
import { classNames } from './utils.js';

export default function SocialIcon(props) {
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
}
