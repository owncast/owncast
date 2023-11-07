import { Meta } from '@storybook/react';
import { Logo } from '../components/ui/Logo/Logo';

const meta = {
  title: 'owncast/Components/Page Logo',
  component: Logo,
  parameters: {
    chromatic: { diffThreshold: 0.8 },
  },
} satisfies Meta<typeof Logo>;

export default meta;

export const LocalServer = {
  args: {
    src: 'http://localhost:8080/logo',
  },
};

export const DemoServer = {
  args: {
    src: 'https://watch.owncast.online/logo',
  },
};

export const NotSquare = {
  args: {
    src: 'https://via.placeholder.com/150x325/FF0000/FFFFFF?text=Rectangle',
  },
};
