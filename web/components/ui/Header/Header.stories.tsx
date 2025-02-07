import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { Header } from './Header';

const meta = {
  title: 'owncast/Layout/Header',
  component: Header,
  parameters: {
    chromatic: { diffThreshold: 0.75 },
  },
} satisfies Meta<typeof Header>;

export default meta;

const Template: StoryFn<typeof Header> = args => (
  <RecoilRoot>
    <Header {...args} />
  </RecoilRoot>
);

export const ChatAvailable = {
  render: Template,

  args: {
    name: 'Example Stream Name',
    chatAvailable: true,
  },
};

export const ChatNotAvailable = {
  render: Template,

  args: {
    name: 'Example Stream Name',
    chatAvailable: false,
  },
};
