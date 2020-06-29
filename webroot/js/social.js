const SOCIAL_PLATFORMS_URLS = {
  default: {
    name: "default",
    urlPrefix: "",
    imgPos: [0,0], // [row,col]
  },

  facebook: {
    name: "Facebook",
    urlPrefix: "http://www.facebook.com/",
    imgPos: [0,1],
  },
  twitter: {
    name: "Twitter",
    urlPrefix: "http://www.twitter.com/",
    imgPos: [0,2],
  },
  instagram: {
    name: "Instagram",
    urlPrefix: "http://www.instagram.com/",
    imgPos: [0,3],
  },
  snapchat: {
    name: "Snapchat",
    urlPrefix: "http://www.snapchat.com/",
    imgPos: [0,4],
  },
  tiktok: {
    name: "TikTok",
    urlPrefix: "http://www.tiktok.com/",
    imgPos: [0,5],
  },
  soundcloud: {
    name: "Soundcloud",
    urlPrefix: "http://www.soundcloud.com/",
    imgPos: [0,6],
  },
  bandcamp: {
    name: "Bandcamp",
    urlPrefix: "http://www.basecamp.com/",
    imgPos: [0,7],
  },
  patreon: {
    name: "Patreon",
    urlPrefix: "http://www.patreon.com/",
    imgPos: [0,1],
  },
  youtube: {
    name: "YouTube",
    urlPrefix: "http://www.youtube.com/",
    imgPos: [0,9 ],
  },
  spotify: {
    name: "Spotify",
    urlPrefix: "http://www.spotify.com/",
    imgPos: [0,10],
  },
  twitch: {
    name: "Twitch",
    urlPrefix: "http://www.twitch.com/",
    imgPos: [0,11],
  },
  paypal: {
    name: "Paypal",
    urlPrefix: "http://www.paypal.com/",
    imgPos: [0,12],
  },
  github: {
    name: "Github",
    urlPrefix: "http://www.github.com/",
    imgPos: [0,13],
  },
  linkedin: {
    name: "LinkedIn",
    urlPrefix: "http://www.linkedin.com/",
    imgPos: [0,14],
  },
  discord: {
    name: "Discord",
    urlPrefix: "http://www.discord.com/",
    imgPos: [0,15],
  },
  mastadon: {
    name: "Mastadon",
    urlPrefix: "http://www.mastadon.com/",
    imgPos: [0,16],
  },
};

Vue.component('social-list', {
  props: ['platforms'],

  template: `
    <ul class="social-list flex">
      <span v-if="this.platforms.length" class="follow-label">Follow me: </span>
      <user-social-icon 
        v-for="(item, index) in this.platforms"
        v-if="item.platform && item.handle"
        v-bind:key="index"
        v-bind:platform="item.platform"
        v-bind:username="item.handle"
      />
    </ul>
  `,

});

Vue.component('user-social-icon', {
  props: ['platform', 'username'],
  data: function() {
    const platformInfo = SOCIAL_PLATFORMS_URLS[this.platform.toLowerCase()] || SOCIAL_PLATFORMS_URLS["default"];
    const imgRow = platformInfo.imgPos && platformInfo.imgPos[0] || 0;
    const imgCol = platformInfo.imgPos && platformInfo.imgPos[1] || 0;
    const useDefault = platformInfo.name === "default";
    return {
      name: platformInfo.name,
      link: platformInfo.name !== "default" ? `${platformInfo.urlPrefix}/${this.username}` : '#',
      style: `--imgRow: -${imgRow}; --imgCol: -${imgCol};`,
      itemClass: {
        "user-social-item": true,
        "flex": true,
        "rounded": useDefault,
        "use-default": useDefault,        
      },
      labelClass: {
        "platform-label": true,
        "visually-hidden": !useDefault,
        "text-indigo-800": true,
      },
    };
  },
  template: `
    <li>
      <a
        v-bind:class="itemClass"
        target="_blank"
        :href="link"
      >
        <span class="platform-icon rounded-lg" :style="style" />
        <span v-bind:class="labelClass">Find @{{username}} on {{platform}}</span>
      </a>
    </li>
  `,
});
