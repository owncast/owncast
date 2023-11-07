import { StoryFn, Meta } from '@storybook/react';
import { IndieAuthModal } from './IndieAuthModal';
import Mock from '../../../stories/assets/mocks/indieauth-modal.png';

const Example = () => (
  <div>
    <IndieAuthModal authenticated displayName="fakeChatName" accessToken="fakeaccesstoken" />
  </div>
);

const meta = {
  title: 'owncast/Modals/IndieAuth',
  component: IndieAuthModal,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
      scale: 0.5,
    },
  },
} satisfies Meta<typeof IndieAuthModal>;

export default meta;

const Template: StoryFn<typeof IndieAuthModal> = () => <Example />;

export const Basic = {
  render: Template,
};
