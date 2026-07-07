import React, { useState } from 'react';
import { Meta, StoryFn } from '@storybook/react';
import { Button } from 'antd';
import { FeatureStreamModal, FeatureStreamModalProps } from './FeatureStreamModal';

export default {
  title: 'owncast/Admin/FeatureStreamModal',
  component: FeatureStreamModal,
  parameters: {
    docs: {
      description: {
        component: 'Modal for featuring/adding Owncast streams with URL validation.',
      },
    },
  },
} as Meta<typeof FeatureStreamModal>;

const Template: StoryFn<FeatureStreamModalProps> = args => {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button onClick={() => setOpen(true)}>Open Feature Stream Modal</Button>
      <FeatureStreamModal
        {...args}
        open={open}
        onCancel={() => setOpen(false)}
        onOk={async (url: string) => {
          console.log('Featuring stream:', url);
          await new Promise(resolve => {
            setTimeout(resolve, 1000);
          }); // Simulate API call
          setOpen(false);
        }}
      />
    </>
  );
};

export const Default = Template.bind({});
Default.args = {};

export const WithError = () => {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button onClick={() => setOpen(true)}>Open Modal with Error</Button>
      <FeatureStreamModal
        open={open}
        onCancel={() => setOpen(false)}
        onOk={async () => {
          throw new Error('Failed to connect to stream');
        }}
      />
    </>
  );
};

export const AlwaysOpen: StoryFn<FeatureStreamModalProps> = () => (
  <FeatureStreamModal
    open
    onCancel={() => console.log('Cancel clicked')}
    onOk={async (url: string) => {
      console.log('Featuring stream:', url);
      await new Promise(resolve => {
        setTimeout(resolve, 1000);
      });
    }}
  />
);
