import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ComponentError } from './ComponentError';

export default {
  title: 'owncast/Components/Component Error',
  component: ComponentError,
  parameters: {
    docs: {
      description: {
        component: `This component is used to display a user-facing fatal error within a component's error boundary. It enables a link to file a bug report. It should have enough detail to help the developers fix the issue, but not be so unapproachable it makes the user scared away.`,
      },
    },
  },
} as ComponentMeta<typeof ComponentError>;

const Template: ComponentStory<typeof ComponentError> = args => <ComponentError {...args} />;

export const DefaultMessage = Template.bind({});
DefaultMessage.args = {
  componentName: 'Test Component',
};

export const Error1 = Template.bind({});
Error1.args = { message: 'This is a test error message.', componentName: 'Test Component' };

export const WithDetails = Template.bind({});
WithDetails.args = {
  message: 'This is a test error message.',
  componentName: 'Test Component',
  details: 'Here are some additional details about the error.',
};

export const CanRetry = Template.bind({});
CanRetry.args = {
  message: 'This is a test error message.',
  componentName: 'Test Component',
  details: 'Here are some additional details about the error.',
  retryFunction: () => {
    console.log('retrying');
  },
};
