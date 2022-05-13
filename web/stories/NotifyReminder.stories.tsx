import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import NotifyReminder from '../components/ui/NotifyReminderPopup/NotifyReminderPopup';
import Mock from './assets/mocks/notify-popup.png';

export default {
  title: 'owncast/Notify Reminder',
  component: NotifyReminder,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
    },
    docs: {
      description: {
        component: `After visiting the page three times this popup reminding you that you can register for live stream notifications shows up.
Clicking it will make the notificaiton modal display. Clicking the "X" will hide the modal and make it never show again.`,
      },
    },
  },
} as ComponentMeta<typeof NotifyReminder>;

const Template: ComponentStory<typeof NotifyReminder> = args => <NotifyReminder {...args} />;

export const Example = Template.bind({});
Example.args = {};
