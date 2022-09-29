import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Modal } from './Modal';

export default {
  title: 'molecules/Modal',
  component: Modal,
} as ComponentMeta<typeof Modal>;

const Template: ComponentStory<typeof Modal> = args => {
  const { children } = args;
  return <Modal {...args}>{children}</Modal>;
};

export const Example = Template.bind({});
Example.args = {
  title: 'Modal example with content nodes',
  visible: true,
  children: <div>Test 123</div>,
};

export const EmbeddedExternalContentExample = Template.bind({});
EmbeddedExternalContentExample.args = {
  title: 'Modal example with URL',
  visible: true,
  url: 'https://owncast.online',
};
