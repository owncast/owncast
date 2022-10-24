import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { ActionButtonMenu } from './ActionButtonMenu';

export default {
  title: 'owncast/Components/Action Buttons/Action Menu',
  component: ActionButtonMenu,
  parameters: {},
} as ComponentMeta<typeof ActionButtonMenu>;

const itemSelected = a => {
  console.log('itemSelected', a);
  action(a.title);
};

const Template: ComponentStory<typeof ActionButtonMenu> = args => (
  <ActionButtonMenu {...args} externalActionSelected={a => itemSelected(a)} />
);

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

export const Example = Template.bind({});
Example.args = {
  actions,
};

export const ShowFollowExample = Template.bind({});
ShowFollowExample.args = {
  actions,
  showFollowItem: true,
};

export const ShowNotifyExample = Template.bind({});
ShowNotifyExample.args = {
  actions,
  showNotifyItem: true,
};

export const ShowNotifyAndFollowExample = Template.bind({});
ShowNotifyAndFollowExample.args = {
  actions,
  showNotifyItem: true,
  showFollowItem: true,
};
