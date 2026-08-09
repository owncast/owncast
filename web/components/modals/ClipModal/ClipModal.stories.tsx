import { Meta, StoryObj } from '@storybook/nextjs';
import { ClipModal } from './ClipModal';

const meta = {
  title: 'owncast/Modals/Clip',
  component: ClipModal,
} satisfies Meta<typeof ClipModal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const TitlePrompt: Story = {
  args: {
    open: true,
    completedClipId: '',
    handleClose: () => {},
    onSaveTitle: () => {},
    onDiscard: () => {},
    saving: false,
    error: '',
  },
};

export const ShareLink: Story = {
  args: {
    ...TitlePrompt.args,
    completedClipId: 'example-clip',
  },
};
