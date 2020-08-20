import { html } from "https://unpkg.com/htm/preact/index.mjs?module";

const SOCIAL_PLATFORMS = {
  default: {
    name: "default",
    imgPos: [0,0], // [row,col]
  },

  facebook: {
    name: "Facebook",
    imgPos: [0,1],
  },
  twitter: {
    name: "Twitter",
    imgPos: [0,2],
  },
  instagram: {
    name: "Instagram",
    imgPos: [0,3],
  },
  snapchat: {
    name: "Snapchat",
    imgPos: [0,4],
  },
  tiktok: {
    name: "TikTok",
    imgPos: [0,5],
  },
  soundcloud: {
    name: "Soundcloud",
    imgPos: [0,6],
  },
  bandcamp: {
    name: "Bandcamp",
    imgPos: [0,7],
  },
  patreon: {
    name: "Patreon",
    imgPos: [0,1],
  },
  youtube: {
    name: "YouTube",
    imgPos: [0,9 ],
  },
  spotify: {
    name: "Spotify",
    imgPos: [0,10],
  },
  twitch: {
    name: "Twitch",
    imgPos: [0,11],
  },
  paypal: {
    name: "Paypal",
    imgPos: [0,12],
  },
  github: {
    name: "Github",
    imgPos: [0,13],
  },
  linkedin: {
    name: "LinkedIn",
    imgPos: [0,14],
  },
  discord: {
    name: "Discord",
    imgPos: [0,15],
  },
  mastodon: {
    name: "Mastodon",
    imgPos: [0,16],
  },
};

export default function SocialIcon(props) {
  const { platform, url } = props;
  const platformInfo = SOCIAL_PLATFORMS[platform.toLowerCase()];
  const inList = !!platformInfo;
  const imgRow = inList ? platformInfo.imgPos[0] : 0;
  const imgCol = inList ? platformInfo.imgPos[1] : 0;

  const name = inList ? platformInfo.name : platform;

  const style = `--imgRow: -${imgRow}; --imgCol: -${imgCol};`;
  const itemClass = {
    "user-social-item": true,
    "flex": true,
    "use-default": !inList,
  };
  const labelClass = {
    "platform-label": true,
    "visually-hidden": inList,
    "text-indigo-800": true,
  };
  return (
    html`
      <a class=${itemClass} target="_blank" href=${url}>
        <span class="platform-icon rounded-lg" style=${style}></span>
        <span class=${labelClass}>Find me on ${name}</span>
      </a>
  `);
}
