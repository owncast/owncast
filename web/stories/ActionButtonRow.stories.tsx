import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import ActionButtonRow from '../components/action-buttons/ActionButtonRow';
import ActionButton from '../components/action-buttons/ActionButton';

export default {
  title: 'owncast/External action button row',
  component: ActionButtonRow,
  parameters: {},
} as ComponentMeta<typeof ActionButtonRow>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ActionButtonRow> = args => (
  <ActionButtonRow>{args.buttons}</ActionButtonRow>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const actions = [
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
];

const buttons = actions.map(action => <ActionButton action={action} />);
export const Example1 = Template.bind({});
Example1.args = {
  buttons,
};
