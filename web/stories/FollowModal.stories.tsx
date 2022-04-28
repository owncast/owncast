import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FollowModal from '../components/modals/FollowModal';

const Example = () => (
  <div>
    <FollowModal />
  </div>
);

export default {
  title: 'owncast/Modals/Follow',
  component: FollowModal,
  parameters: {},
} as ComponentMeta<typeof FollowModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof FollowModal> = args => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
