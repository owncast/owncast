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
  // customStyles is the admin's custom CSS. Theme.tsx renders it last
  // in the appearance cascade (after pluginStyles and the appearance
  // variables), so it wins over plugin styling on overlap.
  customStyles: string;
  // pluginStyles is the concatenated CSS contributed by loaded plugins
  // (manifest.styles + on_page_styles output). Theme.tsx renders it as
  // a baseline <style> block before the appearance variables and
  // customStyles, so admin appearance settings layer on top.
  pluginStyles: string;
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
  // hideFollowersTab hides the public "Followers" tab on the viewer
  // page while leaving the rest of the social features active.
  hideFollowersTab: boolean;
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
  slug: string; // composite unique key: pluginSlug/tabSlug
  pluginSlug: string; // source plugin identifier
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
    pluginStyles: '',
    pluginTabs: [],
    appearanceVariables: new Map(),
    maxSocketPayloadSize: 0,
    federation: {
      enabled: false,
      account: '',
      followerCount: 0,
      hideFollowersTab: false,
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
