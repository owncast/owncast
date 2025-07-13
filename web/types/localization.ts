/**
 * Centralized localization keys for type-safe translation handling.
 * This provides a single source of truth for all translation keys used in the application.
 * Keys are organized by logical sections of the application.
 */
export const Localization = {
  // ===== FRONTEND SECTION =====
  // Keys used in the main user-facing web application

  // Chat interface
  chatOffline: 'Chat is offline',
  chatDisabled: 'Chat is disabled',
  chatWillBeAvailable: 'Chat will be available when the stream is live',

  // Stream information and statistics
  lastLiveAgo: 'Last live ago',
  currentViewers: 'Current viewers',
  maxViewers: 'Max viewers this stream',
  noStreamActive: 'No stream is active',
  healthyStream: 'Healthy Stream',
  playbackHealth: 'Playback Health',

  // User actions and interactions
  notify: 'Notify',
  follow: 'Follow',
  connected: 'Connected',

  // Navigation and accessibility
  skipToPlayer: 'Skip to player',
  skipToContent: 'Skip to page content',
  skipToFooter: 'Skip to footer',

  // Social and external services
  stayUpdated: 'Stay updated!',
  fediverse: 'Add your Owncast instance to the Fediverse',
  owncastDirectory: 'Find an audience on the Owncast Directory',

  // Streaming setup and integration
  useBroadcastingSoftware: 'Use your broadcasting software',
  embedVideo: 'Embed your video onto other sites',

  // Complex HTML translations with variables
  helloWorld: 'hello_world',
  notificationMessage: 'notification_message',
  complexMessage: 'complex_message',

  // ===== ADMIN SECTION =====
  // Keys used in the admin interface

  // Emoji management
  emojis: 'Emojis',
  emojiPageDescription:
    'Here you can upload new custom emojis for usage in the chat. When uploading a new emoji, the filename without extension will be used as emoji name. Additionally, emoji names are case-insensitive. For best results, ensure all emoji have unique names.',
  emojiUploadBulkGuide:
    'Want to upload custom emojis in bulk? Check out our <a href="https://owncast.online/docs/chat/emoji" rel="noopener noreferrer" target="_blank">Emoji guide</a>.',
  uploadNewEmoji: 'Upload new emoji',
  deleteEmoji: 'Delete emoji',

  // Settings and configuration
  settings: 'settings',
  overriddenViaCommandLine: 'Overridden via command line',

  // Logging and monitoring
  info: 'Info',
  warning: 'Warning',
  error: 'Error',
  level: 'Level',
  timestamp: 'Timestamp',
  message: 'Message',
  logs: 'Logs',

  // ===== COMMON SECTION =====
  // Keys shared across both frontend and admin interfaces

  // Basic UI elements
  yes: 'Yes',
  no: 'No',

  // Documentation and help
  documentation: 'Documentation',
  contribute: 'Contribute',
  source: 'Source',

  // Branding
  poweredByOwncast: 'Powered by Owncast',

  // ===== TESTING SECTION =====
  // Keys used for development and testing purposes

  testing: 'testing_string',
  another: 'another_test',
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
