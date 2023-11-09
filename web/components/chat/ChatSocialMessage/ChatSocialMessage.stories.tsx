import { Meta } from '@storybook/react';
import { ChatSocialMessage } from './ChatSocialMessage';

const meta = {
  title: 'owncast/Chat/Messages/Social-fediverse event',
  component: ChatSocialMessage,
  parameters: {},
} satisfies Meta<typeof ChatSocialMessage>;

export default meta;

export const Follow = {
  args: {
    message: {
      type: 'FEDIVERSE_ENGAGEMENT_FOLLOW',
      body: '<p>james followed this live stream.</p>',
      title: 'james@mastodon.social',
      image: 'https://mastodon.social/avatars/original/missing.png',
      link: 'https://mastodon.social/@james',
    },
  },
};

export const Like = {
  args: {
    message: {
      type: 'FEDIVERSE_ENGAGEMENT_LIKE',
      body: '<p>james liked that this stream went live.</p>',
      title: 'james@mastodon.social',
      image: 'https://mastodon.social/avatars/original/missing.png',
      link: 'https://mastodon.social/@james',
    },
  },
};

export const Repost = {
  args: {
    message: {
      type: 'FEDIVERSE_ENGAGEMENT_REPOST',
      body: '<p>james shared this stream with their followers.</p>',
      title: 'james@mastodon.social',
      image: 'https://mastodon.social/avatars/original/missing.png',
      link: 'https://mastodon.social/@james',
    },
  },
};

export const LongAccountName = {
  args: {
    message: {
      type: 'FEDIVERSE_ENGAGEMENT_REPOST',
      body: '<p>james shared this stream with their followers.</p>',
      title: 'littlejimmywilliams@technology.biz.net.org.technology.gov',
      image: 'https://mastodon.social/avatars/original/missing.png',
      link: 'https://mastodon.social/@james',
    },
  },
};

export const InvalidAvatarImage = {
  args: {
    message: {
      type: 'FEDIVERSE_ENGAGEMENT_REPOST',
      body: '<p>james shared this stream with their followers.</p>',
      title: 'james@mastodon.social',
      image: 'https://xx.xx/avatars/original/missing.png',
      link: 'https://mastodon.social/@james',
    },
  },
};
