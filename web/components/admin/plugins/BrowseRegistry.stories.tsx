import React from 'react';
import { Meta, StoryFn } from '@storybook/react';
import { BrowseRegistry, BrowseRegistryProps, RegistryPlugin } from './BrowseRegistry';

export default {
  title: 'owncast/Admin/Plugins/BrowseRegistry',
  component: BrowseRegistry,
  parameters: {
    docs: {
      description: {
        component:
          'The Browse tab of the admin plugins screen. Lists publicly-available registry plugins as cards, each showing the plugin name, version, size, the author display name, summary, optional preview screenshot, declared permissions, and tags. The action button reflects whether the plugin is not installed (Install), already installed (Installed), or has a newer version available (Update).',
      },
    },
  },
} as Meta<typeof BrowseRegistry>;

const noop = async () => {};

const plugins: RegistryPlugin[] = [
  {
    slug: 'welcome-bot',
    name: 'Welcome Bot',
    authorName: 'Gabe Kangas',
    summary: 'Greets new chat participants with a configurable message when they join.',
    tags: ['chat', 'moderation'],
    latest: {
      version: '1.2.0',
      sizeBytes: 48 * 1024,
      manifest: { permissions: ['chat.send', 'chat.read'] },
    },
  },
  {
    slug: 'now-playing',
    name: 'Now Playing',
    authorName: 'Owncast Community',
    summary: 'Shows the currently playing track pulled from an external music service.',
    tags: ['overlay'],
    latest: {
      version: '0.4.1',
      sizeBytes: 120 * 1024,
      manifest: { permissions: ['network.outbound'] },
    },
  },
  {
    // No authorName: the "by ..." line should be omitted entirely rather
    // than rendering an empty or "by undefined" string.
    slug: 'raid-alerts',
    name: 'Raid Alerts',
    summary: 'Plays an on-screen alert when another stream raids yours.',
    latest: {
      version: '2.0.0',
      sizeBytes: 256 * 1024,
    },
  },
];

const Template: StoryFn<BrowseRegistryProps> = args => <BrowseRegistry {...args} />;

// Primary story: a registry with several plugins, each (where available)
// showing its author display name under the title.
export const WithAuthors = Template.bind({});
WithAuthors.args = {
  installedVersions: new Map(),
  registry: plugins,
  registryLoading: false,
  registryError: null,
  onInstall: noop,
  onRetry: () => {},
};

// Exercises the three action-button states alongside the author line:
// Welcome Bot is up to date, Now Playing has an update available, and
// Raid Alerts (no author) is not installed.
export const MixedInstallStates = Template.bind({});
MixedInstallStates.args = {
  ...WithAuthors.args,
  installedVersions: new Map([
    ['welcome-bot', '1.2.0'],
    ['now-playing', '0.3.0'],
  ]),
};

export const SingleAuthoredPlugin = Template.bind({});
SingleAuthoredPlugin.args = {
  ...WithAuthors.args,
  registry: [plugins[0]],
};

export const Loading = Template.bind({});
Loading.args = {
  ...WithAuthors.args,
  registry: [],
  registryLoading: true,
};

export const Empty = Template.bind({});
Empty.args = {
  ...WithAuthors.args,
  registry: [],
};

export const CatalogUnavailable = Template.bind({});
CatalogUnavailable.args = {
  ...WithAuthors.args,
  registry: [],
  registryError: 'dial tcp: lookup registry.owncast.online: no such host',
};
