import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ExternalActionButton from '../components/action-buttons/ExternalActionButton';

export default {
  title: 'owncast/External action button',
  component: ExternalActionButton,
  parameters: {},
} as ComponentMeta<typeof ExternalActionButton>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ExternalActionButton> = args => (
  <ExternalActionButton {...args} />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example1 = Template.bind({});
Example1.args = {
  action: {
    url: 'https://owncast.online/docs',
    title: 'Documentation',
    description: 'Owncast Documentation',
    icon: 'https://owncast.online/images/logo.svg',
    color: '#5232c8',
    openExternally: false,
  },
};

export const Example2 = Template.bind({});
Example2.args = {
  action: {
    url: 'https://opencollective.com/embed/owncast/donate',
    title: 'Support Owncast',
    description: 'Contribute to Owncast',
    icon: 'https://opencollective.com/static/images/opencollective-icon.svg',
    color: '#2b4863',
    openExternally: false,
  },
};
