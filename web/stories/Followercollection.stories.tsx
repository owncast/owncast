import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import FollowerCollection from '../components/FollowersCollection';

export default {
  title: 'owncast/Follower collection',
  component: FollowerCollection,
  parameters: {},
} as ComponentMeta<typeof FollowerCollection>;

const Template: ComponentStory<typeof FollowerCollection> = args => (
  <FollowerCollection {...args} />
);

export const Example = Template.bind({});
Example.args = {
  followers: [
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
    {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
  ],
};
