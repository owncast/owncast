package services

import (
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/a"
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/b"
	"github.com/owncast/owncast/cmd/service-arbiter-example/services/c"
)

type ServiceArbitration struct {
	Services struct {
		A *a.Service
		B *b.Service
		C *c.Service
	}
}

func (a ServiceArbitration) CallServiceA() {
	a.Services.A.Call()
}

func (a ServiceArbitration) CallServiceB() {
	a.Services.B.Call()
}

func (a ServiceArbitration) CallServiceC() {
	a.Services.C.Call()
}
