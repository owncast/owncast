package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/core/rtmp"
	"github.com/owncast/owncast/core/webhooks"
	"github.com/owncast/owncast/internal/activitypub"
	"github.com/owncast/owncast/internal/activitypub/apmodels"
	"github.com/owncast/owncast/internal/activitypub/crypto"
	"github.com/owncast/owncast/internal/activitypub/follower"
	"github.com/owncast/owncast/internal/activitypub/outbox"
	"github.com/owncast/owncast/internal/activitypub/persistence"
	"github.com/owncast/owncast/internal/activitypub/workerpool"
	"github.com/owncast/owncast/metrics"
	"github.com/owncast/owncast/notifications"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/yp"
)

const (
	permDataDir = 0o700
)

func New(pathDataDir string) (app *App, err error) {
	// Create the data directory if needed
	if !utils.DoesFileExists(pathDataDir) {
		if err := os.Mkdir(pathDataDir, permDataDir); err != nil {
			return nil, fmt.Errorf("creating data directory: %v", err)
		}
	}

	// Migrate old (pre 0.1.0) emoji to new location if they exist.
	utils.MigrateCustomEmojiLocations()

	// Otherwise save the default emoji to the data directory.
	if err := data.SetupEmojiDirectory(); err != nil {
		return nil, fmt.Errorf("setting up emoji directory: %v", err)
	}

	// Recreate the temp dir
	if utils.DoesFileExists(config.TempDir) {
		if err := os.RemoveAll(config.TempDir); err != nil {
			return nil, fmt.Errorf("removing temp dir '%s': %v\nCheck permissions!", config.TempDir, err)
		}
	}

	if err := os.Mkdir(config.TempDir, 0o700); err != nil {
		return nil, fmt.Errorf("creating temp dir: %v", err)
	}

	app = &App{}

	if app.Data, err = data.New(config.DatabaseFilePath); err != nil {
		return nil, fmt.Errorf("initializing app data service: %v", err)
	}

	if app.Persistence, err = persistence.New(app.Data); err != nil {
		return nil, fmt.Errorf("initializing app persistence service: %v", err)
	}

	if app.ActivityPub, err = activitypub.New(app.Persistence); err != nil {
		return nil, fmt.Errorf("initializing app activity pub service: %v", err)
	}

	if app.Webhooks, err = webhooks.New(app.Data); err != nil {
		return nil, fmt.Errorf("initializing app webhooks service: %v", err)
	}

	if app.Notifier, err = notifications.New(app.Data); err != nil {
		return nil, fmt.Errorf("initializing app notifier: %v", err)
	}

	if app.Chat, err = chat.New(app.Data, app.Webhooks); err != nil {
		return nil, fmt.Errorf("initializing app controller service: %v", err)
	}

	if app.rtmp, err = rtmp.New(app.Data); err != nil {
		return nil, fmt.Errorf("initializing app RTMP service: %v", err)
	}

	if app.crypto, err = crypto.New(app.Data); err != nil {
		return nil, fmt.Errorf("initializing app crypto service: %v", err)
	}

	if app.models, err = apmodels.New(app.Data, app.crypto); err != nil {
		return nil, fmt.Errorf("initializing app activity pub models service: %v", err)
	}

	if app.workerpool, err = workerpool.New(); err != nil {
		return nil, fmt.Errorf("initializing app worker pool service: %v", err)
	}

	if app.follower, err = follower.New(app.Data, app.crypto, app.models, app.workerpool); err != nil {
		return nil, fmt.Errorf("initializing app follower service: %v", err)
	}

	if app.outbox, err = outbox.New(app.Persistence, app.crypto, app.models, app.follower, app.workerpool); err != nil {
		return nil, fmt.Errorf("initializing app outbox service: %v", err)
	}

	if app.Core, err = core.New(app.ActivityPub, app.Webhooks, app.YP, app.Notifier, app.Chat, app.rtmp); err != nil {
		return nil, fmt.Errorf("initializing app core service: %v", err)
	}

	if app.Metrics, err = metrics.New(app.Core); err != nil {
		return nil, fmt.Errorf("initializing app metrics service: %v", err)
	}

	if app.Controller, err = controllers.New(app.ActivityPub, app.Core, app.Metrics, app.Notifier, app.rtmp, app.follower, app.outbox); err != nil {
		return nil, fmt.Errorf("initializing app controller service: %v", err)
	}

	app.YP = app.Core.YP
	app.Core.Notifier = app.Notifier

	return app, nil
}

type App struct {
	Persistence *persistence.Service
	ActivityPub *activitypub.Service
	Data        *data.Service
	Core        *core.Service
	Chat        *chat.Service
	Metrics     *metrics.Service
	Controller  *controllers.Service
	Webhooks    *webhooks.Service
	Notifier    *notifications.Notifier
	router      *http.ServeMux
	YP          *yp.YP
	rtmp        *rtmp.Service
	follower    *follower.Service
	outbox      *outbox.Service
	crypto      *crypto.Service
	models      *apmodels.Service
	workerpool  *workerpool.Service
}

func (a *App) Serve() error {
	// starts the core
	if err := a.Core.Start(); err != nil {
		return fmt.Errorf("starting core: %v", err)
	}

	go a.Metrics.Start(a.Core.GetStatus)

	if err := a.initRoutes(); err != nil {
		return fmt.Errorf("initializing routes: %v", err)
	}

	return nil
}
