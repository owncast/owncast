package main

import (
	"github.com/owncast/owncast/cmd/service-arbiter-example/services"
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/a"
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/b"
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/c"
)

func New() *App {
	app := &App{}

	app.Services.A = a.New(app.ServiceArbitration)
	app.Services.B = b.New(app.ServiceArbitration)
	app.Services.C = c.New(app.ServiceArbitration)

	app.ServiceArbitration.Services.A = app.Services.A
	app.ServiceArbitration.Services.B = app.Services.B
	app.ServiceArbitration.Services.C = app.Services.C

	return app
}

type App struct {
	services.ServiceArbitration
	Services
}

type Services struct {
	A *a.Service
	B *b.Service
	C *c.Service
}
