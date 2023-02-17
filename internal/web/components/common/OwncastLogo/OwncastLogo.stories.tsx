import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { OwncastLogo } from './OwncastLogo';

export default {
  title: 'owncast/Components/Header Logo',
  component: OwncastLogo,
  parameters: {
    chromatic: { diffThreshold: 0.8 },
  },
} as ComponentMeta<typeof OwncastLogo>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof OwncastLogo> = args => <OwncastLogo {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Logo = Template.bind({});
Logo.args = {
  url: '/logo',
};

export const DemoServer = Template.bind({});
DemoServer.args = {
  url: 'https://watch.owncast.online/logo',
};
