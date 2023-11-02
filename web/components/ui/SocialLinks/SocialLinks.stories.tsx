import { Meta } from '@storybook/react';
import { SocialLinks } from './SocialLinks';

const meta = {
  title: 'owncast/Components/Social links',
  component: SocialLinks,
  parameters: {},
} satisfies Meta<typeof SocialLinks>;

export default meta;

export const Populated = {
  args: {
    links: [
      {
        platform: 'github',
        url: 'https://github.com/owncast/owncast',
        icon: '/img/platformlogos/github.svg',
      },
      {
        platform: 'Documentation',
        url: 'https://owncast.online',
        icon: '/img/platformlogos/link.svg',
      },
      {
        platform: 'mastodon',
        url: 'https://fosstodon.org/users/owncast',
        icon: '/img/platformlogos/mastodon.svg',
      },
    ],
  },
};

export const Empty = {
  args: {
    links: [],
  },
};
