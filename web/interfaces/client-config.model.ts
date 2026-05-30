export interface ClientConfig {
  name: string;
  title?: string;
  summary: string;
  offlineMessage?: string;
  logo: string;
  tags: string[];
  nsfw: boolean;
  extraPageContent: string;
  socialHandles: SocialHandle[];
  chatDisabled: boolean;
  chatRequireAuthentication: boolean;
  externalActions: any[];
  // customStyles is the admin's CSS plus the concatenated content of
  // every loaded plugin's manifest.styles entries (the host
  // pre-merges them server-side). Theme.tsx renders this as one
  // inline <style> block.
  customStyles: string;
  // pluginTabs is the list of viewer-page tabs contributed by
  // loaded plugins via manifest.tabs. DesktopContent / MobileContent
  // render one tab per entry alongside the built-in tabs.
  pluginTabs: PluginTab[];
  appearanceVariables: Map<string, string>;
  maxSocketPayloadSize: number;
  federation: Federation;
  notifications: Notifications;
  authentication: Authentication;
  socketHostOverride?: string;
}

interface Authentication {
  indieAuthEnabled: boolean;
}

interface Federation {
  enabled: boolean;
  account: string;
  followerCount: number;
}

interface Notifications {
  browser: Browser;
}

interface Browser {
  enabled: boolean;
  publicKey: string;
}

interface SocialHandle {
  platform: string;
  url: string;
  icon: string;
}

// PluginTab is one viewer-page tab contributed by a plugin via
// manifest.tabs. Mirrors models.PluginTab on the backend.
export interface PluginTab {
  slug: string;
  title: string;
  html: string;
}

export function makeEmptyClientConfig(): ClientConfig {
  return {
    name: '',
    summary: '',
    offlineMessage: '',
    logo: '',
    tags: [],
    nsfw: false,
    extraPageContent: '',
    socialHandles: [],
    chatDisabled: false,
    chatRequireAuthentication: false,
    externalActions: [],
    customStyles: '',
    pluginTabs: [],
    appearanceVariables: new Map(),
    maxSocketPayloadSize: 0,
    federation: {
      enabled: false,
      account: '',
      followerCount: 0,
    },
    notifications: {
      browser: {
        enabled: false,
        publicKey: '',
      },
    },
    authentication: {
      indieAuthEnabled: false,
    },
  };
}
