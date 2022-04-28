import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import UserDropdownMenu from '../components/UserDropdownMenu';

export default {
  title: 'owncast/User settings dropdown menu',
  component: UserDropdownMenu,
  parameters: {},
} as ComponentMeta<typeof UserDropdownMenu>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof UserDropdownMenu> = args => <UserDropdownMenu />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example = Template.bind({});
