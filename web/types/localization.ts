/**
 * Centralized localization keys for type-safe translation handling.
 * This provides a single source of truth for all translation keys used in the application.
 * Keys are organized by logical sections using TypeScript namespaces.
 */
export const Localization = {
  /**
   * Frontend keys used in the main user-facing web application
   */
  Frontend: {
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

    // Errors
    componentError: 'component_error',

    // Browser notifications - organized by component
    BrowserNotifyModal: {
      unsupported: 'browser_notify_unsupported',
      unsupportedLocal: 'browser_notify_unsupported_local',
      iosTitle: 'browser_notify_ios_title',
      iosDescription: 'browser_notify_ios_description',
      iosShareButton: 'browser_notify_ios_share_button',
      iosAddToHomeScreen: 'browser_notify_ios_add_to_home_screen',
      iosAddButton: 'browser_notify_ios_add_button',
      iosNameAndTap: 'browser_notify_ios_name_and_tap',
      iosComeBack: 'browser_notify_ios_come_back',
      iosAllowPrompt: 'browser_notify_ios_allow_prompt',
      permissionWantsTo: 'browser_notify_permission_wants_to',
      showNotifications: 'browser_notify_show_notifications',
      allowButton: 'browser_notify_allow_button',
      blockButton: 'browser_notify_block_button',
      enabledTitle: 'browser_notify_enabled_title',
      enabledDescription: 'browser_notify_enabled_description',
      deniedTitle: 'browser_notify_denied_title',
      deniedDescription: 'browser_notify_denied_description',
      mainDescription: 'browser_notify_main_description',
      learnMore: 'browser_notify_learn_more',
      errorTitle: 'browser_notify_error_title',
      errorMessage: 'browser_notify_error_message',
    },

    // Offline banner messages
    offlineBasic: 'offline_basic',
    offlineNotifyOnly: 'offline_notify_only',
    offlineFediverseOnly: 'offline_fediverse_only',
    offlineNotifyAndFediverse: 'offline_notify_and_fediverse',
  },

  /**
   * Admin keys used in the admin interface
   */
  Admin: {
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

    Chat: {
      moderationMessagesSent: 'Admin.Chat.moderationMessagesSent',
      moderationMessagesSent_one: 'Admin.Chat.moderationMessagesSent_one',
    },

    // EditInstanceDetails component specific keys
    EditInstanceDetails: {
      offlineMessageDescription: 'Admin.EditInstanceDetails.offlineMessageDescription',
      directoryDescription: 'Admin.EditInstanceDetails.directoryDescription',
      serverUrlRequiredForDirectory: 'Admin.EditInstanceDetails.serverUrlRequiredForDirectory',
    },

    // Logging and monitoring
    info: 'Info',
    warning: 'Warning',
    error: 'Error',
    level: 'Level',
    timestamp: 'Timestamp',
    message: 'Message',
    logs: 'Logs',
  },

  /**
   * Common keys shared across both frontend and admin interfaces
   */
  Common: {
    // Basic UI elements
    yes: 'Yes',
    no: 'No',

    // Documentation and help
    documentation: 'Documentation',
    contribute: 'Contribute',
    source: 'Source',

    // Branding
    poweredByOwncast: 'Powered by Owncast',
    poweredByOwncastVersion: 'powered_by_owncast_version',
  },

  /**
   * Testing keys used for development and testing purposes
   */
  Testing: {
    testing: 'testing_string',
    another: 'another_test',
    simpleKey: 'Testing.simpleKey',
    itemCount: 'Testing.itemCount',
    messageCount: 'Testing.messageCount',
    noPluralKey: 'Testing.noPluralKey',
  },
} as const;

/**
 * Helper type to extract all nested values from the Localization object
 */
type NestedValues<T> = T extends object
  ? {
      [K in keyof T]: T[K] extends string ? T[K] : NestedValues<T[K]>;
    }[keyof T]
  : never;

/**
 * Type representing all valid localization keys.
 * This ensures type safety when using translation keys with nested structure.
 */
export type LocalizationKey = NestedValues<typeof Localization>;

/**
 * Helper type to get the actual string value from a localization key.
 * This can be useful for type inference in components.
 */
export type LocalizationValue<T extends LocalizationKey> = T;
