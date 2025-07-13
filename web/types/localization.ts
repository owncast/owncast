/**
 * Centralized localization keys for type-safe translation handling.
 * This provides a single source of truth for all translation keys used in the application.
 */
export const Localization = {
  // Common UI strings
  chatOffline: 'Chat is offline',
  chatDisabled: 'Chat is disabled',
  chatWillBeAvailable: 'Chat will be available when the stream is live',

  // Testing and demo strings
  testing: 'testing_string',
  another: 'another_test',

  // Complex HTML translations with variables
  helloWorld: 'hello_world',
  notificationMessage: 'notification_message',
  complexMessage: 'complex_message',

  // Stream and timing related
  lastLiveAgo: 'Last live ago',
  currentViewers: 'Current viewers',
  maxViewers: 'Max viewers this stream',

  // Common actions and labels
  notify: 'Notify',
  follow: 'Follow',
  connected: 'Connected',
  yes: 'Yes',
  no: 'No',

  // Status messages
  noStreamActive: 'No stream is active',
  healthyStream: 'Healthy Stream',
  playbackHealth: 'Playback Health',

  // Navigation and accessibility
  skipToPlayer: 'Skip to player',
  skipToContent: 'Skip to page content',
  skipToFooter: 'Skip to footer',

  // Documentation and help
  documentation: 'Documentation',
  contribute: 'Contribute',
  source: 'Source',

  // Branding
  poweredByOwncast: 'Powered by Owncast',

  // Error and logging
  info: 'Info',
  warning: 'Warning',
  error: 'Error',
  level: 'Level',
  timestamp: 'Timestamp',
  message: 'Message',
  logs: 'Logs',

  // Settings and configuration
  settings: 'settings',
  overriddenViaCommandLine: 'Overridden via command line',

  // Social and directory
  stayUpdated: 'Stay updated!',
  fediverse: 'Add your Owncast instance to the Fediverse',
  owncastDirectory: 'Find an audience on the Owncast Directory',

  // Streaming setup
  useBroadcastingSoftware: 'Use your broadcasting software',
  embedVideo: 'Embed your video onto other sites',

  // Emoji admin page
  emojis: 'Emojis',
  emojiPageDescription: 'Here you can upload new custom emojis for usage in the chat. When uploading a new emoji, the filename without extension will be used as emoji name. Additionally, emoji names are case-insensitive. For best results, ensure all emoji have unique names.',
  emojiUploadBulkGuide: 'Want to upload custom emojis in bulk? Check out our <a href="https://owncast.online/docs/chat/emoji" rel="noopener noreferrer" target="_blank">Emoji guide</a>.',
  uploadNewEmoji: 'Upload new emoji',
  deleteEmoji: 'Delete emoji',

  // Placeholder for future translations
  simpleKey: 'simple_key',
} as const;

/**
 * Type representing all valid localization keys.
 * This ensures type safety when using translation keys.
 */
export type LocalizationKey = (typeof Localization)[keyof typeof Localization];

/**
 * Helper type to get the actual string value from a localization key.
 * This can be useful for type inference in components.
 */
export type LocalizationValue<T extends LocalizationKey> = T;
