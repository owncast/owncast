import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { ChatSocialMessage } from './ChatSocialMessage';

export default {
  title: 'owncast/Chat/Messages/Social-fediverse event',
  component: ChatSocialMessage,
  parameters: {},
} as ComponentMeta<typeof ChatSocialMessage>;

const Template: ComponentStory<typeof ChatSocialMessage> = args => <ChatSocialMessage {...args} />;

export const Follow = Template.bind({});
Follow.args = {
  message: {
    type: 'FEDIVERSE_ENGAGEMENT_FOLLOW',
    body: 'james followed this live stream.',
    title: 'james@mastodon.social',
    image: 'https://mastodon.social/avatars/original/missing.png',
    link: 'https://mastodon.social/@james',
  },
};

export const Like = Template.bind({});
Like.args = {
  message: {
    type: 'FEDIVERSE_ENGAGEMENT_LIKE',
    body: 'james liked that this stream went live.',
    title: 'james@mastodon.social',
    image: 'https://mastodon.social/avatars/original/missing.png',
    link: 'https://mastodon.social/@james',
  },
};

export const Repost = Template.bind({});
Repost.args = {
  message: {
    type: 'FEDIVERSE_ENGAGEMENT_REPOST',
    body: 'james shared this stream with their followers.',
    title: 'james@mastodon.social',
    image: 'https://mastodon.social/avatars/original/missing.png',
    link: 'https://mastodon.social/@james',
  },
};
