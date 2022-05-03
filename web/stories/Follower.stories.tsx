import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import SingleFollower from '../components/Follower';

export default {
  title: 'owncast/Follower',
  component: SingleFollower,
  parameters: {},
} as ComponentMeta<typeof SingleFollower>;

const Template: ComponentStory<typeof SingleFollower> = args => <SingleFollower {...args} />;

export const Example = Template.bind({});
Example.args = {
  follower: {
    name: 'John Doe',
    description: 'User',
    username: '@account@domain.tld',
    image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
    link: 'https://yahoo.com',
  },
};
