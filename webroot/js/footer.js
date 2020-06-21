Vue.component('owncast-footer', {
  template: `
    <footer class="flex border-t border-gray-500 border-solid">
      <span>
        <a href="https://github.com/gabek/owncast" target="_blank">About Owncast</a>
      </span>          
    </footer>
  `,
});


const SOCIAL_PLATFORMS_URLS = {
  default: {
    name: "default",
    urlPrefix: "",
    imgPos: [0,0], // [row,col]
  },

  facebook: {
    name: "Facebook",
    urlPrefix: "http://www.facebook.com/",
    imgPos: [0,1], // [row,col]
  },
  twitter: {
    name: "Twitter",
    urlPrefix: "http://www.twitter.com/",
    imgPos: [0,2], // [row,col]
  },
  instagram: {
    name: "Instagram",
    urlPrefix: "http://www.instagram.com/",
    imgPos: [0,3], // [row,col]
  },
  instagram: {
    name: "Snapchat",
    urlPrefix: "http://www.snapchat.com/",
    imgPos: [0,4], // [row,col]
  },
  tiktok: {
    name: "TikTok",
    urlPrefix: "http://www.tiktok.com/",
    imgPos: [0,5], // [row,col]
  },
  soundcloud: {
    name: "Soundcloud",
    urlPrefix: "http://www.soundcloud.com/",
    imgPos: [0,6], // [row,col]
  },
  basecamp: {
    name: "Base Camp",
    urlPrefix: "http://www.basecamp.com/",
    imgPos: [0,7], // [row,col]
  },
  patreon: {
    name: "Patreon",
    urlPrefix: "http://www.patreon.com/",
    imgPos: [0,1], // [row,col]
  },
  youtube: {
    name: "YouTube",
    urlPrefix: "http://www.youtube.com/",
    imgPos: [0,9 ], // [row,col]
  },
  spotify: {
    name: "Spotify",
    urlPrefix: "http://www.spotify.com/",
    imgPos: [0,10], // [row,col]
  },
  twitch: {
    name: "Twitch",
    urlPrefix: "http://www.twitch.com/",
    imgPos: [0,11], // [row,col]
  },
  paypal: {
    name: "Paypal",
    urlPrefix: "http://www.paypal.com/",
    imgPos: [0,12], // [row,col]
  },
  github: {
    name: "Github",
    urlPrefix: "http://www.github.com/",
    imgPos: [0,13], // [row,col]
  },
  linkedin: {
    name: "LinkedIn",
    urlPrefix: "http://www.linkedin.com/",
    imgPos: [0,14], // [row,col]
  },
  discord: {
    name: "Discord",
    urlPrefix: "http://www.discord.com/",
    imgPos: [0,15], // [row,col]
  },
  mastadon: {
    name: "Mastadon",
    urlPrefix: "http://www.mastadon.com/",
    imgPos: [0,16], // [row,col]
  },
};

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
    <a
      v-bind:class="itemClass"
      target="_blank"
      :href="link"
    >
      <span 
        class="platform-icon rounded-lg"
        :style="style"
      />
      <span v-bind:class="labelClass">Find @{{username}} on {{platform}}</span>
    </a>
  `,
});
