import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import * as FollowerComponent from '../components/Follower';

export default {
  title: 'owncast/Follower',
  component: FollowerComponent,
  parameters: {},
} as ComponentMeta<typeof FollowerComponent>;

const Template: ComponentStory<typeof FollowerComponent> = args => <FollowerComponent {...args} />;

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
