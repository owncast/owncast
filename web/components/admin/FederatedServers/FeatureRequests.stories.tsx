import React from 'react';
import { Meta, StoryFn } from '@storybook/react';
import { FeatureRequests, FeatureRequestsProps } from './FeatureRequests';

export default {
  title: 'owncast/Admin/Featured Streams/FeatureRequests',
  component: FeatureRequests,
  parameters: {
    docs: {
      description: {
        component:
          'Lists incoming requests from other Owncast servers asking to feature this server’s stream, with approve/reject actions. Renders nothing when there are no pending requests.',
      },
    },
  },
} as Meta<typeof FeatureRequests>;

const noop = async () => {};

const Template: StoryFn<FeatureRequestsProps> = args => <FeatureRequests {...args} />;

export const WithRequests = Template.bind({});
WithRequests.args = {
  onApprove: noop,
  onReject: noop,
  requests: [
    {
      link: 'https://goodnight.example.com/federation/user/goodnight',
      name: 'Goodnight TV',
      username: 'goodnight@goodnight.example.com',
      image: 'https://picsum.photos/seed/20/64/64',
    },
    {
      link: 'https://music.example.com/federation/user/music',
      name: 'Music Studio Live',
      username: 'music@music.example.com',
    },
  ],
};

export const Loading = Template.bind({});
Loading.args = {
  onApprove: noop,
  onReject: noop,
  loading: true,
  requests: [],
};
