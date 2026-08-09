import { Meta, StoryObj } from '@storybook/nextjs';
import { ClipButton } from './ClipButton';

const meta = {
  title: 'owncast/Action buttons/Clip',
  component: ClipButton,
} satisfies Meta<typeof ClipButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Ready: Story = {
  args: {
    onClick: () => {},
  },
};

export const Capturing: Story = {
  args: {
    active: true,
    remainingSeconds: 95,
    onClick: () => {},
  },
};
