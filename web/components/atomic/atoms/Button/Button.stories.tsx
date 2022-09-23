import React from 'react';
import { action } from '@storybook/addon-actions';
import { ComponentStory, ComponentMeta } from '@storybook/react';

import { disableControls } from '../../../../utils/storybook-utils';
import { Button, ButtonProps } from './Button';

export default {
  title: 'atoms/Button',
  component: Button,
  argTypes: { ...disableControls('onClick') },
} as ComponentMeta<typeof Button>;

const Template: ComponentStory<typeof Button> = args => <Button {...args} />;

export const WithIcon = Template.bind({});
WithIcon.args = {
  iconSrc: 'https://owncast.online/images/logo.svg',
  iconAltText: '',
  onClick: action('button clicked'),
  text: 'foo',
} as ButtonProps;

export const WithoutIcon = Template.bind({});
WithoutIcon.args = {
  onClick: action('button clicked'),
  text: 'foo',
} as ButtonProps;
