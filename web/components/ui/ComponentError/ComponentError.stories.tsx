import { Meta } from '@storybook/react';
import { ComponentError } from './ComponentError';

const meta = {
  title: 'owncast/Components/Component Error',
  component: ComponentError,
  parameters: {
    docs: {
      description: {
        component: `This component is used to display a user-facing fatal error within a component's error boundary. It enables a link to file a bug report. It should have enough detail to help the developers fix the issue, but not be so unapproachable it makes the user scared away.`,
      },
    },
  },
} satisfies Meta<typeof ComponentError>;

export default meta;

export const DefaultMessage = {
  args: {
    componentName: 'Test Component',
  },
};

export const Error1 = {
  args: { message: 'This is a test error message.', componentName: 'Test Component' },
};

export const WithDetails = {
  args: {
    message: 'This is a test error message.',
    componentName: 'Test Component',
    details: 'Here are some additional details about the error.',
  },
};

export const CanRetry = {
  args: {
    message: 'This is a test error message.',
    componentName: 'Test Component',
    details: 'Here are some additional details about the error.',
    retryFunction: () => {
      console.log('retrying');
    },
  },
};
