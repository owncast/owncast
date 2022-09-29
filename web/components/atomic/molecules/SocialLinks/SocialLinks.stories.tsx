import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { SocialLinks } from './SocialLinks';

export default {
  title: 'molecules/Social links',
  component: SocialLinks,
  parameters: {},
} as ComponentMeta<typeof SocialLinks>;

const Template: ComponentStory<typeof SocialLinks> = args => <SocialLinks {...args} />;

export const Populated = Template.bind({});
Populated.args = {
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
};

export const Empty = Template.bind({});
Empty.args = {
  links: [],
};
