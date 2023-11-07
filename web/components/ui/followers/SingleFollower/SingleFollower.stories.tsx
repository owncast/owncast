import { Meta } from '@storybook/react';
import { SingleFollower } from './SingleFollower';
import SingleFollowerMock from '../../../../stories/assets/mocks/single-follower.png';

const meta = {
  title: 'owncast/Components/Followers/Single Follower',
  component: SingleFollower,
  parameters: {
    design: {
      type: 'image',
      url: SingleFollowerMock,
    },
    docs: {
      description: {
        component: `Represents a single follower.`,
      },
    },
  },
} satisfies Meta<typeof SingleFollower>;

export default meta;

export const Example = {
  args: {
    follower: {
      name: 'John Doe',
      description: 'User',
      username: '@account@domain.tld',
      image: 'https://avatars0.githubusercontent.com/u/1234?s=460&v=4',
      link: 'https://yahoo.com',
    },
  },
};
