import { Meta } from '@storybook/react';
import ConfigFederation from './config-federation';
// import Mock from '../../../stories/assets/mocks/chatmessage-action.png';

const meta = {
  title: 'owncast/Admin/Config/Federation',
  component: ConfigFederation,
  parameters: {},
} satisfies Meta<typeof ConfigFederation>;

export default meta;

const Template = arguments_ => <ConfigFederation {...arguments_} />;

export const Default = Template.bind({});

Default.args = {};
