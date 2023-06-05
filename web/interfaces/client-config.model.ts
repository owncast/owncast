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
  externalActions: any[];
  customStyles: string;
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
    externalActions: [],
    customStyles: '',
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
