import React from 'react';
import { Menu, Dropdown } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { ComponentStory, ComponentMeta } from '@storybook/react';

const menu = (
  <Menu>
    <Menu.Item key="0">
      <a href="https://owncast.online">1st menu item</a>
    </Menu.Item>
    <Menu.Item key="1">
      <a href="https://directory.owncast.online">2nd menu item</a>
    </Menu.Item>
    <Menu.Divider />
    <Menu.Item key="3">3rd menu item</Menu.Item>
  </Menu>
);

const DropdownExample = () => (
  <Dropdown overlay={menu} trigger={['click']}>
    <button type="button" className="ant-dropdown-link" onClick={e => e.preventDefault()}>
      Click me <DownOutlined />
    </button>
  </Dropdown>
);

export default {
  title: 'owncast/Dropdown',
  component: Dropdown,
  parameters: {},
} as ComponentMeta<typeof Dropdown>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof Dropdown> = args => <DropdownExample />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
