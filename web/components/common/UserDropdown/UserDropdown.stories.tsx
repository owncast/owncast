import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { UserDropdown } from './UserDropdown';

export default {
  title: 'owncast/Components/User settings menu',
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
};

export const ChatDisabled = Template.bind({});
ChatDisabled.args = {
  username: 'test-user',
};
