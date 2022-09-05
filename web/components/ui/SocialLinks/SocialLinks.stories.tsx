import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { SocialLinks } from './SocialLinks';

export default {
  title: 'owncast/Components/Social links',
  component: SocialLinks,
  parameters: {},
} as ComponentMeta<typeof SocialLinks>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof SocialLinks> = args => <SocialLinks {...args} />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Populated = Template.bind({});
Populated.args = {
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
};

export const Empty = Template.bind({});
Empty.args = {
  links: [],
};
