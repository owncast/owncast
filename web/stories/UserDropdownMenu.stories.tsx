import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { UserDropdown } from '../components/common';
import { ChatState } from '../interfaces/application-state';

export default {
  title: 'owncast/User settings menu',
  component: UserDropdown,
  parameters: {},
} as ComponentMeta<typeof UserDropdown>;

// This component uses Recoil internally so wrap it in a RecoilRoot.
const Example = args => (
  <RecoilRoot>
    <UserDropdown {...args} />
  </RecoilRoot>
);

const Template: ComponentStory<typeof UserDropdown> = args => <Example {...args} />;

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
