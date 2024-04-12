export default {
  title: 'owncast/Screenshots/Windows/Firefox',

  parameters: {
    chromatic: {
      disableSnapshot: true,
    },
  },
};

export const Offline = {
  render: () => (
    <img
      src="screenshots/windows-10-firefox-desktop-default-offline.png"
      alt="Windows Firefox offline"
      width="100%"
    />
  ),

  name: 'Offline',
};

export const Online = {
  render: () => (
    <img
      src="screenshots/windows-10-firefox-desktop-default-online.png"
      alt="Windows Firefox online"
      width="100%"
    />
  ),

  name: 'Online',
};
