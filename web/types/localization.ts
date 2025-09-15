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
    helloWorld: 'Frontend.helloWorld',
    notificationMessage: 'Frontend.notificationMessage',
    complexMessage: 'Frontend.complexMessage',

    // Errors
    componentError: 'Frontend.componentError',

    // Browser notifications - organized by component
    BrowserNotifyModal: {
      unsupported: 'Frontend.BrowserNotifyModal.unsupported',
      unsupportedLocal: 'Frontend.BrowserNotifyModal.unsupportedLocal',
      iosTitle: 'Frontend.BrowserNotifyModal.iosTitle',
      iosDescription: 'Frontend.BrowserNotifyModal.iosDescription',
      iosShareButton: 'Frontend.BrowserNotifyModal.iosShareButton',
      iosAddToHomeScreen: 'Frontend.BrowserNotifyModal.iosAddToHomeScreen',
      iosAddButton: 'Frontend.BrowserNotifyModal.iosAddButton',
      iosNameAndTap: 'Frontend.BrowserNotifyModal.iosNameAndTap',
      iosComeBack: 'Frontend.BrowserNotifyModal.iosComeBack',
      iosAllowPrompt: 'Frontend.BrowserNotifyModal.iosAllowPrompt',
      permissionWantsTo: 'Frontend.BrowserNotifyModal.permissionWantsTo',
      showNotifications: 'Frontend.BrowserNotifyModal.showNotifications',
      allowButton: 'Frontend.BrowserNotifyModal.allowButton',
      blockButton: 'Frontend.BrowserNotifyModal.blockButton',
      enabledTitle: 'Frontend.BrowserNotifyModal.enabledTitle',
      enabledDescription: 'Frontend.BrowserNotifyModal.enabledDescription',
      deniedTitle: 'Frontend.BrowserNotifyModal.deniedTitle',
      deniedDescription: 'Frontend.BrowserNotifyModal.deniedDescription',
      mainDescription: 'Frontend.BrowserNotifyModal.mainDescription',
      learnMore: 'Frontend.BrowserNotifyModal.learnMore',
      errorTitle: 'Frontend.BrowserNotifyModal.errorTitle',
      errorMessage: 'Frontend.BrowserNotifyModal.errorMessage',
    },

    // Name change modal - organized by component
    NameChangeModal: {
      description: 'Frontend.NameChangeModal.description',
      placeholder: 'Frontend.NameChangeModal.placeholder',
      buttonText: 'Frontend.NameChangeModal.buttonText',
      colorLabel: 'Frontend.NameChangeModal.colorLabel',
      authInfo: 'Frontend.NameChangeModal.authInfo',
      overLimit: 'Frontend.NameChangeModal.overLimit',
    },

    // Offline banner messages
    offlineBasic: 'Frontend.offlineBasic',
    offlineNotifyOnly: 'Frontend.offlineNotifyOnly',
    offlineFediverseOnly: 'Frontend.offlineFediverseOnly',
    offlineNotifyAndFediverse: 'Frontend.offlineNotifyAndFediverse',
  },

  /**
   * Admin keys used in the admin interface
   */
  Admin: {
    // Emoji management
    emojis: 'Admin.emojis',
    emojiPageDescription: 'Admin.emojiPageDescription',
    emojiUploadBulkGuide: 'Admin.emojiUploadBulkGuide',
    uploadNewEmoji: 'Admin.uploadNewEmoji',
    deleteEmoji: 'Admin.deleteEmoji',

    // Settings and configuration
    settings: 'Admin.settings',
    overriddenViaCommandLine: 'Admin.overriddenViaCommandLine',

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

    // VideoVariantForm component specific keys
    VideoVariantForm: {
      bitrateDisabledPassthrough: 'Admin.VideoVariantForm.bitrateDisabledPassthrough',
      bitrateValueKbps: 'Admin.VideoVariantForm.bitrateValueKbps',
      bitrateGoodForSlow: 'Admin.VideoVariantForm.bitrateGoodForSlow',
      bitrateGoodForMost: 'Admin.VideoVariantForm.bitrateGoodForMost',
      bitrateGoodForHigh: 'Admin.VideoVariantForm.bitrateGoodForHigh',
    },

    // Logging and monitoring
    info: 'Admin.info',
    warning: 'Admin.warning',
    error: 'Admin.error',
    level: 'Admin.level',
    timestamp: 'Admin.timestamp',
    message: 'Admin.message',
    logs: 'Admin.logs',
  },

  /**
   * Common keys shared across both frontend and admin interfaces
   */
  Common: {
    // Basic UI elements
    yes: 'Common.yes',
    no: 'Common.no',

    // Documentation and help
    documentation: 'Common.documentation',
    contribute: 'Common.contribute',
    source: 'Common.source',

    // Branding
    poweredByOwncast: 'Common.poweredByOwncast',
    poweredByOwncastVersion: 'Common.poweredByOwncastVersion',
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
