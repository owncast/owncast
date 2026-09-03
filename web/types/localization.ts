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
    unmute: 'Frontend.unmute',

    // Navigation and accessibility
    skipToPlayer: 'Skip to player',
    skipToContent: 'Skip to page content',
    skipToFooter: 'Skip to footer',
    showAllTabs: 'Frontend.showAllTabs',

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

    // Public schedule tab
    Schedule: {
      tab: 'Frontend.Schedule.tab',
      addToCalendar: 'Frontend.Schedule.addToCalendar',
      today: 'Frontend.Schedule.today',
      month: 'Frontend.Schedule.month',
      week: 'Frontend.Schedule.week',
      list: 'Frontend.Schedule.list',
      loading: 'Frontend.Schedule.loading',
      error: 'Frontend.Schedule.error',
      retry: 'Frontend.Schedule.retry',
      noEvents: 'Frontend.Schedule.noEvents',
      cancelled: 'Frontend.Schedule.cancelled',
      timezone: 'Frontend.Schedule.timezone',
      duration: 'Frontend.Schedule.duration',
      liveNow: 'Frontend.Schedule.liveNow',
      countdownLiveIn: 'Frontend.Schedule.countdownLiveIn',
      countdownLiveNow: 'Frontend.Schedule.countdownLiveNow',
      countdownDays: 'Frontend.Schedule.countdownDays',
      countdownHours: 'Frontend.Schedule.countdownHours',
      countdownMinutes: 'Frontend.Schedule.countdownMinutes',
      countdownSeconds: 'Frontend.Schedule.countdownSeconds',
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

    // Appearance page
    Appearance: {
      pluginStylingActive: 'Admin.Appearance.pluginStylingActive',
      pluginStylingDescription: 'Admin.Appearance.pluginStylingDescription',
      alsoSetByPlugin: 'Admin.Appearance.alsoSetByPlugin',
      alsoSetByPluginTooltip: 'Admin.Appearance.alsoSetByPluginTooltip',
    },

    // AutoplaySelector component specific keys
    Autoplay: {
      title: 'Admin.Autoplay.title',
      description: 'Admin.Autoplay.description',
      optionOffLabel: 'Admin.Autoplay.optionOffLabel',
      optionOffDescription: 'Admin.Autoplay.optionOffDescription',
      optionAlwaysLabel: 'Admin.Autoplay.optionAlwaysLabel',
      optionAlwaysDescription: 'Admin.Autoplay.optionAlwaysDescription',
      optionSoundOnlyLabel: 'Admin.Autoplay.optionSoundOnlyLabel',
      optionSoundOnlyDescription: 'Admin.Autoplay.optionSoundOnlyDescription',
    },

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
      searchPlaceholder: 'Admin.Help.searchPlaceholder',
      upgradeAvailable: 'Admin.Help.upgradeAvailable',
      upgradeLink: 'Admin.Help.upgradeLink',
      gettingStarted: 'Admin.Help.gettingStarted',
      gettingStartedLink: 'Admin.Help.gettingStartedLink',
      fixTitle: 'Admin.Help.fixTitle',
      fixDescription: 'Admin.Help.fixDescription',
      fixButton: 'Admin.Help.fixButton',
      communityTitle: 'Admin.Help.communityTitle',
      communityDescription: 'Admin.Help.communityDescription',
      communityButton: 'Admin.Help.communityButton',
      bugTitle: 'Admin.Help.bugTitle',
      bugDescription: 'Admin.Help.bugDescription',
      bugButton: 'Admin.Help.bugButton',
      commonTasks: 'Admin.Help.commonTasks',
      groupSetup: 'Admin.Help.groupSetup',
      groupCustomize: 'Admin.Help.groupCustomize',
      groupGrow: 'Admin.Help.groupGrow',
      groupExtend: 'Admin.Help.groupExtend',
      openSettings: 'Admin.Help.openSettings',
      learnMore: 'Admin.Help.learnMore',
      taskBroadcasting: 'Admin.Help.taskBroadcasting',
      taskVideoQuality: 'Admin.Help.taskVideoQuality',
      taskStorage: 'Admin.Help.taskStorage',
      taskSsl: 'Admin.Help.taskSsl',
      taskWebsite: 'Admin.Help.taskWebsite',
      taskChat: 'Admin.Help.taskChat',
      taskNotifications: 'Admin.Help.taskNotifications',
      taskFediverse: 'Admin.Help.taskFediverse',
      taskDirectory: 'Admin.Help.taskDirectory',
      taskEmbed: 'Admin.Help.taskEmbed',
      taskPlugins: 'Admin.Help.taskPlugins',
      taskApis: 'Admin.Help.taskApis',
      supportInfo: 'Admin.Help.supportInfo',
      supportInfoDescription: 'Admin.Help.supportInfoDescription',
      copySupportInfo: 'Admin.Help.copySupportInfo',
      copied: 'Admin.Help.copied',
      moreResources: 'Admin.Help.moreResources',
      documentation: 'Admin.Help.documentation',
      releaseNotes: 'Admin.Help.releaseNotes',
      discussions: 'Admin.Help.discussions',
      fediverse: 'Admin.Help.fediverse',
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
      chartNoData: 'Admin.ViewerInfo.chartNoData',
      viewers: 'Admin.ViewerInfo.viewers',
    },

    // Per-client playback health table on the viewer info page
    PlaybackClients: {
      title: 'Admin.PlaybackClients.title',
      description: 'Admin.PlaybackClients.description',
      player: 'Admin.PlaybackClients.player',
      location: 'Admin.PlaybackClients.location',
      watchTime: 'Admin.PlaybackClients.watchTime',
      state: 'Admin.PlaybackClients.state',
      playing: 'Admin.PlaybackClients.playing',
      paused: 'Admin.PlaybackClients.paused',
      buffering: 'Admin.PlaybackClients.buffering',
      seeking: 'Admin.PlaybackClients.seeking',
      ended: 'Admin.PlaybackClients.ended',
      unknownState: 'Admin.PlaybackClients.unknownState',
      speed: 'Admin.PlaybackClients.speed',
      quality: 'Admin.PlaybackClients.quality',
      latency: 'Admin.PlaybackClients.latency',
      segmentDownload: 'Admin.PlaybackClients.segmentDownload',
      errors: 'Admin.PlaybackClients.errors',
      qualityChanges: 'Admin.PlaybackClients.qualityChanges',
      sourceClient: 'Admin.PlaybackClients.sourceClient',
      sourceServer: 'Admin.PlaybackClients.sourceServer',
      sourceUnknown: 'Admin.PlaybackClients.sourceUnknown',
      serverUnmeasurable: 'Admin.PlaybackClients.serverUnmeasurable',
      serverUnmeasurableDescription: 'Admin.PlaybackClients.serverUnmeasurableDescription',
      lastMeasured: 'Admin.PlaybackClients.lastMeasured',
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
      streamTitle: 'Admin.FeaturedStreams.streamTitle',
      url: 'Admin.FeaturedStreams.url',
      status: 'Admin.FeaturedStreams.status',
      online: 'Admin.FeaturedStreams.online',
      offline: 'Admin.FeaturedStreams.offline',
      pendingApproval: 'Admin.FeaturedStreams.pendingApproval',
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

      // Incoming feature requests (servers asking to feature this stream)
      featureRequestsTitle: 'Admin.FeaturedStreams.featureRequestsTitle',
      featureRequestsDescription: 'Admin.FeaturedStreams.featureRequestsDescription',
      approveButton: 'Admin.FeaturedStreams.approveButton',
      rejectButton: 'Admin.FeaturedStreams.rejectButton',
      featureRequestApproved: 'Admin.FeaturedStreams.featureRequestApproved',
      featureRequestRejected: 'Admin.FeaturedStreams.featureRequestRejected',
      failedToApprove: 'Admin.FeaturedStreams.failedToApprove',
      failedToReject: 'Admin.FeaturedStreams.failedToReject',

      // Tabs splitting streams you feature from directories featuring you
      streamsYouFeatureTab: 'Admin.FeaturedStreams.streamsYouFeatureTab',
      featuringYouTab: 'Admin.FeaturedStreams.featuringYouTab',

      // Directories that are featuring/listing this server
      approvalResent: 'Admin.FeaturedStreams.approvalResent',
      directoryListingsDescription: 'Admin.FeaturedStreams.directoryListingsDescription',
      directoryListingsEmpty: 'Admin.FeaturedStreams.directoryListingsEmpty',
      directoryListingsTitle: 'Admin.FeaturedStreams.directoryListingsTitle',
      directoryRemoved: 'Admin.FeaturedStreams.directoryRemoved',
      failedToRemoveDirectory: 'Admin.FeaturedStreams.failedToRemoveDirectory',
      failedToResendApproval: 'Admin.FeaturedStreams.failedToResendApproval',
      removeFromDirectoryButton: 'Admin.FeaturedStreams.removeFromDirectoryButton',
      removeFromDirectoryConfirm: 'Admin.FeaturedStreams.removeFromDirectoryConfirm',
      resendApprovalButton: 'Admin.FeaturedStreams.resendApprovalButton',
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
      autoplayUpdated: 'Admin.StatusMessages.autoplayUpdated',
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
      webhookSecret: 'Admin.Webhooks.webhookSecret',
      createNewWebhook: 'Admin.Webhooks.createNewWebhook',
      webhookUrlPlaceholder: 'Admin.Webhooks.webhookUrlPlaceholder',
      selectEvents: 'Admin.Webhooks.selectEvents',
      selectAll: 'Admin.Webhooks.selectAll',
    },

    // Schedule page
    Schedule: {
      title: 'Admin.Schedule.title',
      pageDescription: 'Admin.Schedule.pageDescription',
      enableLabel: 'Admin.Schedule.enableLabel',
      enableTip: 'Admin.Schedule.enableTip',
      countdownLabel: 'Admin.Schedule.countdownLabel',
      countdownTip: 'Admin.Schedule.countdownTip',
      chatOpenMinutesLabel: 'Admin.Schedule.chatOpenMinutesLabel',
      chatOpenMinutesTip: 'Admin.Schedule.chatOpenMinutesTip',
      chatOpenMinutesDisabled: 'Admin.Schedule.chatOpenMinutesDisabled',
      chatOpenMinutes5: 'Admin.Schedule.chatOpenMinutes5',
      chatOpenMinutes10: 'Admin.Schedule.chatOpenMinutes10',
      chatOpenMinutes30: 'Admin.Schedule.chatOpenMinutes30',
      chatOpenMinutes60: 'Admin.Schedule.chatOpenMinutes60',
      reminderLabel: 'Admin.Schedule.reminderLabel',
      reminderTip: 'Admin.Schedule.reminderTip',
      eventReminderLabel: 'Admin.Schedule.eventReminderLabel',
      eventReminderTip: 'Admin.Schedule.eventReminderTip',
      firstReminderLabel: 'Admin.Schedule.firstReminderLabel',
      firstReminderTip: 'Admin.Schedule.firstReminderTip',
      secondReminderLabel: 'Admin.Schedule.secondReminderLabel',
      secondReminderTip: 'Admin.Schedule.secondReminderTip',
      reminderDisabled: 'Admin.Schedule.reminderDisabled',
      reminder15Minutes: 'Admin.Schedule.reminder15Minutes',
      reminder30Minutes: 'Admin.Schedule.reminder30Minutes',
      reminder60Minutes: 'Admin.Schedule.reminder60Minutes',
      reminder2Hours: 'Admin.Schedule.reminder2Hours',
      reminder24Hours: 'Admin.Schedule.reminder24Hours',
      addEvent: 'Admin.Schedule.addEvent',
      editEvent: 'Admin.Schedule.editEvent',
      editAction: 'Admin.Schedule.editAction',
      upcomingEvents: 'Admin.Schedule.upcomingEvents',
      recurringSchedules: 'Admin.Schedule.recurringSchedules',
      nextOccurrences: 'Admin.Schedule.nextOccurrences',
      nameLabel: 'Admin.Schedule.nameLabel',
      descriptionLabel: 'Admin.Schedule.descriptionLabel',
      repeatsLabel: 'Admin.Schedule.repeatsLabel',
      doesNotRepeat: 'Admin.Schedule.doesNotRepeat',
      weekly: 'Admin.Schedule.weekly',
      onDays: 'Admin.Schedule.onDays',
      startingFrom: 'Admin.Schedule.startingFrom',
      atTime: 'Admin.Schedule.atTime',
      dateLabel: 'Admin.Schedule.dateLabel',
      timeLabel: 'Admin.Schedule.timeLabel',
      endsOnOptional: 'Admin.Schedule.endsOnOptional',
      durationLabel: 'Admin.Schedule.durationLabel',
      timezoneLabel: 'Admin.Schedule.timezoneLabel',
      cancelAction: 'Admin.Schedule.cancelAction',
      cancelEventAction: 'Admin.Schedule.cancelEventAction',
      deleteEventConfirm: 'Admin.Schedule.deleteEventConfirm',
      deleteAction: 'Admin.Schedule.deleteAction',
      deleteEventAction: 'Admin.Schedule.deleteEventAction',
      deleteSeriesConfirm: 'Admin.Schedule.deleteSeriesConfirm',
      savedToast: 'Admin.Schedule.savedToast',
      deletedToast: 'Admin.Schedule.deletedToast',
      cancelledToast: 'Admin.Schedule.cancelledToast',
      uneditableRule: 'Admin.Schedule.uneditableRule',
      namePlaceholder: 'Admin.Schedule.namePlaceholder',
      statusScheduled: 'Admin.Schedule.statusScheduled',
      statusCancelled: 'Admin.Schedule.statusCancelled',
      statusRecurring: 'Admin.Schedule.statusRecurring',
      sidebarTitle: 'Admin.Schedule.sidebarTitle',
      columnName: 'Admin.Schedule.columnName',
      columnWhen: 'Admin.Schedule.columnWhen',
      columnDuration: 'Admin.Schedule.columnDuration',
      columnStatus: 'Admin.Schedule.columnStatus',
      columnRepeats: 'Admin.Schedule.columnRepeats',
      durationValue: 'Admin.Schedule.durationValue',
      recurrenceOn: 'Admin.Schedule.recurrenceOn',
      recurrenceAt: 'Admin.Schedule.recurrenceAt',
      recurrenceUntil: 'Admin.Schedule.recurrenceUntil',
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

    // Users management page
    Users: {
      pageTitle: 'Admin.Users.pageTitle',
      pageDescription: 'Admin.Users.pageDescription',
      searchPlaceholder: 'Admin.Users.searchPlaceholder',
      banReason: 'Admin.Users.banReason',
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
      browsePreviewAlt: 'Admin.Plugins.browsePreviewAlt',
      browseAuthor: 'Admin.Plugins.browseAuthor',
      browseFilterPermissions: 'Admin.Plugins.browseFilterPermissions',
      browseFilterAuthor: 'Admin.Plugins.browseFilterAuthor',
      browseFilterCategory: 'Admin.Plugins.browseFilterCategory',
      browseFilteredEmpty: 'Admin.Plugins.browseFilteredEmpty',
      browseClearFilters: 'Admin.Plugins.browseClearFilters',
      browseSearch: 'Admin.Plugins.browseSearch',
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
      configTab: 'Admin.Plugins.configTab',
      configSave: 'Admin.Plugins.configSave',
      configSaveError: 'Admin.Plugins.configSaveError',
      configSaved: 'Admin.Plugins.configSaved',
      configLoadError: 'Admin.Plugins.configLoadError',
      authSettingsTab: 'Admin.Plugins.authSettingsTab',
      authSettingsDescription: 'Admin.Plugins.authSettingsDescription',
      authSettingsAccessMode: 'Admin.Plugins.authSettingsAccessMode',
      authSettingsWebsiteOnly: 'Admin.Plugins.authSettingsWebsiteOnly',
      authSettingsWebsiteOnlyDescription: 'Admin.Plugins.authSettingsWebsiteOnlyDescription',
      authSettingsWebsiteAndStream: 'Admin.Plugins.authSettingsWebsiteAndStream',
      authSettingsWebsiteAndStreamDescription:
        'Admin.Plugins.authSettingsWebsiteAndStreamDescription',
      authSettingsWebsiteStreamAndStatus: 'Admin.Plugins.authSettingsWebsiteStreamAndStatus',
      authSettingsWebsiteStreamAndStatusDescription:
        'Admin.Plugins.authSettingsWebsiteStreamAndStatusDescription',
      authSettingsSave: 'Admin.Plugins.authSettingsSave',
      authSettingsSaved: 'Admin.Plugins.authSettingsSaved',
      authSettingsLoadError: 'Admin.Plugins.authSettingsLoadError',
      authSettingsSaveError: 'Admin.Plugins.authSettingsSaveError',
      configEmpty: 'Admin.Plugins.configEmpty',
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
      // Display names for the canonical plugin category taxonomy. Slugs
      // mirror the registry's category list; keep in lock-step with the
      // categoryNameKey map in BrowseRegistry.tsx.
      Categories: {
        chatBots: 'Admin.Plugins.Categories.chatBots',
        chatFilters: 'Admin.Plugins.Categories.chatFilters',
        moderation: 'Admin.Plugins.Categories.moderation',
        authentication: 'Admin.Plugins.Categories.authentication',
        themes: 'Admin.Plugins.Categories.themes',
        overlays: 'Admin.Plugins.Categories.overlays',
        notifications: 'Admin.Plugins.Categories.notifications',
        integrations: 'Admin.Plugins.Categories.integrations',
        video: 'Admin.Plugins.Categories.video',
        analytics: 'Admin.Plugins.Categories.analytics',
        games: 'Admin.Plugins.Categories.games',
        viewerUi: 'Admin.Plugins.Categories.viewerUi',
        adminUtilities: 'Admin.Plugins.Categories.adminUtilities',
        examples: 'Admin.Plugins.Categories.examples',
        other: 'Admin.Plugins.Categories.other',
      },
      // Permission identifiers mirror services/plugins/hostfns.go; keep
      // the keys here in lock-step with those constants.
      Permissions: {
        storageKv: 'Admin.Plugins.Permissions.storageKv',
        storageUpload: 'Admin.Plugins.Permissions.storageUpload',
        storageFs: 'Admin.Plugins.Permissions.storageFs',
        storageSql: 'Admin.Plugins.Permissions.storageSql',
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
        fediverseInbound: 'Admin.Plugins.Permissions.fediverseInbound',
        fediversePost: 'Admin.Plugins.Permissions.fediversePost',
        videoconfigRead: 'Admin.Plugins.Permissions.videoconfigRead',
        videoconfigWrite: 'Admin.Plugins.Permissions.videoconfigWrite',
        uiModify: 'Admin.Plugins.Permissions.uiModify',
        chatFilter: 'Admin.Plugins.Permissions.chatFilter',
        usersRegister: 'Admin.Plugins.Permissions.usersRegister',
        authGate: 'Admin.Plugins.Permissions.authGate',
      },
      // Short plain-language labels for each permission. Shown on
      // permission Tags in the plugins list with the full description
      // surfaced via tooltip on hover.
      PermissionNames: {
        storageKv: 'Admin.Plugins.PermissionNames.storageKv',
        storageUpload: 'Admin.Plugins.PermissionNames.storageUpload',
        storageFs: 'Admin.Plugins.PermissionNames.storageFs',
        storageSql: 'Admin.Plugins.PermissionNames.storageSql',
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
        fediverseInbound: 'Admin.Plugins.PermissionNames.fediverseInbound',
        fediversePost: 'Admin.Plugins.PermissionNames.fediversePost',
        videoconfigRead: 'Admin.Plugins.PermissionNames.videoconfigRead',
        videoconfigWrite: 'Admin.Plugins.PermissionNames.videoconfigWrite',
        uiModify: 'Admin.Plugins.PermissionNames.uiModify',
        chatFilter: 'Admin.Plugins.PermissionNames.chatFilter',
        usersRegister: 'Admin.Plugins.PermissionNames.usersRegister',
        authGate: 'Admin.Plugins.PermissionNames.authGate',
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
