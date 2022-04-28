import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import Header from '../components/ui/Header/Header';

export default {
  title: 'owncast/Header',
  component: Header,
  parameters: {},
} as ComponentMeta<typeof Header>;

const Template: ComponentStory<typeof Header> = args => <Header {...args} />;

export const Example = Template.bind({});
Example.args = {
  name: 'Example Stream Name',
};
