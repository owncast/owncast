import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { Header } from './Header';

export default {
  title: 'owncast/Layout/Header',
  component: Header,
  parameters: {
    chromatic: { diffThreshold: 0.75 },
  },
} as ComponentMeta<typeof Header>;

const Template: ComponentStory<typeof Header> = args => (
  <RecoilRoot>
    <Header {...args} />
  </RecoilRoot>
);

export const ChatAvailable = Template.bind({});
ChatAvailable.args = {
  name: 'Example Stream Name',
  chatAvailable: true,
};

export const ChatNotAvailable = Template.bind({});
ChatNotAvailable.args = {
  name: 'Example Stream Name',
  chatAvailable: false,
};
