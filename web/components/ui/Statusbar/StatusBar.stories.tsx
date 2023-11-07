import { Meta } from '@storybook/react';
import { subHours } from 'date-fns';
import { Statusbar } from './Statusbar';

const meta = {
  title: 'owncast/Player/Status bar',
  component: Statusbar,
  parameters: {},
} satisfies Meta<typeof Statusbar>;

export default meta;

export const Online = {
  args: {
    online: true,
    viewerCount: 42,
    lastConnectTime: subHours(new Date(), 3),
  },
};

export const Offline = {
  args: {
    online: false,
    lastDisconnectTime: subHours(new Date(), 3),
  },
};
