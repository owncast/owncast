import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import ChatTextField from '../components/chat/ChatTextField/ChatTextField';

export default {
  title: 'owncast/Chat/Input text field',
  component: ChatTextField,
  parameters: {
    docs: {
      description: {
        component: `
- This is a element using \`contentEditable\` in order to support rendering emoji images inline.
- Emoji button shows emoji picker.
- Should show one line by default, but grow to two lines as needed.
- The Send button should be hidden for desktop layouts and be shown for mobile layouts.`,
      },
    },
  },
} as ComponentMeta<typeof ChatTextField>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatTextField> = args => (
  <RecoilRoot>
    <ChatTextField {...args} />
  </RecoilRoot>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example = Template.bind({});

export const LongerMessage = Template.bind({});
LongerMessage.args = {
  value:
    'Lorem ipsum dolor sit amet, <img src="https://watch.owncast.online/img/emoji/bluntparrot.gif" width="40px" /> consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
};

LongerMessage.parameters = {
  docs: {
    description: {
      story: 'Should display two lines of text and scroll to display more.',
    },
  },
};
