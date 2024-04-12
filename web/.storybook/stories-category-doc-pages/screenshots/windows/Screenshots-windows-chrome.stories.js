export default {
  title: 'owncast/Screenshots/Windows/Chrome',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/windows-11-chrome-desktop-default-offline.png"
      alt="Windows Chrome offline"
      width="100%"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/windows-11-chrome-desktop-default-online.png"
      alt="Windows Chrome online"
      width="100%"
    />
  ),

  name: 'Online',
};
