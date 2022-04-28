import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatTextField from '../components/chat/ChatTextField';

const Example = () => (
  <div>
    <ChatTextField />
  </div>
);

export default {
  title: 'owncast/Chat/Input text field',
  component: ChatTextField,
  parameters: {},
} as ComponentMeta<typeof ChatTextField>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatTextField> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
