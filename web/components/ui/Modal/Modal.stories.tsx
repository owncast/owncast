import { StoryFn, Meta } from '@storybook/react';
import { Modal } from './Modal';

const meta = {
  title: 'owncast/Modals/Container',
  component: Modal,
  parameters: {
    docs: {
      description: {
        component: `This is the popup modal container that all modal content is rendered inside. It can be passed content nodes to render, or a URL to show an iframe.`,
      },
    },
  },
} satisfies Meta<typeof Modal>;

export default meta;

const Template: StoryFn<typeof Modal> = args => {
  const { children } = args;
  return <Modal {...args}>{children}</Modal>;
};

export const Example = {
  render: Template,

  args: {
    title: 'Modal example with content nodes',
    visible: true,
    children: <div>Test 123</div>,
  },
};

export const UrlExample = {
  render: Template,

  args: {
    title: 'Modal example with URL',
    visible: true,
    url: 'https://owncast.online',
  },
};
