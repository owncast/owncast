import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { OfflineBanner } from './OfflineBanner';
import OfflineState from '../../../stories/assets/mocks/offline-state.png';

const meta = {
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
} satisfies Meta<typeof OfflineBanner>;

export default meta;

const Template: StoryFn<typeof OfflineBanner> = args => (
  <RecoilRoot>
    <OfflineBanner {...args} />
  </RecoilRoot>
);

export const ExampleDefaultWithNotifications = {
  render: Template,

  args: {
    streamName: 'Cool stream 42',
    notificationsEnabled: true,
    lastLive: new Date(),
  },
};

export const ExampleDefaultWithDateAndFediverse = {
  render: Template,

  args: {
    streamName: 'Dull stream 31337',
    lastLive: new Date(),
    notificationsEnabled: false,
    fediverseAccount: 'streamer@coolstream.biz',
  },
};

export const ExampleCustomWithDateAndNotifications = {
  render: Template,

  args: {
    streamName: 'Dull stream 31337',
    customText:
      'This is some example offline text that a streamer can leave for a visitor of the page.',
    lastLive: new Date(),
    notificationsEnabled: true,
  },
};

export const ExampleDefaultWithNotificationsAndFediverse = {
  render: Template,

  args: {
    streamName: 'Cool stream 42',
    notificationsEnabled: true,
    fediverseAccount: 'streamer@coolstream.biz',
    lastLive: new Date(),
  },
};

export const ExampleDefaultWithoutNotifications = {
  render: Template,

  args: {
    streamName: 'Cool stream 42',
    notificationsEnabled: false,
    lastLive: new Date(),
  },
};

export const ExampleCustomTextWithoutNotifications = {
  render: Template,

  args: {
    streamName: 'Dull stream 31337',
    customText:
      'This is some example offline text that a streamer can leave for a visitor of the page.',
  },
};
