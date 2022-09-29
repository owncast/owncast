import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { IconLink } from './IconLink';

export default {
  title: 'atoms/Icon link',
  component: IconLink,
} as ComponentMeta<typeof IconLink>;

const Template: ComponentStory<typeof IconLink> = args => <IconLink {...args} />;

export const GithHubLink = Template.bind({});
GithHubLink.args = {
  title: 'github',
  alt: 'github',
  href: 'https://github.com/owncast/owncast',
  icon: 'https://watch.owncast.online/img/platformlogos/github.svg',
};
