import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { IndieAuthModal } from './IndieAuthModal';
import Mock from '../../../stories/assets/mocks/indieauth-modal.png';

const Example = () => (
  <div>
    <IndieAuthModal authenticated displayName="fakeChatName" accessToken="fakeaccesstoken" />
  </div>
);

export default {
  title: 'owncast/Modals/IndieAuth',
  component: IndieAuthModal,
  parameters: {
    design: {
      type: 'image',
      url: Mock,
      scale: 0.5,
    },
  },
} as ComponentMeta<typeof IndieAuthModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof IndieAuthModal> = () => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
