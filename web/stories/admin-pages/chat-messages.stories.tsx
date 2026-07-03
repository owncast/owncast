import { Meta, StoryObj } from '@storybook/nextjs';
import { http, HttpResponse } from 'msw';
import ChatMessages from '../../pages/admin/chat/messages';

// The chat messages admin page is the only Table usage with row selection,
// which is styling that regressed across antd versions before (selected rows
// were invisible in Firefox with antd 5 defaults). WithSelection keeps a
// visual baseline of the selected-row treatment.

const user = (id: string, displayName: string, displayColor: number) => ({
  id,
  displayName,
  displayColor,
  createdAt: '2026-06-01T12:00:00Z',
  previousNames: [displayName],
  scopes: [],
  isBot: false,
  authenticated: false,
});

const messages = [
  {
    id: 'msg-1',
    timestamp: '2026-07-01T18:00:00Z',
    body: 'Hello everyone, stream starting soon!',
    user: user('user-1', 'StreamerFan42', 1),
    visible: true,
    hiddenAt: null,
  },
  {
    id: 'msg-2',
    timestamp: '2026-07-01T18:01:30Z',
    body: 'The audio is a little quiet on my end.',
    user: user('user-2', 'NightOwl', 3),
    visible: true,
    hiddenAt: null,
  },
  {
    id: 'msg-3',
    timestamp: '2026-07-01T18:02:10Z',
    body: 'This message was hidden by a moderator.',
    user: user('user-3', 'SpamAccount', 5),
    visible: false,
    hiddenAt: '2026-07-01T18:03:00Z',
  },
];

const sleep = (ms: number) => {
  const { promise, resolve } = Promise.withResolvers<void>();
  setTimeout(resolve, ms);
  return promise;
};

const meta = {
  title: 'owncast/Admin/Pages/Chat Messages',
  component: ChatMessages,
  parameters: {
    msw: {
      handlers: [http.get('*/api/admin/chat/messages', () => HttpResponse.json(messages))],
    },
  },
} satisfies Meta<typeof ChatMessages>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

// Select the first row so the selected-row styling is part of the snapshot.
export const WithSelection: Story = {
  play: async ({ canvasElement }) => {
    const deadline = Date.now() + 5000;
    let checkbox: HTMLInputElement | null = null;
    while (!checkbox && Date.now() < deadline) {
      checkbox = canvasElement.querySelector('.ant-table-tbody input[type="checkbox"]');
      if (!checkbox) {
        // eslint-disable-next-line no-await-in-loop
        await sleep(100);
      }
    }
    checkbox?.click();
    await sleep(200);
  },
};
