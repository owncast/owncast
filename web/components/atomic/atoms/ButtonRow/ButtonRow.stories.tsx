import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';

import { Button } from '../Button/Button';
import { ButtonRow, ButtonRowProps } from './ButtonRow';
import {action} from "@storybook/addon-actions";
import {disableControls} from "../../../../utils/storybook-utils";

export default {
  title: 'atoms/ButtonRow',
  component: ButtonRow,
  argTypes: {
    ...disableControls('children'),
  },
} as ComponentMeta<typeof ButtonRow>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ButtonRow> = ({ children }) => (<ButtonRow>{children}</ButtonRow>);

export const buttonRow = Template.bind({});
buttonRow.args = {
  children: [
		<Button text="foo" onClick={action('foo clicked')} />,
		<Button text="bar" onClick={action('bar clicked')} />,
		<Button text="tin" onClick={action('tin clicked')} />,
	],
} as ButtonRowProps;
