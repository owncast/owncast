import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Logo } from '../components/ui/Logo/Logo';

export default {
  title: 'owncast/Components/Page Logo',
  component: Logo,
  parameters: {
    chromatic: { diffThreshold: 0.8 },
  },
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

export const NotSquare = Template.bind({});
NotSquare.args = {
  src: 'https://via.placeholder.com/150x325/FF0000/FFFFFF?text=Rectangle',
};
