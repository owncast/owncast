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
  mastadon: {
    name: "Mastadon",
    imgPos: [0,16],
  },
};

Vue.component('social-list', {
  props: ['platforms'],

  template: `
    <ul class="social-list flex" v-if="this.platforms.length">
      <span class="follow-label">Follow me: </span>
      <user-social-icon 
        v-for="(item, index) in this.platforms"
        v-if="item.platform && item.url"
        v-bind:key="index"
        v-bind:platform="item.platform"
        v-bind:url="item.url"
      />
    </ul>
  `,

});

Vue.component('user-social-icon', {
  props: ['platform', 'url'],
  data: function() {
    const platformInfo = SOCIAL_PLATFORMS[this.platform.toLowerCase()];
    const inList = !!platformInfo;
    const imgRow = inList ? platformInfo.imgPos[0] : 0;
    const imgCol = inList ? platformInfo.imgPos[1] : 0;
    return {
      name: inList ? platformInfo.name : this.platform,
      link: this.url,

      style: `--imgRow: -${imgRow}; --imgCol: -${imgCol};`,
      itemClass: {
        "user-social-item": true,
        "flex": true,
        "use-default": !inList,        
      },
      labelClass: {
        "platform-label": true,
        "visually-hidden": inList,
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
        <span v-bind:class="labelClass">Find me on {{platform}}</span>
      </a>
    </li>
  `,
});
