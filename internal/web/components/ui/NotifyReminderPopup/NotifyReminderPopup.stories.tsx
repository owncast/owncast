/* eslint-disable no-alert */
import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { NotifyReminderPopup } from './NotifyReminderPopup';
import Mock from '../../../stories/assets/mocks/notify-popup.png';

const Example = args => (
  <div style={{ margin: '20px', marginTop: '130px' }}>
    <NotifyReminderPopup {...args}>
      <button type="button">notify button</button>
    </NotifyReminderPopup>
  </div>
);

export default {
  title: 'owncast/Components/Notify Reminder',
  component: NotifyReminderPopup,
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
} as ComponentMeta<typeof NotifyReminderPopup>;

const Template: ComponentStory<typeof NotifyReminderPopup> = args => <Example {...args} />;

export const Active = Template.bind({});
Active.args = {
  open: true,
  notificationClicked: () => {
    alert('notification clicked');
  },
  notificationClosed: () => {
    alert('notification closed');
  },
};

export const InActive = Template.bind({});
InActive.args = {
  open: false,
};
