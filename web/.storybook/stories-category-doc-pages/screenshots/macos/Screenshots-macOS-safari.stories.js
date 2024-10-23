export default {
  title: 'owncast/Screenshots/macOS/Safari',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-safari-desktop-default-offline.png"
      alt="macOS Safari offline"
      width="100%"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/os-x-ventura-safari-desktop-default-online.png"
      alt="macOS Safari online"
      width="100%"
    />
  ),

  name: 'Online',
};
