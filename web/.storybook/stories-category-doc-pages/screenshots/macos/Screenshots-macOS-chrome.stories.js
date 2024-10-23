export default {
  title: 'owncast/Screenshots/macOS/Chrome',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-chrome-desktop-default-offline.png"
      alt="macOS Chrome offline"
      width="100%"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-chrome-desktop-default-online.png"
      alt="macOS Chrome online"
      width="100%"
    />
  ),

  name: 'Online',
};
