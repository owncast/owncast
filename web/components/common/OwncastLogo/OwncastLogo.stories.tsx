import { Meta } from '@storybook/react';
import { OwncastLogo } from './OwncastLogo';

const meta = {
  title: 'owncast/Components/Header Logo',
  component: OwncastLogo,
  parameters: {
    chromatic: { diffThreshold: 0.8 },
  },
} satisfies Meta<typeof OwncastLogo>;

export default meta;

export const Logo = {
  args: {
    url: '/logo',
  },
};

export const DemoServer = {
  args: {
    url: 'https://watch.owncast.online/logo',
  },
};
