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
    chatEmbedTitle: 'Frontend.chatEmbedTitle',
    chatOffline: 'Chat is offline',
    chatDisabled: 'Chat is disabled',
    chatWillBeAvailable: 'Chat will be available when the stream is live',

    // Stream information and statistics
    lastLiveAgo: 'Last live {{timeAgo}} ago',
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
      learnMoreAboutNotifications: 'Frontend.BrowserNotifyModal.learnMoreAboutNotifications',
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

    // Header component
    Header: {
      skipToPlayer: 'Frontend.Header.skipToPlayer',
      skipToOfflineMessage: 'Frontend.Header.skipToOfflineMessage',
      skipToContent: 'Frontend.Header.skipToContent',
      skipToFooter: 'Frontend.Header.skipToFooter',
      chatWillBeAvailable: 'Frontend.Header.chatWillBeAvailable',
      chatOffline: 'Frontend.Header.chatOffline',
    },

    // Footer component
    Footer: {
      documentation: 'Frontend.Footer.documentation',
      contribute: 'Frontend.Footer.contribute',
      source: 'Frontend.Footer.source',
    },

    // Featured Streams tab component
    StreamsTab: {
      streams: 'Frontend.StreamsTab.streams',
      loadingStreams: 'Frontend.StreamsTab.loadingStreams',
      errorLoadingStreams: 'Frontend.StreamsTab.errorLoadingStreams',
      noFeaturedStreams: 'Frontend.StreamsTab.noFeaturedStreams',
      live: 'Frontend.StreamsTab.live',
      offline: 'Frontend.StreamsTab.offline',
    },

    // Chat message components
    Chat: {
      userJoined: 'Frontend.Chat.userJoined',
      userLeft: 'Frontend.Chat.userLeft',
      nameChangeText: 'Frontend.Chat.nameChangeText',
      moderatorNotification: 'Frontend.Chat.moderatorNotification',
      authenticateToChat: 'Frontend.Chat.authenticateToChat',
    },

    // Follow modal component
    FollowModal: {
      description: 'Frontend.FollowModal.description',
      learnFediverse: 'Frontend.FollowModal.learnFediverse',
      newToYou: 'Frontend.FollowModal.newToYou',
      instructions: 'Frontend.FollowModal.instructions',
      placeholder: 'Frontend.FollowModal.placeholder',
      redirectMessage: 'Frontend.FollowModal.redirectMessage',
      joinFediverse: 'Frontend.FollowModal.joinFediverse',
      followButton: 'Frontend.FollowModal.followButton',
      followError: 'Frontend.FollowModal.followError',
      unableToFollow: 'Frontend.FollowModal.unableToFollow',
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

    // EditFavicon component specific keys
    EditFavicon: {
      label: 'Admin.EditFavicon.label',
      tip: 'Admin.EditFavicon.tip',
      resetButton: 'Admin.EditFavicon.resetButton',
      resetConfirmTitle: 'Admin.EditFavicon.resetConfirmTitle',
      resetConfirmOk: 'Admin.EditFavicon.resetConfirmOk',
      resetConfirmCancel: 'Admin.EditFavicon.resetConfirmCancel',
    },

    // VideoVariantForm component specific keys
    VideoVariantForm: {
      bitrateDisabledPassthrough: 'Admin.VideoVariantForm.bitrateDisabledPassthrough',
      bitrateValueKbps: 'Admin.VideoVariantForm.bitrateValueKbps',
      bitrateGoodForSlow: 'Admin.VideoVariantForm.bitrateGoodForSlow',
      bitrateGoodForMost: 'Admin.VideoVariantForm.bitrateGoodForMost',
      bitrateGoodForHigh: 'Admin.VideoVariantForm.bitrateGoodForHigh',
    },

    // Hardware monitoring page
    HardwareInfo: {
      title: 'Admin.HardwareInfo.title',
      pleaseWait: 'Admin.HardwareInfo.pleaseWait',
      noDetails: 'Admin.HardwareInfo.noDetails',
      cpu: 'Admin.HardwareInfo.cpu',
      memory: 'Admin.HardwareInfo.memory',
      disk: 'Admin.HardwareInfo.disk',
      used: 'Admin.HardwareInfo.used',
    },

    // Help page
    Help: {
      title: 'Admin.Help.title',
      configureInstance: 'Admin.Help.configureInstance',
      learnMore: 'Admin.Help.learnMore',
      configureBroadcasting: 'Admin.Help.configureBroadcasting',
      embedStream: 'Admin.Help.embedStream',
      customizeWebsite: 'Admin.Help.customizeWebsite',
      tweakVideo: 'Admin.Help.tweakVideo',
      useStorage: 'Admin.Help.useStorage',
      foundBug: 'Admin.Help.foundBug',
      bugPlease: 'Admin.Help.bugPlease',
      letUsKnow: 'Admin.Help.letUsKnow',
      generalQuestion: 'Admin.Help.generalQuestion',
      generalAnswered: 'Admin.Help.generalAnswered',
      faq: 'Admin.Help.faq',
      orExist: 'Admin.Help.orExist',
      discussions: 'Admin.Help.discussions',
      buildAddons: 'Admin.Help.buildAddons',
      buildTools: 'Admin.Help.buildTools',
      developerApis: 'Admin.Help.developerApis',
      troubleshooting: 'Admin.Help.troubleshooting',
      fixProblems: 'Admin.Help.fixProblems',
      documentation: 'Admin.Help.documentation',
      readDocs: 'Admin.Help.readDocs',
      commonTasks: 'Admin.Help.commonTasks',
      other: 'Admin.Help.other',
    },

    // Log table component
    LogTable: {
      level: 'Admin.LogTable.level',
      info: 'Admin.LogTable.info',
      warning: 'Admin.LogTable.warning',
      error: 'Admin.LogTable.error',
      timestamp: 'Admin.LogTable.timestamp',
      message: 'Admin.LogTable.message',
      logs: 'Admin.LogTable.logs',
    },

    // News feed component
    NewsFeed: {
      link: 'Admin.NewsFeed.link',
      noNews: 'Admin.NewsFeed.noNews',
      title: 'Admin.NewsFeed.title',
    },

    // Viewer info page
    ViewerInfo: {
      title: 'Admin.ViewerInfo.title',
      currentStream: 'Admin.ViewerInfo.currentStream',
      last12Hours: 'Admin.ViewerInfo.last12Hours',
      last24Hours: 'Admin.ViewerInfo.last24Hours',
      last7Days: 'Admin.ViewerInfo.last7Days',
      last30Days: 'Admin.ViewerInfo.last30Days',
      last3Months: 'Admin.ViewerInfo.last3Months',
      last6Months: 'Admin.ViewerInfo.last6Months',
      currentViewers: 'Admin.ViewerInfo.currentViewers',
      maxViewersThisStream: 'Admin.ViewerInfo.maxViewersThisStream',
      maxViewersLastStream: 'Admin.ViewerInfo.maxViewersLastStream',
      maxViewers: 'Admin.ViewerInfo.maxViewers',
      pleaseWait: 'Admin.ViewerInfo.pleaseWait',
      noData: 'Admin.ViewerInfo.noData',
      viewers: 'Admin.ViewerInfo.viewers',
    },

    // Featured streams admin component
    FeaturedStreams: {
      // Page titles and descriptions
      pageTitle: 'Admin.FeaturedStreams.pageTitle',
      pageDescription: 'Admin.FeaturedStreams.pageDescription',
      pageDescriptionSecondary: 'Admin.FeaturedStreams.pageDescriptionSecondary',

      // Button labels
      featureStreamButton: 'Admin.FeaturedStreams.featureStreamButton',
      unfeatureButton: 'Admin.FeaturedStreams.unfeatureButton',

      // Modal content
      modalTitle: 'Admin.FeaturedStreams.modalTitle',
      streamUrlLabel: 'Admin.FeaturedStreams.streamUrlLabel',
      streamUrlPlaceholder: 'Admin.FeaturedStreams.streamUrlPlaceholder',
      streamUrlHelp: 'Admin.FeaturedStreams.streamUrlHelp',
      featureStreamAction: 'Admin.FeaturedStreams.featureStreamAction',

      // Table headers
      streamName: 'Admin.FeaturedStreams.streamName',
      url: 'Admin.FeaturedStreams.url',
      status: 'Admin.FeaturedStreams.status',
      online: 'Admin.FeaturedStreams.online',
      offline: 'Admin.FeaturedStreams.offline',
      lastChecked: 'Admin.FeaturedStreams.lastChecked',
      never: 'Admin.FeaturedStreams.never',
      added: 'Admin.FeaturedStreams.added',
      actions: 'Admin.FeaturedStreams.actions',
      totalStreams: 'Admin.FeaturedStreams.totalStreams',

      // Confirmation dialog
      unfeatureConfirm: 'Admin.FeaturedStreams.unfeatureConfirm',
      confirmYes: 'Admin.FeaturedStreams.confirmYes',
      confirmNo: 'Admin.FeaturedStreams.confirmNo',

      // Requirements
      streamRequirements: 'Admin.FeaturedStreams.streamRequirements',
      requirementOwncast: 'Admin.FeaturedStreams.requirementOwncast',
      requirementHttps: 'Admin.FeaturedStreams.requirementHttps',
      requirementDefaultPort: 'Admin.FeaturedStreams.requirementDefaultPort',
      requirementFeaturedStreams: 'Admin.FeaturedStreams.requirementFeaturedStreams',

      // Validation errors
      enterStreamUrl: 'Admin.FeaturedStreams.enterStreamUrl',
      enterValidUrl: 'Admin.FeaturedStreams.enterValidUrl',
      onlyHttpsSupported: 'Admin.FeaturedStreams.onlyHttpsSupported',
      onlyDefaultPortSupported: 'Admin.FeaturedStreams.onlyDefaultPortSupported',
      invalidUrl: 'Admin.FeaturedStreams.invalidUrl',

      // Success/Error messages
      streamFeaturedSuccess: 'Admin.FeaturedStreams.streamFeaturedSuccess',
      streamUnfeaturedSuccess: 'Admin.FeaturedStreams.streamUnfeaturedSuccess',
      failedToFeature: 'Admin.FeaturedStreams.failedToFeature',
      failedToUnfeature: 'Admin.FeaturedStreams.failedToUnfeature',
      unsupportedFeaturedStreams: 'Admin.FeaturedStreams.unsupportedFeaturedStreams',

      // Warnings
      socialFeaturesRequired: 'Admin.FeaturedStreams.socialFeaturesRequired',
      socialFeaturesRequiredDesc: 'Admin.FeaturedStreams.socialFeaturesRequiredDesc',
      federationSettings: 'Admin.FeaturedStreams.federationSettings',
    },

    // Logging and monitoring
    info: 'Admin.info',
    warning: 'Admin.warning',
    error: 'Admin.error',
    level: 'Admin.level',
    timestamp: 'Admin.timestamp',
    message: 'Admin.message',
    logs: 'Admin.logs',

    // Form status messages
    StatusMessages: {
      updated: 'Admin.StatusMessages.updated',
      tagsUpdated: 'Admin.StatusMessages.tagsUpdated',
      variantsUpdated: 'Admin.StatusMessages.variantsUpdated',
      videoCodecUpdated: 'Admin.StatusMessages.videoCodecUpdated',
      latencyBufferUpdated: 'Admin.StatusMessages.latencyBufferUpdated',
      deletingEmoji: 'Admin.StatusMessages.deletingEmoji',
      emojiDeleted: 'Admin.StatusMessages.emojiDeleted',
      convertingEmoji: 'Admin.StatusMessages.convertingEmoji',
      uploadingEmoji: 'Admin.StatusMessages.uploadingEmoji',
      emojiUploadedSuccessfully: 'Admin.StatusMessages.emojiUploadedSuccessfully',
      thereWasAnError: 'Admin.StatusMessages.thereWasAnError',
      fileSizeTooBig: 'Admin.StatusMessages.fileSizeTooBig',
      fileTypeNotSupported: 'Admin.StatusMessages.fileTypeNotSupported',
      pleaseEnterTag: 'Admin.StatusMessages.pleaseEnterTag',
      tagAlreadyUsed: 'Admin.StatusMessages.tagAlreadyUsed',
    },

    // Actions page
    Actions: {
      title: 'Admin.Actions.title',
      description: 'Admin.Actions.description',
      readMoreLink: 'Admin.Actions.readMoreLink',
      createNewAction: 'Admin.Actions.createNewAction',
      createNewActionTitle: 'Admin.Actions.createNewActionTitle',
      editActionTitle: 'Admin.Actions.editActionTitle',
      modalDescription: 'Admin.Actions.modalDescription',
      onlyHttpsSupported: 'Admin.Actions.onlyHttpsSupported',
      readMoreAboutActions: 'Admin.Actions.readMoreAboutActions',
      selectActionType: 'Admin.Actions.selectActionType',
      linkOrEmbedUrl: 'Admin.Actions.linkOrEmbedUrl',
      customHtml: 'Admin.Actions.customHtml',
      htmlEmbedPlaceholder: 'Admin.Actions.htmlEmbedPlaceholder',
      urlPlaceholder: 'Admin.Actions.urlPlaceholder',
      titlePlaceholder: 'Admin.Actions.titlePlaceholder',
      descriptionPlaceholder: 'Admin.Actions.descriptionPlaceholder',
      iconPlaceholder: 'Admin.Actions.iconPlaceholder',
      optionalBackgroundColor: 'Admin.Actions.optionalBackgroundColor',
      openExternally: 'Admin.Actions.openExternally',
    },

    // Webhooks page
    Webhooks: {
      createNewWebhook: 'Admin.Webhooks.createNewWebhook',
      webhookUrlPlaceholder: 'Admin.Webhooks.webhookUrlPlaceholder',
      selectEvents: 'Admin.Webhooks.selectEvents',
      selectAll: 'Admin.Webhooks.selectAll',
    },

    // Access Tokens page
    AccessTokens: {
      createNewAccessToken: 'Admin.AccessTokens.createNewAccessToken',
      nameDescription: 'Admin.AccessTokens.nameDescription',
      namePlaceholder: 'Admin.AccessTokens.namePlaceholder',
      selectPermissions: 'Admin.AccessTokens.selectPermissions',
      cannotEditAfterCreated: 'Admin.AccessTokens.cannotEditAfterCreated',
      selectAll: 'Admin.AccessTokens.selectAll',
    },

    // Plugins admin (overview list, per-plugin detail, sidebar submenu)
    Plugins: {
      sidebarTitle: 'Admin.Plugins.sidebarTitle',
      overview: 'Admin.Plugins.overview',
      pageTitle: 'Admin.Plugins.pageTitle',
      pageDescription: 'Admin.Plugins.pageDescription',
      refresh: 'Admin.Plugins.refresh',
      pluginColumn: 'Admin.Plugins.pluginColumn',
      permissionsColumn: 'Admin.Plugins.permissionsColumn',
      statusColumn: 'Admin.Plugins.statusColumn',
      enabledColumn: 'Admin.Plugins.enabledColumn',
      none: 'Admin.Plugins.none',
      pluginFailedToLoad: 'Admin.Plugins.pluginFailedToLoad',
      statusError: 'Admin.Plugins.statusError',
      statusRunning: 'Admin.Plugins.statusRunning',
      statusEnabledNotLoaded: 'Admin.Plugins.statusEnabledNotLoaded',
      statusDisabled: 'Admin.Plugins.statusDisabled',
      statusAutoDisabled: 'Admin.Plugins.statusAutoDisabled',
      statusAutoDisabledTooltip: 'Admin.Plugins.statusAutoDisabledTooltip',
      statusPendingApproval: 'Admin.Plugins.statusPendingApproval',
      statusPendingApprovalTooltip: 'Admin.Plugins.statusPendingApprovalTooltip',
      approveButton: 'Admin.Plugins.approveButton',
      approveTooltip: 'Admin.Plugins.approveTooltip',
      uploadButton: 'Admin.Plugins.uploadButton',
      uploadSuccess: 'Admin.Plugins.uploadSuccess',
      tabInstalled: 'Admin.Plugins.tabInstalled',
      tabBrowse: 'Admin.Plugins.tabBrowse',
      browseEmpty: 'Admin.Plugins.browseEmpty',
      browseInstall: 'Admin.Plugins.browseInstall',
      browseInstalled: 'Admin.Plugins.browseInstalled',
      browseUpdate: 'Admin.Plugins.browseUpdate',
      browseUnavailableTitle: 'Admin.Plugins.browseUnavailableTitle',
      browseUnavailableDescription: 'Admin.Plugins.browseUnavailableDescription',
      browseUnavailableRetry: 'Admin.Plugins.browseUnavailableRetry',
      updateAvailable: 'Admin.Plugins.updateAvailable',
      updateConfirmTitle: 'Admin.Plugins.updateConfirmTitle',
      updateConfirmOk: 'Admin.Plugins.updateConfirmOk',
      updateConfirmCancel: 'Admin.Plugins.updateConfirmCancel',
      installConfirmTitle: 'Admin.Plugins.installConfirmTitle',
      installConfirmPrompt: 'Admin.Plugins.installConfirmPrompt',
      installConfirmNoPermissions: 'Admin.Plugins.installConfirmNoPermissions',
      installConfirmEnable: 'Admin.Plugins.installConfirmEnable',
      installConfirmCancel: 'Admin.Plugins.installConfirmCancel',
      installEnabledSuccess: 'Admin.Plugins.installEnabledSuccess',
      uninstallTooltip: 'Admin.Plugins.uninstallTooltip',
      uninstallAria: 'Admin.Plugins.uninstallAria',
      uninstallConfirmTitle: 'Admin.Plugins.uninstallConfirmTitle',
      uninstallConfirmDescription: 'Admin.Plugins.uninstallConfirmDescription',
      uninstallConfirmOk: 'Admin.Plugins.uninstallConfirmOk',
      uninstallConfirmCancel: 'Admin.Plugins.uninstallConfirmCancel',
      uninstallSuccess: 'Admin.Plugins.uninstallSuccess',
      reloadTooltip: 'Admin.Plugins.reloadTooltip',
      openPluginAdmin: 'Admin.Plugins.openPluginAdmin',
      toggleAria: 'Admin.Plugins.toggleAria',
      pluginErrorTitle: 'Admin.Plugins.pluginErrorTitle',
      configure: 'Admin.Plugins.configure',
      permissionsTab: 'Admin.Plugins.permissionsTab',
      commandsTab: 'Admin.Plugins.commandsTab',
      instructionsTab: 'Admin.Plugins.instructionsTab',
      instructionsLoadError: 'Admin.Plugins.instructionsLoadError',
      permissionColumnHeader: 'Admin.Plugins.permissionColumnHeader',
      commandColumnHeader: 'Admin.Plugins.commandColumnHeader',
      descriptionColumnHeader: 'Admin.Plugins.descriptionColumnHeader',
      allowedHostsLabel: 'Admin.Plugins.allowedHostsLabel',
      modOnlyTag: 'Admin.Plugins.modOnlyTag',
      noPermissionsTitle: 'Admin.Plugins.noPermissionsTitle',
      noPermissionsDescription: 'Admin.Plugins.noPermissionsDescription',
      notFoundTitle: 'Admin.Plugins.notFoundTitle',
      notFoundDescription: 'Admin.Plugins.notFoundDescription',
      errorTitle: 'Admin.Plugins.errorTitle',
      // Permission identifiers mirror services/plugins/hostfns.go; keep
      // the keys here in lock-step with those constants.
      Permissions: {
        storageKv: 'Admin.Plugins.Permissions.storageKv',
        storageUpload: 'Admin.Plugins.Permissions.storageUpload',
        storageFs: 'Admin.Plugins.Permissions.storageFs',
        chatSend: 'Admin.Plugins.Permissions.chatSend',
        chatHistory: 'Admin.Plugins.Permissions.chatHistory',
        chatModerate: 'Admin.Plugins.Permissions.chatModerate',
        networkFetch: 'Admin.Plugins.Permissions.networkFetch',
        eventsEmit: 'Admin.Plugins.Permissions.eventsEmit',
        httpServe: 'Admin.Plugins.Permissions.httpServe',
        httpSse: 'Admin.Plugins.Permissions.httpSse',
        serverRead: 'Admin.Plugins.Permissions.serverRead',
        notificationsSend: 'Admin.Plugins.Permissions.notificationsSend',
        usersRead: 'Admin.Plugins.Permissions.usersRead',
        usersModerate: 'Admin.Plugins.Permissions.usersModerate',
        fediversePost: 'Admin.Plugins.Permissions.fediversePost',
        videoconfigRead: 'Admin.Plugins.Permissions.videoconfigRead',
        videoconfigWrite: 'Admin.Plugins.Permissions.videoconfigWrite',
        uiModify: 'Admin.Plugins.Permissions.uiModify',
        chatFilter: 'Admin.Plugins.Permissions.chatFilter',
      },
      // Short plain-language labels for each permission. Shown on
      // permission Tags in the plugins list with the full description
      // surfaced via tooltip on hover.
      PermissionNames: {
        storageKv: 'Admin.Plugins.PermissionNames.storageKv',
        storageUpload: 'Admin.Plugins.PermissionNames.storageUpload',
        storageFs: 'Admin.Plugins.PermissionNames.storageFs',
        chatSend: 'Admin.Plugins.PermissionNames.chatSend',
        chatHistory: 'Admin.Plugins.PermissionNames.chatHistory',
        chatModerate: 'Admin.Plugins.PermissionNames.chatModerate',
        networkFetch: 'Admin.Plugins.PermissionNames.networkFetch',
        eventsEmit: 'Admin.Plugins.PermissionNames.eventsEmit',
        httpServe: 'Admin.Plugins.PermissionNames.httpServe',
        httpSse: 'Admin.Plugins.PermissionNames.httpSse',
        serverRead: 'Admin.Plugins.PermissionNames.serverRead',
        notificationsSend: 'Admin.Plugins.PermissionNames.notificationsSend',
        usersRead: 'Admin.Plugins.PermissionNames.usersRead',
        usersModerate: 'Admin.Plugins.PermissionNames.usersModerate',
        fediversePost: 'Admin.Plugins.PermissionNames.fediversePost',
        videoconfigRead: 'Admin.Plugins.PermissionNames.videoconfigRead',
        videoconfigWrite: 'Admin.Plugins.PermissionNames.videoconfigWrite',
        uiModify: 'Admin.Plugins.PermissionNames.uiModify',
        chatFilter: 'Admin.Plugins.PermissionNames.chatFilter',
      },
    },
  },

  /**
   * Common keys shared across both frontend and admin interfaces
   */
  Common: {
    // Branding
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
