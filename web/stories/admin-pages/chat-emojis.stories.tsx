import { Meta, StoryObj } from '@storybook/nextjs';
import ChatEmojis from '../../pages/admin/chat/emojis';

// Render smoke for the admin Chat Emojis page: a story that throws fails the
// Storybook build, catching page-level breakage on every PR. Renders with
// default (empty) server config context, no backend.
const meta = {
  title: 'owncast/Admin/Pages/Chat Emojis',
  component: ChatEmojis,
} satisfies Meta<typeof ChatEmojis>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
