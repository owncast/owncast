import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import PageLogo from '../components/PageLogo';

export default {
  title: 'owncast/Page Logo',
  component: PageLogo,
  parameters: {},
} as ComponentMeta<typeof PageLogo>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof PageLogo> = args => <PageLogo {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Logo = Template.bind({});
Logo.args = {
  url: '/logo',
};

export const DemoServer = Template.bind({});
DemoServer.args = {
  url: 'https://watch.owncast.online/logo',
};
