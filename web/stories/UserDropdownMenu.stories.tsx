import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import UserDropdownMenu from '../components/UserDropdownMenu';
import { ChatState } from '../interfaces/application-state';

export default {
  title: 'owncast/User settings menu',
  component: UserDropdownMenu,
  parameters: {},
} as ComponentMeta<typeof UserDropdownMenu>;

// This component uses Recoil internally so wrap it in a RecoilRoot.
const Example = args => (
  <RecoilRoot>
    <UserDropdownMenu {...args} />
  </RecoilRoot>
);

const Template: ComponentStory<typeof UserDropdownMenu> = args => <Example {...args} />;

export const ChatEnabled = Template.bind({});
ChatEnabled.args = {
  username: 'test-user',
  chatState: ChatState.Available,
};

export const ChatDisabled = Template.bind({});
ChatDisabled.args = {
  username: 'test-user',
  chatState: ChatState.NotAvailable,
};
