import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { OfflineBanner } from './OfflineBanner';
import OfflineState from '../../../stories/assets/mocks/offline-state.png';

export default {
  title: 'owncast/Layout/Offline Banner',
  component: OfflineBanner,
  parameters: {
    design: {
      type: 'image',
      url: OfflineState,
      scale: 0.5,
    },
    docs: {
      description: {
        component: `When the stream is offline the player should be replaced by this banner that can support custom text and notify actions.`,
      },
    },
  },
} as ComponentMeta<typeof OfflineBanner>;

const Template: ComponentStory<typeof OfflineBanner> = args => <OfflineBanner {...args} />;

export const ExampleDefault = Template.bind({});
ExampleDefault.args = {
  name: 'Cool stream 42',
  text: 'To get notifications when <server name> is back online you can follow or ask for notifications.',
};

export const ExampleCustom = Template.bind({});
ExampleCustom.args = {
  name: 'Dull stream 31337',
  text: 'This is some example offline text that a streamer can leave for a visitor of the page.',
};
