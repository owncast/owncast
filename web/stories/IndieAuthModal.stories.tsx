import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import IndieAuthModal from '../components/modals/IndieAuthModal';
import Mock from './assets/mocks/indieauth-modal.png';

const Example = () => (
  <div>
    <IndieAuthModal />
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
const Template: ComponentStory<typeof IndieAuthModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
