import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ExternalActionButtonRow from '../components/action-buttons/ExternalActionButtonRow';

export default {
  title: 'owncast/External action button row',
  component: ExternalActionButtonRow,
  parameters: {},
} as ComponentMeta<typeof ExternalActionButtonRow>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ExternalActionButtonRow> = args => (
  <ExternalActionButtonRow {...args} />
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Example1 = Template.bind({});
Example1.args = {
  actions: [
    {
      url: 'https://owncast.online/docs',
      title: 'Documentation',
      description: 'Owncast Documentation',
      icon: 'https://owncast.online/images/logo.svg',
      color: '#5232c8',
      openExternally: false,
    },
    {
      url: 'https://opencollective.com/embed/owncast/donate',
      title: 'Support Owncast',
      description: 'Contribute to Owncast',
      icon: 'https://opencollective.com/static/images/opencollective-icon.svg',
      color: '#2b4863',
      openExternally: false,
    },
  ],
};
