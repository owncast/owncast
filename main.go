package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	// Embed IANA timezone data so recurrence expansion with named zones
	// (e.g. America/Los_Angeles) works on hosts without a system tzdata,
	// notably stock Windows.
	_ "time/tzdata"

	log "github.com/sirupsen/logrus"

	fediverseauth "github.com/owncast/owncast/auth/fediverse"
	indieauthlib "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/scheduleeventsrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/pluginhost"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/dispatcher"
	"github.com/owncast/owncast/services/notifications"
	"github.com/owncast/owncast/services/rtmp"
	"github.com/owncast/owncast/services/schedule"
	"github.com/owncast/owncast/services/scheduler"
	"github.com/owncast/owncast/services/stalefeaturedcheckservice"
	"github.com/owncast/owncast/services/stream"
	"github.com/owncast/owncast/services/webhooks"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers"
	"github.com/owncast/owncast/webserver/handlers/admin"
	"github.com/owncast/owncast/webserver/handlers/auth/fediverse"
	"github.com/owncast/owncast/webserver/handlers/auth/indieauth"
	"github.com/owncast/owncast/webserver/handlers/moderation"
	"github.com/owncast/owncast/webserver/router"
	"github.com/owncast/owncast/webserver/router/middleware"
	"github.com/owncast/owncast/yp"
)

var (
	dbFile                         = flag.String("database", "", "Path to the database file.")
	logDirectory                   = flag.String("logdir", "", "Directory where logs will be written to")
	backupDirectory                = flag.String("backupdir", "", "Directory where backups will be written to")
	enableDebugOptions             = flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")
	enableVerboseLogging           = flag.Bool("enableVerboseLogging", false, "Enable additional logging.")
	restoreDatabaseFile            = flag.String("restoreDatabase", "", "Restore an Owncast database backup")
	newAdminPassword               = flag.String("adminpassword", "", "Set your admin password")
	newStreamKey                   = flag.String("streamkey", "", "Set a temporary stream key for this session")
	webServerPortOverride          = flag.String("webserverport", "", "Force the web server to listen on a specific port")
	webServerIPOverride            = flag.String("webserverip", "", "Force web server to listen on this IP address")
	rtmpPortOverride               = flag.Int("rtmpport", 0, "Set listen port for the RTMP server")
	followerValidationIntervalSecs = flag.Int("followervalidationinterval", 0, "Set follower validation interval in seconds")
)

// nolint:cyclop
func main() {
	flag.Parse()

	// Construct the runtime config and overlay CLI flag values onto the
	// defaults. The resulting *Config is threaded into every consumer
	// that needs runtime configuration values.
	cfg := config.NewDefault()

	if *logDirectory != "" {
		cfg.LogDirectory = *logDirectory
	}

	if *backupDirectory != "" {
		cfg.BackupDirectory = *backupDirectory
	}

	// Create the data directory if needed
	if !utils.DoesFileExists("data") {
		if err := os.Mkdir("./data", 0o700); err != nil {
			log.Fatalln("Cannot create data directory", err)
		}
	}

	// Migrate old (pre 0.1.0) emoji to new location if they exist.
	utils.MigrateCustomEmojiLocations()

	// Otherwise save the default emoji to the data directory.
	if err := datastore.SetupEmojiDirectory(); err != nil {
		log.Fatalln("Cannot set up emoji directory", err)
	}

	// Recreate the temp dir
	if utils.DoesFileExists(cfg.TempDir) {
		err := os.RemoveAll(cfg.TempDir)
		if err != nil {
			log.Fatalln("Unable to remove temp dir! Check permissions.", cfg.TempDir, err)
		}
	}
	if err := os.Mkdir(cfg.TempDir, 0o700); err != nil {
		log.Fatalln("Unable to create temp dir!", err)
	}

	configureLogging(cfg.LogDirectory, *enableDebugOptions, *enableVerboseLogging)
	log.Infoln(config.GetReleaseString())

	// Allows a user to restore a specific database backup
	if *restoreDatabaseFile != "" {
		databaseFile := cfg.DatabaseFilePath
		if *dbFile != "" {
			databaseFile = *dbFile
		}

		if err := utils.Restore(*restoreDatabaseFile, databaseFile); err != nil {
			log.Fatalln(err)
		}

		log.Println("Database has been restored.  Restart Owncast.")
		log.Exit(0)
	}

	cfg.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		cfg.DatabaseFilePath = *dbFile
	}

	dataStore, err := datastore.SetupPersistence(cfg.DatabaseFilePath, cfg.BackupDirectory)
	if err != nil {
		log.Fatalln("failed to open database", err)
	}

	// Composition root.
	//
	// Every service, handler, and repository is constructed here and
	// passed via explicit Deps structs to its consumers. main.go is the
	// only place that knows the concrete service implementations; no
	// other package reaches for a singleton via .Get() or a package-level
	// global. The order below reflects the dependency layering:
	//
	//   1. Repositories — wrap dataStore; depend on nothing else.
	//   2. ActivityPub helper builders — pure helpers built from repos.
	//   3. One-shot bootstrap — notificationsRepository.Setup() seeds
	//      the browser-push keys before any service starts.
	//   4. Leaf services that depend only on repositories — yp,
	//      middleware, indieauth, rtmp.
	//   5. Cycle-pair services — webhooks + chat + yp all need
	//      streamSvc.GetStatus, but stream needs them too. Each is
	//      constructed with a nil callback and rewired below via
	//      SetGetStatus once streamSvc exists.
	//   6. activitypub + stream — depend on the cycle-pair services.
	//   7. SetGetStatus fill-in — resolves the cycle.
	//   8. Late services — metrics + fediverseAuth, depend on stream
	//      and chat.
	//   9. HTTP handler set — admin, fediverse, indieauth, moderation,
	//      then the top-level *Handlers that wires everything for the
	//      router.
	//   10. router.Start — blocks; serves until shutdown.
	//
	// New deps land here, not as package-level globals. Construction
	// cycles get the SetGetStatus pattern, not a fresh shim.

	// Stage 1: repositories.
	configRepository := configrepository.New(dataStore)
	authRepository := authrepository.New(dataStore)
	followersRepository := followersrepository.New(dataStore)
	webhookRepository := webhookrepository.New(dataStore)
	chatMessageRepository := chatmessagerepository.New(dataStore)
	userRepository := userrepository.New(dataStore)
	notificationsRepository := notificationsrepository.New(dataStore, configRepository)
	federatedServersRepository := federatedserversrepository.New(dataStore)
	scheduleEventsRepository := scheduleeventsrepository.New(dataStore)

	// Expose globals for helper code that still uses package-level
	// Get accessors (the featured-streams ActivityPub paths). Long term
	// these callers move to dependency injection.
	configrepository.SetGlobalInstance(configRepository)
	federatedserversrepository.SetGlobalInstance(federatedServersRepository)

	handleCommandLineFlags(cfg, configRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cacheContainer := cache.New()
	defer cacheContainer.Stop()

	// Stage 2: ActivityPub helper types (Signer → Builder → Resolver).
	// Signer is the seed; Builder depends on Signer for actor public-key
	// embedding; Resolver depends on both for signing outbound IRI
	// fetches.
	apSigner := apcrypto.New(apcrypto.Deps{ConfigRepository: configRepository})
	apBuilder := apmodels.New(apmodels.Deps{ConfigRepository: configRepository, Signer: apSigner})
	apResolver := apresolvers.New(apresolvers.Deps{ConfigRepository: configRepository, Builder: apBuilder, Signer: apSigner})

	// Stage 3: one-shot bootstrap. Seeds browser-push keys + default
	// notification config before any service that reads them starts.
	notificationsRepository.Setup()

	// Stage 4: leaf services. Each depends only on repositories + cfg.
	ypSvc := yp.New(yp.Deps{
		GetStatus:        nil, // stage 7
		ConfigRepository: configRepository,
	})

	mw := middleware.New(middleware.Deps{
		ConfigRepository: configRepository,
		AuthRepository:   authRepository,
		UserRepository:   userRepository,
	})
	indieauthSvc := indieauthlib.New(indieauthlib.Deps{
		ConfigRepository: configRepository,
	})

	rtmpSvc := rtmp.New(rtmp.Deps{
		ConfigRepository: configRepository,
		Config:           cfg,
	})

	// Shared in-process event dispatcher. Constructed before its producers
	// and consumers so they all receive it via Deps (no post-construction
	// wiring): webhooks publishes events to it, chat runs inbound messages
	// through its filter chain, and the plugin host subscribes.
	eventDispatcher := dispatcher.New()

	// Stage 5: cycle-pair services. webhooks + chat both want
	// streamSvc.GetStatus, but streamSvc consumes them in turn. Build
	// them with a nil GetStatus and rewire in stage 7.
	webhooksSvc := webhooks.New(webhooks.Deps{
		GetStatus:         nil, // stage 7
		Followers:         followersRepository,
		ConfigRepository:  configRepository,
		WebhookRepository: webhookRepository,
		Events:            eventDispatcher,
	})

	chatSvc := chat.New(chat.Deps{
		GetStatus:             nil, // stage 7
		Webhooks:              webhooksSvc,
		Datastore:             dataStore,
		ConfigRepository:      configRepository,
		AuthRepository:        authRepository,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
		Events:                eventDispatcher,
	})

	// Stage 6: cycle-pair consumers. activitypub + stream sit on top of
	// the stage-5 services.
	apSvc := activitypub.New(activitypub.Deps{
		Datastore:           dataStore,
		Webhooks:            webhooksSvc,
		Chat:                chatSvc,
		ConfigRepository:    configRepository,
		FollowersRepository: followersRepository,
		Builder:             apBuilder,
		Signer:              apSigner,
		Resolver:            apResolver,
		Events:              eventDispatcher,
		Config:              cfg,
	})
	apSvc.Start()

	streamSvc := stream.New(stream.Deps{
		Rtmp:             rtmpSvc,
		Activitypub:      apSvc,
		Webhooks:         webhooksSvc,
		Chat:             chatSvc,
		YP:               ypSvc,
		Datastore:        dataStore,
		ConfigRepository: configRepository,
		Config:           cfg,
	})

	// Stage 7: resolve the stream-status cycle. webhooks/chat/yp each
	// hold a func() Status that streamSvc now provides.
	webhooksSvc.SetGetStatus(streamSvc.GetStatus)
	chatSvc.SetGetStatus(streamSvc.GetStatus)
	ypSvc.SetGetStatus(streamSvc.GetStatus)

	// Materializes scheduled stream occurrences, keeps the next-event answer
	// warm for the status endpoint, and closes chat when a pre-opened event
	// never starts.
	scheduleSvc := schedule.New(schedule.Deps{
		ScheduleEventsRepository:         scheduleEventsRepository,
		GetStatus:                        streamSvc.GetStatus,
		GetChatOpenMinutes:               configRepository.GetScheduleChatOpenMinutes,
		GetScheduleEnabled:               configRepository.GetScheduleEnabled,
		GetFederationEnabled:             configRepository.GetFederationEnabled,
		FederateScheduledEvent:           apSvc.SendScheduledEvent,
		GetScheduleReminderMessage:       configRepository.GetScheduleReminderMessage,
		GetScheduleFirstReminderMinutes:  configRepository.GetScheduleFirstReminderMinutes,
		GetScheduleSecondReminderMinutes: configRepository.GetScheduleSecondReminderMinutes,
		GetServerURL:                     configRepository.GetServerURL,
		NotifyScheduledEvent: func(message string) {
			notificationService, err := notifications.New(dataStore, configRepository)
			if err != nil {
				log.Errorln(err)
				return
			}
			notificationService.NotifyScheduledEvent(message)
		},
		SetStreamTitle: configRepository.SetStreamTitle,
		Webhooks:       webhooksSvc,
		OnMissedEventWarning: func(_ *models.ScheduledEvent) {
			if err := chatSvc.SendSystemMessage(schedule.MissedEventChatMessage, false); err != nil {
				log.Errorf("unable to send missed scheduled stream chat message: %v", err)
			}
		},
	})
	chatSvc.SetIsScheduledChatOpen(scheduleSvc.IsChatOpen)

	if err := streamSvc.Start(ctx); err != nil {
		log.Fatalln("failed to start the stream service", err)
	}
	defer streamSvc.Stop(ctx)
	scheduleSvc.Start()
	defer scheduleSvc.Stop()

	// The central scheduler owns all recurring background work. Services
	// expose an idempotent job method plus an interval constant, and get
	// registered here.
	schedulerSvc, err := scheduler.New()
	if err != nil {
		log.Fatalln("failed to create the scheduler", err)
	}

	// Marks federated peer servers offline when they stop sending
	// stream-status pings, keeping the featured-streams directory honest.
	staleFeaturedSvc := stalefeaturedcheckservice.New(stalefeaturedcheckservice.Deps{
		ConfigRepository:           configRepository,
		FederatedServersRepository: federatedServersRepository,
	})

	schedulerJobs := []scheduler.Job{
		{Name: "stale-featured-servers", Interval: stalefeaturedcheckservice.CheckInterval, RunAtStart: true, Run: staleFeaturedSvc.Run},
		{Name: "chat-data-pruner", Interval: chat.DataPruneInterval, RunAtStart: true, Run: chatSvc.RunDataPruner},
	}
	for _, job := range schedulerJobs {
		if err := schedulerSvc.Register(job); err != nil {
			log.Fatalln("failed to register scheduler job", job.Name, err)
		}
	}
	schedulerSvc.Start()
	defer schedulerSvc.Stop()

	// Stage 8: late services. metrics polls stream + chat, fediverseAuth
	// owns OTP state for the chat-side handler.
	metricsSvc := metrics.New(metrics.Deps{
		Stream:                streamSvc,
		Chat:                  chatSvc,
		ConfigRepository:      configRepository,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})
	go metricsSvc.Start()

	// Plugin host. Optional infrastructure: a failure here disables plugins
	// but must not abort startup. Constructs the WASM plugin runtime, wires
	// its host functions to the services above, and exposes the
	// /plugins/<name>/* HTTP handler mounted by the router below.
	var pluginContentHandler http.Handler
	var pluginAdminHandler http.Handler
	var viewerAuthGate func(http.Handler) http.Handler
	pluginHostInstance, err := pluginhost.New(ctx, pluginhost.Deps{
		Datastore:               dataStore,
		Chat:                    chatSvc,
		Stream:                  streamSvc,
		Activitypub:             apSvc,
		Events:                  eventDispatcher,
		ConfigRepository:        configRepository,
		UserRepository:          userRepository,
		AuthRepository:          authRepository,
		NotificationsRepository: notificationsRepository,
		ChatMessageRepository:   chatMessageRepository,
		RequireAdminAuth:        mw.RequireAdminAuth,
		IsAdminRequest:          mw.IsAdminRequest,
	})
	if err != nil {
		log.Errorln("plugin host failed to start; plugins disabled:", err)
	} else {
		pluginContentHandler = pluginHostInstance.Handler()
		pluginAdminHandler = pluginHostInstance.AdminHandler()
		viewerAuthGate = pluginHostInstance.AuthGateMiddleware()
		// The gate's access policy can disable directory listing, so YP has
		// to consult it before pinging or answering /api/yp.
		ypSvc.SetDirectoryAvailable(pluginHostInstance.DirectoryAvailable)
		defer pluginHostInstance.Stop(ctx)
	}

	// Stage 9: HTTP handler set. *Handlers is the dispatcher the router
	// binds methods on; the sub-handlers (admin, fediverse, indieauth,
	// moderation) hold their own narrower deps.
	// Plugin host getters: nil when the plugin host is disabled (boot
	// off or failed). Wired into both the admin handler set (for the
	// appearance style-contributors report) and the viewer handler set.
	var pluginActions func() []models.ExternalAction
	var pluginCSSContent func() []byte
	var pluginJSContent func() []byte
	var pluginPageContent func(*http.Request) []byte
	var pluginTabs func(*http.Request) []models.PluginTab
	var pluginStyleContributors func() []models.PluginStyleInfo
	if pluginHostInstance != nil {
		pluginActions = pluginHostInstance.Actions
		pluginCSSContent = pluginHostInstance.StylesContent
		pluginJSContent = pluginHostInstance.ScriptsContent
		pluginPageContent = pluginHostInstance.PageContent
		pluginTabs = pluginHostInstance.Tabs
		pluginStyleContributors = pluginHostInstance.StyleContributors
	}

	adminHandlers := admin.New(admin.Deps{
		Stream:                   streamSvc,
		Rtmp:                     rtmpSvc,
		Activitypub:              apSvc,
		Webhooks:                 webhooksSvc,
		Chat:                     chatSvc,
		Metrics:                  metricsSvc,
		ConfigRepository:         configRepository,
		AuthRepository:           authRepository,
		FollowersRepository:      followersRepository,
		WebhookRepository:        webhookRepository,
		ChatMessageRepository:    chatMessageRepository,
		UserRepository:           userRepository,
		ScheduleEventsRepository: scheduleEventsRepository,
		Schedule:                 scheduleSvc,
		APBuilder:                apBuilder,
		APSigner:                 apSigner,
		Config:                   cfg,
		PluginStyleContributors:  pluginStyleContributors,
	})

	fediverseAuthSvc := fediverseauth.New()
	fediverseAuthSvc.Start()

	fediverseHandler := fediverse.New(fediverse.Deps{
		Activitypub:      apSvc,
		Chat:             chatSvc,
		FediverseAuth:    fediverseAuthSvc,
		ConfigRepository: configRepository,
		UserRepository:   userRepository,
	})

	indieauthHandler := indieauth.New(indieauth.Deps{
		Chat:           chatSvc,
		UserRepository: userRepository,
		IndieAuth:      indieauthSvc,
		Middleware:     mw,
	})

	moderationHandler := moderation.New(moderation.Deps{
		Chat:                  chatSvc,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})

	h := handlers.NewHandlers(handlers.Deps{
		Cache:                    cacheContainer,
		Stream:                   streamSvc,
		Chat:                     chatSvc,
		Admin:                    adminHandlers,
		Activitypub:              apSvc,
		Fediverse:                fediverseHandler,
		IndieAuth:                indieauthHandler,
		Moderation:               moderationHandler,
		Middleware:               mw,
		YP:                       ypSvc,
		Metrics:                  metricsSvc,
		ConfigRepository:         configRepository,
		FollowersRepository:      followersRepository,
		ChatMessageRepository:    chatMessageRepository,
		UserRepository:           userRepository,
		NotificationsRepository:  notificationsRepository,
		ScheduleEventsRepository: scheduleEventsRepository,
		Schedule:                 scheduleSvc,
		APBuilder:                apBuilder,
		Config:                   cfg,
		PluginActions:            pluginActions,
		PluginCSSContent:         pluginCSSContent,
		PluginJSContent:          pluginJSContent,
		PluginPageContent:        pluginPageContent,
		PluginTabs:               pluginTabs,
	})

	// Stage 10: serve. Blocks until shutdown.
	if err := router.Start(cfg, *enableVerboseLogging, h, mw, apSvc.Controllers(), pluginContentHandler, pluginAdminHandler, viewerAuthGate); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}
}

func handleCommandLineFlags(cfg *config.Config, configRepository configrepository.ConfigRepository) {
	if *newAdminPassword != "" {
		if err := configRepository.SetAdminPassword(*newAdminPassword); err != nil {
			log.Errorln("Error setting your admin password.", err)
			log.Exit(1)
		} else {
			log.Infoln("Admin password changed")
		}
	}

	if *newStreamKey != "" {
		log.Println("Temporary stream key is set for this session.")
		cfg.TemporaryStreamKey = *newStreamKey
	}

	// Set the web server port
	if *webServerPortOverride != "" {
		portNumber, err := strconv.Atoi(*webServerPortOverride)
		if err != nil {
			log.Warnln(err)
			return
		}

		log.Println("Saving new web server port number to", portNumber)
		if err := configRepository.SetHTTPPortNumber(float64(portNumber)); err != nil {
			log.Errorln(err)
		}
	}
	cfg.WebServerPort = configRepository.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		if err := configRepository.SetHTTPListenAddress(*webServerIPOverride); err != nil {
			log.Errorln(err)
		}
	}
	cfg.WebServerIP = configRepository.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		if err := configRepository.SetRTMPPortNumber(float64(*rtmpPortOverride)); err != nil {
			log.Errorln(err)
		}
	}

	// Set the follower validation interval
	if *followerValidationIntervalSecs > 0 {
		cfg.FollowerValidationInterval = time.Duration(*followerValidationIntervalSecs) * time.Second
		log.Printf("Follower validation interval set to %v", cfg.FollowerValidationInterval)
	}
}

func configureLogging(logDirectory string, enableDebugFeatures bool, enableVerboseLogging bool) {
	logging.Setup(logDirectory, enableDebugFeatures, enableVerboseLogging)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
