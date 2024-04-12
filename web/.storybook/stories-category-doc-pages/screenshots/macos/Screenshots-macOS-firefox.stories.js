export default {
  title: 'owncast/Screenshots/macOS/Firefox',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-firefox-desktop-default-offline.png"
      alt="macOS Chrome offline"
      width="100%"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-firefox-desktop-default-online.png"
      alt="macOS Chrome online"
      width="100%"
    />
  ),

  name: 'Online',
};
