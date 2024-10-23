export default {
  title: 'owncast/Screenshots/Android Landscape/Stock Browser',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/android-10.0-android-browser-samsung-galaxy-s20-ultra-landscape-offline.png"
      alt="Android Browser offline"
      height="1000px"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/android-10.0-android-browser-samsung-galaxy-s20-ultra-landscape-online.png"
      alt="Android Browser offline"
      height="1000px"
    />
  ),

  name: 'Online',
};
