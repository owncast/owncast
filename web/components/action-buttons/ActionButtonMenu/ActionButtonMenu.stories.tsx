import { StoryFn, Meta } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { ActionButtonMenu } from './ActionButtonMenu';

const meta = {
  title: 'owncast/Components/Action Buttons/Action Menu',
  component: ActionButtonMenu,
  parameters: {},
} satisfies Meta<typeof ActionButtonMenu>;

export default meta;

const itemSelected = a => {
  console.log('itemSelected', a);
  action(a.title);
};

const Template: StoryFn<typeof ActionButtonMenu> = args => (
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

export const Example = {
  render: Template,

  args: {
    actions,
  },
};

export const ShowFollowExample = {
  render: Template,

  args: {
    actions,
    showFollowItem: true,
  },
};

export const ShowNotifyExample = {
  render: Template,

  args: {
    actions,
    showNotifyItem: true,
  },
};

export const ShowNotifyAndFollowExample = {
  render: Template,

  args: {
    actions,
    showNotifyItem: true,
    showFollowItem: true,
  },
};
