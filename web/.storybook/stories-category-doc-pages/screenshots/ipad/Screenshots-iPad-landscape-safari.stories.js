export default {
  title: 'owncast/Screenshots/iPad Landscape/Safari',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/ios-16-mobile-safari-ipad-pro-11-2022-landscape-offline.png"
      alt="macOS Safari offline"
      width="800px"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/ios-16-mobile-safari-ipad-pro-11-2022-landscape-online.png"
      alt="macOS Safari online"
      width="800px"
    />
  ),

  name: 'Online',
};
