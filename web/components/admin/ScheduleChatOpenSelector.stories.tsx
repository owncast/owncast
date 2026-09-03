import type { Meta, StoryFn } from '@storybook/nextjs';
import { ScheduleChatOpenSelector } from './ScheduleChatOpenSelector';

const meta = {
  title: 'owncast/Admin/Schedule chat open selector',
  component: ScheduleChatOpenSelector,
} satisfies Meta<typeof ScheduleChatOpenSelector>;

export default meta;

const Template: StoryFn<typeof ScheduleChatOpenSelector> = args => (
  <ScheduleChatOpenSelector {...args} />
);

export const Example = {
  render: Template,
  args: {
    value: 10,
    label: 'Open chat before stream',
    tip: 'Let viewers join chat before a scheduled stream begins.',
    options: [
      { value: 0, label: 'Disabled' },
      { value: 5, label: '5 minutes' },
      { value: 10, label: '10 minutes' },
      { value: 30, label: '30 minutes' },
      { value: 60, label: '60 minutes' },
    ],
    onChange: () => {},
  },
};
