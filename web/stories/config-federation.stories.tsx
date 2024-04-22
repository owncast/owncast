import { Meta } from '@storybook/react';
import ConfigFederation from '../pages/admin/config-federation';

const meta = {
  title: 'owncast/Admin/Config/Federation',
  component: ConfigFederation,
  parameters: {},
} satisfies Meta<typeof ConfigFederation>;

export default meta;

const Template = arguments_ => <ConfigFederation {...arguments_} />;

export const Default = Template.bind({});

Default.args = {};
