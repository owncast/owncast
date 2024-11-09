export default {
  title: 'owncast/Screenshots/Android Portrait/Stock Browser',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/android-13.0-android-browser-google-pixel-7-pro-portrait-offline.png"
      alt="Android Browser offline"
      height="1000px"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/android-13.0-android-browser-google-pixel-7-pro-portrait-online.png"
      alt="Android Browser offline"
      height="1000px"
    />
  ),

  name: 'Online',
};
