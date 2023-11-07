import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { ChatModerationDetailsModal } from './ChatModerationDetailsModal';

const mocks = {
  mocks: [
    {
      // The "matcher" determines if this
      // mock should respond to the current
      // call to fetch().
      matcher: {
        name: 'response',
        url: 'glob:/api/moderation/chat/user/*',
      },
      // If the "matcher" matches the current
      // fetch() call, the fetch response is
      // built using this "response".
      response: {
        status: 200,
        body: {
          user: {
            id: 'hjFPU967R',
            displayName: 'focused-snyder',
            displayColor: 2,
            createdAt: '2022-07-12T13:08:31.406505322-07:00',
            previousNames: ['focused-snyder'],
            scopes: ['MODERATOR'],
            isBot: false,
            authenticated: false,
          },
          connectedClients: [
            {
              messageCount: 3,
              userAgent:
                'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36',
              connectedAt: '2022-07-20T16:45:07.796685618-07:00',
              geo: 'N/A',
            },
          ],
          messages: [
            {
              id: 'bQp8UJR4R',
              timestamp: '2022-07-20T16:53:41.938083228-07:00',
              user: null,
              body: 'test message 3',
            },
            {
              id: 'ubK88Jg4R',
              timestamp: '2022-07-20T16:53:39.675531279-07:00',
              user: null,
              body: 'test message 2',
            },
            {
              id: '20v8UJRVR',
              timestamp: '2022-07-20T16:53:37.551084121-07:00',
              user: null,
              body: 'test message 1',
            },
          ],
        },
      },
    },
  ],
};

const meta = {
  title: 'owncast/Chat/Moderation modal',
  component: ChatModerationDetailsModal,
  parameters: {
    fetchMock: mocks,
    docs: {
      description: {
        component: `This should be a modal that gives the moderator more details about the user such as:
- When the user was created
- Other names they've used
- If they're authenticated, and using what method (IndieAuth, FediAuth)
`,
      },
    },
  },
} satisfies Meta<typeof ChatModerationDetailsModal>;

export default meta;

const Template: StoryFn<typeof ChatModerationDetailsModal> = () => (
  <RecoilRoot>
    <ChatModerationDetailsModal userId="testuser123" accessToken="fakeaccesstoken4839" />
  </RecoilRoot>
);

export const Example = {
  render: Template,
};
