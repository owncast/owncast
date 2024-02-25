import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { OfflineEmbed } from './OfflineEmbed';
import OfflineState from '../../../stories/assets/mocks/offline-state.png';

const meta = {
  title: 'owncast/Layout/Offline Embed',
  component: OfflineEmbed,
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
} satisfies Meta<typeof OfflineEmbed>;

export default meta;

const Template: StoryFn<typeof OfflineEmbed> = args => (
  <RecoilRoot>
    <OfflineEmbed {...args} />
  </RecoilRoot>
);

export const ExampleDefaultWithNotifications = {
  render: Template,

  args: {
    streamName: 'Cool stream 42',
    subtitle: 'This stream rocks. You should watch it.',
    image: 'https://placehold.co/600x400/orange/white',
  },
};
