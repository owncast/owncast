import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import Logo from '../components/ui/Logo/Logo';

export default {
  title: 'owncast/Components/Page Logo',
  component: Logo,
  parameters: {},
} as ComponentMeta<typeof Logo>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof Logo> = args => <Logo {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const LocalServer = Template.bind({});
LocalServer.args = {
  src: 'http://localhost:8080/logo',
};

export const DemoServer = Template.bind({});
DemoServer.args = {
  src: 'https://watch.owncast.online/logo',
};

export const RandomImage = Template.bind({});
RandomImage.args = {
  src: 'https://picsum.photos/600/500',
};
