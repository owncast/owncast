/* eslint-disable no-alert */
import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import NotifyReminder from '../components/ui/NotifyReminderPopup/NotifyReminderPopup';
import Mock from './assets/mocks/notify-popup.png';

const Example = args => (
  <div style={{ margin: '20px', marginTop: '130px' }}>
    <NotifyReminder {...args}>
      <button type="button">notify button</button>
    </NotifyReminder>
  </div>
);

export default {
  title: 'owncast/Components/Notify Reminder',
  component: NotifyReminder,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `After visiting the page three times this popup reminding you that you can register for live stream notifications shows up.
Clicking it will make the notification modal display. Clicking the "X" will hide the modal and make it never show again.`,
      },
    },
  },
} as ComponentMeta<typeof NotifyReminder>;

const Template: ComponentStory<typeof NotifyReminder> = args => <Example {...args} />;

export const Active = Template.bind({});
Active.args = {
  visible: true,
  notificationClicked: () => {
    alert('notification clicked');
  },
  notificationClosed: () => {
    alert('notification closed');
  },
};

export const InActive = Template.bind({});
InActive.args = {
  visible: false,
};
