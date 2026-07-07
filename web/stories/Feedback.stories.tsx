import { Meta, StoryObj } from '@storybook/nextjs';
import { Button, Modal, message, notification } from 'antd';

// Exercises antd's imperative feedback surfaces (message toasts,
// notifications, Modal.confirm) so they have a visual baseline. These render
// outside the normal component tree, which is exactly the plumbing that
// changes across antd major versions, so a story that shows them is the
// cheapest regression net.
const FeedbackDemo = () => (
  <div style={{ height: '400px' }}>
    <Button type="primary" onClick={() => message.success('Settings saved.', 0)}>
      Show a message toast
    </Button>
  </div>
);

const sleep = (ms: number) => {
  const { promise, resolve } = Promise.withResolvers<void>();
  setTimeout(resolve, ms);
  return promise;
};

const meta = {
  title: 'owncast/Feedback',
  component: FeedbackDemo,
  parameters: {
    chromatic: { delay: 1000 },
  },
} satisfies Meta<typeof FeedbackDemo>;

export default meta;
type Story = StoryObj<typeof meta>;

// Fire every imperative feedback type at once. Durations are 0 so nothing
// auto-dismisses before the snapshot is taken.
export const ToastsAndConfirm: Story = {
  play: async () => {
    message.success('Settings saved.', 0);
    message.error('Unable to reach the server.', 0);
    notification.open({
      message: 'New follower',
      description: 'someone@fediverse.example followed your stream.',
      duration: 0,
    });
    Modal.confirm({
      title: 'Delete this stream key?',
      content: 'Anyone using this key will no longer be able to stream.',
      okText: 'Delete',
      okType: 'danger',
    });
    // Let the open animations settle before any snapshot.
    await sleep(500);
  },
};
