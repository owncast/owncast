import { StoryFn, Meta } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { ActionButton } from './ActionButton';

const meta = {
  title: 'owncast/Components/Action Buttons/Single button',
  component: ActionButton,
  parameters: {
    docs: {
      description: {
        component: `An **Action Button** or **External Action Button** is a button that is used to trigger either an internal or external action. Many will show a modal, but they can also open a new tab to allow navigating to external pages. They are rendered horizontally within the Action Button Row.`,
      },
    },
  },
} satisfies Meta<typeof ActionButton>;

export default meta;

const itemSelected = a => {
  console.log('itemSelected', a);
  action(a.title);
};

const Template: StoryFn<typeof ActionButton> = args => (
  <ActionButton externalActionSelected={itemSelected} {...args} />
);

export const Example1 = {
  render: Template,

  args: {
    action: {
      url: 'https://owncast.online/docs',
      title: 'Documentation',
      description: 'Owncast Documentation',
      icon: 'https://owncast.online/images/logo.svg',
      color: '#5232c8',
      openExternally: false,
    },
  },
};

export const Example2 = {
  render: Template,

  args: {
    action: {
      url: 'https://opencollective.com/embed/owncast/donate',
      title: 'Support Owncast',
      description: 'Contribute to Owncast',
      icon: 'https://opencollective.com/static/images/opencollective-icon.svg',
      color: '#2b4863',
      openExternally: false,
    },
  },
};
