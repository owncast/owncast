import React from 'react';
import { Meta, StoryFn } from '@storybook/nextjs';
import { InstallConfirmModal, InstallConfirmModalProps } from './InstallConfirmModal';
import { Plugin, PluginPermission } from '../../../interfaces/plugin';

export default {
  title: 'owncast/Admin/Plugins/InstallConfirmModal',
  component: InstallConfirmModal,
  parameters: {
    docs: {
      description: {
        component:
          'Post-install confirmation listing the permissions a plugin manifest declares, shown before the admin enables it.',
      },
    },
  },
} as Meta<typeof InstallConfirmModal>;

const basePlugin: Plugin = {
  slug: 'example-plugin',
  name: 'Example Plugin',
  path: '/plugins/example-plugin',
  enabled: false,
  loaded: false,
  discoveredAt: '2026-01-01T00:00:00Z',
};

const Template: StoryFn<InstallConfirmModalProps> = args => <InstallConfirmModal {...args} />;

export const WithPermissions = Template.bind({});
WithPermissions.args = {
  plugin: {
    ...basePlugin,
    permissions: [
      PluginPermission.ChatSend,
      PluginPermission.UsersRead,
      PluginPermission.NetworkFetch,
    ],
    allowedHosts: ['api.example.com', '*.cdn.example.com'],
  },
  onCancel: () => {},
  onEnable: () => {},
};

export const NoPermissions = Template.bind({});
NoPermissions.args = {
  plugin: basePlugin,
  onCancel: () => {},
  onEnable: () => {},
};
