import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ChatContainer from '../components/chat/ChatContainer';

const Example = () => (
  <div>
    <ChatContainer />
  </div>
);

export default {
  title: 'owncast/Chat/Chat messages container',
  component: ChatContainer,
  parameters: {},
} as ComponentMeta<typeof ChatContainer>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ChatContainer> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
