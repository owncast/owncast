import { Meta } from '@storybook/react';
import { ContentHeader } from './ContentHeader';

const meta = {
  title: 'owncast/Components/Content Header',
  component: ContentHeader,
  parameters: {},
} satisfies Meta<typeof ContentHeader>;

export default meta;

export const Example = {
  args: {
    name: 'My Awesome Owncast Stream',
    summary: 'A calvacade of glorious sights and sounds',
    tags: ['word', 'tag with spaces', 'music'],
    logo: 'https://watch.owncast.online/logo',
    links: [
      {
        platform: 'github',
        url: 'https://github.com/owncast/owncast',
        icon: 'https://watch.owncast.online/img/platformlogos/github.svg',
      },
      {
        platform: 'Documentation',
        url: 'https://owncast.online',
        icon: 'https://watch.owncast.online/img/platformlogos/link.svg',
      },
      {
        platform: 'mastodon',
        url: 'https://fosstodon.org/users/owncast',
        icon: 'https://watch.owncast.online/img/platformlogos/mastodon.svg',
      },
    ],
  },
};

export const LongContent = {
  args: {
    name: 'My Awesome Owncast Stream, streaming the best of streams and some lorem ipsum too',
    summary:
      'A calvacade of glorious sights and sounds. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
    tags: [
      'word',
      'tag with spaces',
      'music',
      'more tags',
      'a bunch',
      'keep going',
      'and more',
      'just a few more',
      'video games',
      'things',
      'stuff',
      'ok some more',
      'this should do it',
    ],
    logo: 'https://watch.owncast.online/logo',
    links: [
      {
        platform: 'github',
        url: 'https://github.com/owncast/owncast',
        icon: 'https://watch.owncast.online/img/platformlogos/github.svg',
      },
      {
        platform: 'Documentation',
        url: 'https://owncast.online',
        icon: 'https://watch.owncast.online/img/platformlogos/link.svg',
      },
      {
        platform: 'mastodon',
        url: 'https://fosstodon.org/users/owncast',
        icon: 'https://watch.owncast.online/img/platformlogos/mastodon.svg',
      },
    ],
  },
};
