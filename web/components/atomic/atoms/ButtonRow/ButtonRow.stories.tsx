import React from 'react';
import { action } from '@storybook/addon-actions';
import { ComponentStory, ComponentMeta } from '@storybook/react';

import { disableControls } from '../../../../utils/storybook-utils';
import { Button } from '../Button/Button';
import { ButtonRow as ButtonRowComponent, ButtonRowProps } from './ButtonRow';

export default {
  title: 'atoms/ButtonRow',
  component: ButtonRowComponent,
  argTypes: {
    ...disableControls('children'),
  },
} as ComponentMeta<typeof ButtonRowComponent>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof ButtonRowComponent> = ({ children }) => (
  <ButtonRowComponent>{children}</ButtonRowComponent>
);

export const ButtonRow = Template.bind({});
ButtonRow.args = {
  children: [
    <Button text="foo" onClick={action('foo clicked')} />,
    <Button text="bar" onClick={action('bar clicked')} />,
    <Button text="tin" onClick={action('tin clicked')} />,
  ],
} as ButtonRowProps;
