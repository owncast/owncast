import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import BrowserNotifyModal from '../components/modals/BrowserNotifyModal';

const Example = () => (
  <div>
    <AuthModal />
  </div>
);

export default {
  title: 'owncast/Modals/Browser Push Notifications',
  component: BrowserNotifyModal,
  parameters: {},
} as ComponentMeta<typeof BrowserNotifyModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof BrowserNotifyModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
