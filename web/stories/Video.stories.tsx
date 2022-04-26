import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import VideoPoster from '../components/video/VideoPoster';

export default {
  title: 'owncast/VideoPoster',
  component: VideoPoster,
  parameters: {},
} as ComponentMeta<typeof VideoPoster>;

const VideoPosterExample = () => <VideoPoster />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof VideoPoster> = args => <VideoPosterExample />;

export const Basic = Template.bind({});
