import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { Button, ButtonProps } from './Button';
import {action} from "@storybook/addon-actions";
import {disableControls} from "../../../../utils/storybook-utils";

export default {
  title: 'atoms/Button',
  component: Button,
	argTypes: { ...disableControls('onClick') }
} as ComponentMeta<typeof Button>;

const Template: ComponentStory<typeof Button> = args => <Button {...args} />;

export const withIcon = Template.bind({});
withIcon.args = {
	iconSrc: 'https://owncast.online/images/logo.svg',
	iconAltText: '',
	onClick: action('button clicked'),
	text: 'foo',
} as ButtonProps;

export const withoutIcon = Template.bind({});
withoutIcon.args = {
	onClick: action('button clicked'),
	text: 'foo',
} as ButtonProps;
