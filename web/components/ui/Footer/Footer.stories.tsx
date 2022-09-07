import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Footer } from './Footer';

export default {
  title: 'owncast/Layout/Footer',
  component: Footer,
  parameters: {},
} as ComponentMeta<typeof Footer>;

const Template: ComponentStory<typeof Footer> = args => <Footer {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example = Template.bind({});
Example.args = {
  version: 'v1.2.3',
};
