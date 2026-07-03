import { Meta, StoryObj } from '@storybook/nextjs';
import ActionButtons from './ActionButtons';

// The viewer's action button row (external actions, follow, notify) is
// config-gated inside the Main layout stories, so this story keeps it
// unconditionally visible.
const noop = () => {};

const meta = {
  title: 'owncast/Components/Action buttons row',
  component: ActionButtons,
} satisfies Meta<typeof ActionButtons>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    supportFediverseFeatures: true,
    supportsBrowserNotifications: true,
    showNotifyReminder: false,
    externalActions: [
      {
        title: 'Support the stream',
        description: 'Buy a coffee',
        url: 'https://example.com/support',
        color: '#5232c8',
        openExternally: true,
      },
      {
        title: 'Merch',
        url: 'https://example.com/merch',
      },
    ],
    setShowFollowModal: noop,
    setShowNotifyModal: noop,
    disableNotifyReminderPopup: noop,
    externalActionSelected: noop,
  },
};
