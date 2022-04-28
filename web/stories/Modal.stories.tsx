import React, { useState } from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Modal, Button } from 'antd';

const Usage = () => (
  <Modal title="Basic Modal" visible>
    <p>Some contents...</p>
    <p>Some contents...</p>
    <p>Some contents...</p>
  </Modal>
);

export default {
  title: 'owncast/Modal container',
  component: Modal,
  parameters: {},
} as ComponentMeta<typeof Modal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof Modal> = args => <Usage />;

export const Example = Template.bind({});
