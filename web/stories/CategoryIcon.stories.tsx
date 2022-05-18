import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import CategoryIcon from '../components/ui/CategoryIcon/CategoryIcon';

export default {
  title: 'owncast/Components/Category icon',
  component: CategoryIcon,
  parameters: {},
} as ComponentMeta<typeof CategoryIcon>;

const Template: ComponentStory<typeof CategoryIcon> = args => <CategoryIcon {...args} />;

export const Game = Template.bind({});
Game.args = {
  tags: ['games', 'fun'],
};

export const Chat = Template.bind({});
Chat.args = {
  tags: ['blah', 'chat', 'games', 'fun'],
};

export const Conference = Template.bind({});
Conference.args = {
  tags: ['blah', 'conference', 'games', 'fun'],
};
