// TS types for elements on the Config pages

// for dropdown
export interface SocialHandleDropdownItem {
  icon: string;
  platform: string;
  key: string;
}

export type FieldUpdaterFunc = (args: UpdateArgs) => void;

export interface UpdateArgs {
  value: any;
  fieldName?: string;
  path?: string;
}

export interface ApiPostArgs {
  apiPath: string;
  data: object;
  onSuccess?: (arg: any) => void;
  onError?: (arg: any) => void;
}

export interface ConfigDirectoryFields {
  enabled: boolean;
  instanceUrl: string;
}

export interface ConfigInstanceDetailsFields {
  customStyles: string;
  customJavascript: string;
  extraPageContent: string;
  logo: string;
  name: string;
  nsfw: boolean;
  socialHandles: SocialHandle[];
  streamTitle: string;
  summary: string;
  offlineMessage: string;
  tags: string[];
  title: string;
  welcomeMessage: string;
  appearanceVariables: AppearanceVariables;
}

export type CpuUsageLevel = 1 | 2 | 3 | 4 | 5;

// from data
export interface SocialHandle {
  platform: string;
  url: string;
}

export interface VideoVariant {
  key?: number; // unique identifier generated on client side just for ant table rendering
  cpuUsageLevel: CpuUsageLevel;
  framerate: number;

  audioPassthrough: boolean;
  audioBitrate: number;
  videoPassthrough: boolean;
  videoBitrate: number;

  scaledWidth: number;
  scaledHeight: number;

  name: string;
}
export interface VideoSettingsFields {
  latencyLevel: number;
  videoQualityVariants: VideoVariant[];
  cpuUsageLevel: CpuUsageLevel;
}

export interface S3Field {
  acl?: string;
  accessKey: string;
  bucket: string;
  enabled: boolean;
  endpoint: string;
  region: string;
  secret: string;
  pathPrefix: string;
  forcePathStyle: boolean;
}

type AppearanceVariables = {
  [key: string]: string;
};

export interface ExternalAction {
  title: string;
  description: string;
  url: string;
  openExternally: boolean;
}

// PluginStyleInfo describes the page styling one enabled plugin
// contributes. The Appearance config uses it to tell the admin that
// plugin styles are combined with their own colors, and to flag the
// swatches a plugin also sets. declaredVars holds the theme custom
// properties the plugin declares (without the leading `--`, e.g.
// "theme-color-action"); it can be empty when a plugin styles the page
// without touching a recognized appearance token.
export interface PluginStyleInfo {
  slug: string;
  name: string;
  declaredVars: string[];
}

export interface Federation {
  enabled: boolean;
  isPrivate: boolean;
  username: string;
  goLiveMessage: string;
  showEngagement: boolean;
  hideFollowersTab: boolean;
  blockedDomains: string[];
}

export interface BrowserNotification {
  enabled: boolean;
  goLiveMessage: string;
}

export interface DiscordNotification {
  enabled: boolean;
  webhook: string;
  goLiveMessage: string;
}

export interface NotificationsConfig {
  browser: BrowserNotification;
  discord: DiscordNotification;
}

export interface Health {
  healthy: boolean;
  healthyPercentage: number;
  message: string;
  representation: number;
}

export interface StreamKey {
  key: string;
  comment: string;
}

export interface ConfigDetails {
  externalActions: ExternalAction[];
  styleContributors: PluginStyleInfo[];
  ffmpegPath: string;
  instanceDetails: ConfigInstanceDetailsFields;
  rtmpServerPort: string;
  rtmpServerAddress: string;
  s3: S3Field;
  streamKeys: StreamKey[];
  streamKeyOverridden: boolean;
  adminPassword: string;
  videoSettings: VideoSettingsFields;
  webServerPort: string;
  socketHostOverride: string;
  videoServingEndpoint: string;
  yp: ConfigDirectoryFields;
  supportedCodecs: string[];
  videoCodec: string;
  forbiddenUsernames: string[];
  suggestedUsernames: string[];
  chatDisabled: boolean;
  chatSpamProtectionEnabled: boolean;
  chatSlurFilterEnabled: boolean;
  chatRequireAuthentication: boolean;
  federation: Federation;
  notifications: NotificationsConfig;
  chatJoinMessagesEnabled: boolean;
  chatEstablishedUserMode: boolean;
  hideViewerCount: boolean;
  disableSearchIndexing: boolean;
}
