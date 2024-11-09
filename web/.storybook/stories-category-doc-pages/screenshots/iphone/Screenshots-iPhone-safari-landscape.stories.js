export default {
  title: 'owncast/Screenshots/iPhone/Safari/Landscape',

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
      alt="iPhone Safari offline"
      height="1000px"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/ios-16-mobile-safari-ipad-pro-11-2022-landscape-online.png"
      alt="iPhone Safari online"
      height="1000px"
    />
  ),

  name: 'Online',
};
