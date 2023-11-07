import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { ChatTextField } from './ChatTextField';
import Mockup from '../../../stories/assets/mocks/chatinput-mock.png';

const mockResponse = JSON.parse(
  `[{"name":"Reaper-gg.png","url":"https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=OC"},{"name":"Reaper-hi.png","url":"https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=XX"},{"name":"Reaper-hype.png","url":"https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=TX"},{"name":"Reaper-lol.png","url":"https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=CA"},{"name":"Reaper-love.png","url":"https://ui-avatars.com/api/?background=0D8ABC&color=fff&name=OK"}]`,
);
const mocks = {
  mocks: [
    {
      // The "matcher" determines if this
      // mock should respond to the current
      // call to fetch().
      matcher: {
        name: 'response',
        url: 'glob:/api/emoji',
      },
      // If the "matcher" matches the current
      // fetch() call, the fetch response is
      // built using this "response".
      response: {
        status: 200,
        body: mockResponse,
      },
    },
  ],
};

const meta = {
  title: 'owncast/Chat/Input text field',
  component: ChatTextField,
  parameters: {
    fetchMock: mocks,
    chromatic: { diffThreshold: 0.9 },

    design: {
      type: 'image',
      url: Mockup,
    },
    docs: {
      description: {
        component: `
- This is a element using \`contentEditable\` in order to support rendering emoji images inline.
- Emoji button shows emoji picker.
- Should show one line by default, but grow to two lines as needed.`,
      },
    },
  },
} satisfies Meta<typeof ChatTextField>;

export default meta;

const Template: StoryFn<typeof ChatTextField> = args => (
  <RecoilRoot>
    <ChatTextField {...args} />
  </RecoilRoot>
);

export const Example = {
  render: Template,

  args: {
    enabled: true,
  },
};

export const LongerMessage = {
  render: Template,

  args: {
    enabled: true,
    defaultText:
      'Lorem ipsum dolor sit amet,  consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
  },

  parameters: {
    docs: {
      description: {
        story: 'Should display two lines of text and scroll to display more.',
      },
    },
  },
};

export const DisabledChat = {
  render: Template,

  args: {
    enabled: false,
  },

  parameters: {
    docs: {
      description: {
        story: 'Should not allow you to type anything and should state that chat is disabled.',
      },
    },
  },
};
