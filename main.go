package main

import (
	"context"
	"flag"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	indieauthlib "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/logging"
	"github.com/owncast/owncast/persistence/authrepository"
	"github.com/owncast/owncast/persistence/chatmessagerepository"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/notificationsrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/persistence/webhookrepository"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/persistence/followersrepository"
	apresolvers "github.com/owncast/owncast/services/activitypub/resolvers"
	"github.com/owncast/owncast/services/cache"
	"github.com/owncast/owncast/services/chat"
	chatevents "github.com/owncast/owncast/services/chat/events"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/services/rtmp"
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

	if *logDirectory != "" {
		config.LogDirectory = *logDirectory
	}

	if *backupDirectory != "" {
		config.BackupDirectory = *backupDirectory
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
	if utils.DoesFileExists(config.TempDir) {
		err := os.RemoveAll(config.TempDir)
		if err != nil {
			log.Fatalln("Unable to remove temp dir! Check permissions.", config.TempDir, err)
		}
	}
	if err := os.Mkdir(config.TempDir, 0o700); err != nil {
		log.Fatalln("Unable to create temp dir!", err)
	}

	configureLogging(*enableDebugOptions, *enableVerboseLogging)
	log.Infoln(config.GetReleaseString())

	// Allows a user to restore a specific database backup
	if *restoreDatabaseFile != "" {
		databaseFile := config.DatabaseFilePath
		if *dbFile != "" {
			databaseFile = *dbFile
		}

		if err := utils.Restore(*restoreDatabaseFile, databaseFile); err != nil {
			log.Fatalln(err)
		}

		log.Println("Database has been restored.  Restart Owncast.")
		log.Exit(0)
	}

	config.EnableDebugFeatures = *enableDebugOptions

	if *dbFile != "" {
		config.DatabaseFilePath = *dbFile
	}

	if err := datastore.SetupPersistence(config.DatabaseFilePath); err != nil {
		log.Fatalln("failed to open database", err)
	}

	// Composition root: construct services here and inject them into the
	// components that consume them. As more packages migrate off package-
	// level singletons into services/<domain>/, their constructors join
	// this block.
	configRepository := configrepository.New(datastore.GetDatastore())
	authRepository := authrepository.New(datastore.GetDatastore())
	followersRepository := followersrepository.New(datastore.GetDatastore())
	webhookRepository := webhookrepository.New(datastore.GetDatastore())
	chatMessageRepository := chatmessagerepository.New(datastore.GetDatastore())
	userRepository := userrepository.New(datastore.GetDatastore())

	handleCommandLineFlags(configRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cacheContainer := cache.New()
	defer cacheContainer.Stop()

	// Install the package-level configRepository handle in every helper
	// package that has stateless lookups. main.go is the only place that
	// knows the concrete repo — all other callers receive it via Deps.
	apmodels.SetConfigRepository(configRepository)
	apcrypto.SetConfigRepository(configRepository)
	apresolvers.SetConfigRepository(configRepository)
	chatevents.SetConfigRepository(configRepository)
	metrics.SetConfigRepository(configRepository)
	metrics.SetUserRepository(userRepository)
	yp.SetConfigRepository(configRepository)
	middleware.SetConfigRepository(configRepository)
	middleware.SetAuthRepository(authRepository)
	middleware.SetUserRepository(userRepository)
	indieauthlib.SetConfigRepository(configRepository)
	notificationsrepository.SetConfigRepository(configRepository)
	metrics.SetChatMessageRepository(chatMessageRepository)

	rtmpSvc := rtmp.New(rtmp.Deps{
		ConfigRepository: configRepository,
	})

	// Webhooks needs the stream status callback and the followers
	// repository. Construct it before activitypub (which uses it) and
	// before stream (which uses it to dispatch start/stop events).
	webhooksSvc := webhooks.New(webhooks.Deps{
		GetStatus:         nil, // wired below once streamSvc exists
		Followers:         followersRepository,
		ConfigRepository:  configRepository,
		WebhookRepository: webhookRepository,
	})

	// chat is constructed first because activitypub.Deps and
	// stream.Deps both need it; webhooks doesn't (yet — once chat
	// migrates to receive webhook events via the dispatcher, this
	// changes).
	chatSvc := chat.New(chat.Deps{
		GetStatus:             nil, // wired below once streamSvc exists
		Webhooks:              webhooksSvc,
		ConfigRepository:      configRepository,
		AuthRepository:        authRepository,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})

	apSvc := activitypub.New(activitypub.Deps{
		Datastore:           datastore.GetDatastore(),
		Webhooks:            webhooksSvc,
		Chat:                chatSvc,
		ConfigRepository:    configRepository,
		FollowersRepository: followersRepository,
	})
	apSvc.Start()

	streamSvc := stream.New(stream.Deps{
		Rtmp:             rtmpSvc,
		Activitypub:      apSvc,
		Webhooks:         webhooksSvc,
		Chat:             chatSvc,
		ConfigRepository: configRepository,
	})

	// Now that stream + AP are constructed, finish wiring the webhook
	// service deps and the chat getStatus (small construction cycle:
	// webhooks needs stream.GetStatus and ap.Followers; chat needs
	// stream.GetStatus; stream and AP both need *webhooks.Service +
	// *chat.Service).
	webhooksSvc.SetDeps(webhooks.Deps{
		GetStatus:         streamSvc.GetStatus,
		Followers:         followersRepository,
		ConfigRepository:  configRepository,
		WebhookRepository: webhookRepository,
	})
	chatSvc.SetGetStatus(streamSvc.GetStatus)

	if err := streamSvc.Start(ctx); err != nil {
		log.Fatalln("failed to start the stream service", err)
	}
	defer streamSvc.Stop(ctx)

	go metrics.Start(streamSvc, chatSvc)

	adminHandlers := admin.New(admin.Deps{
		Stream:                streamSvc,
		Rtmp:                  rtmpSvc,
		Activitypub:           apSvc,
		Webhooks:              webhooksSvc,
		Chat:                  chatSvc,
		ConfigRepository:      configRepository,
		AuthRepository:        authRepository,
		FollowersRepository:   followersRepository,
		WebhookRepository:     webhookRepository,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})

	fediverseHandler := fediverse.New(fediverse.Deps{
		Activitypub:      apSvc,
		Chat:             chatSvc,
		ConfigRepository: configRepository,
		UserRepository:   userRepository,
	})

	indieauthHandler := indieauth.New(indieauth.Deps{
		Chat:           chatSvc,
		UserRepository: userRepository,
	})

	moderationHandler := moderation.New(moderation.Deps{
		Chat:                  chatSvc,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})

	h := handlers.NewHandlers(handlers.Deps{
		Cache:                 cacheContainer,
		Stream:                streamSvc,
		Chat:                  chatSvc,
		Admin:                 adminHandlers,
		Activitypub:           apSvc,
		Fediverse:             fediverseHandler,
		IndieAuth:             indieauthHandler,
		Moderation:            moderationHandler,
		ConfigRepository:      configRepository,
		FollowersRepository:   followersRepository,
		ChatMessageRepository: chatMessageRepository,
		UserRepository:        userRepository,
	})

	if err := router.Start(*enableVerboseLogging, h, apSvc.Controllers()); err != nil {
		log.Fatalln("failed to start/run the router", err)
	}
}

func handleCommandLineFlags(configRepository configrepository.ConfigRepository) {
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
		config.TemporaryStreamKey = *newStreamKey
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
	config.WebServerPort = configRepository.GetHTTPPortNumber()

	// Set the web server ip
	if *webServerIPOverride != "" {
		log.Println("Saving new web server listen IP address to", *webServerIPOverride)
		if err := configRepository.SetHTTPListenAddress(*webServerIPOverride); err != nil {
			log.Errorln(err)
		}
	}
	config.WebServerIP = configRepository.GetHTTPListenAddress()

	// Set the rtmp server port
	if *rtmpPortOverride > 0 {
		log.Println("Saving new RTMP server port number to", *rtmpPortOverride)
		if err := configRepository.SetRTMPPortNumber(float64(*rtmpPortOverride)); err != nil {
			log.Errorln(err)
		}
	}

	// Set the follower validation interval
	if *followerValidationIntervalSecs > 0 {
		config.FollowerValidationInterval = time.Duration(*followerValidationIntervalSecs) * time.Second
		log.Printf("Follower validation interval set to %v", config.FollowerValidationInterval)
	}
}

func configureLogging(enableDebugFeatures bool, enableVerboseLogging bool) {
	logging.Setup(enableDebugFeatures, enableVerboseLogging)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
