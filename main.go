package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router"
)

func main() {
	// logrus.SetReportCaller(true)
	log.Println(core.GetVersion())

	//TODO: potentially load the config from a flag like:
	//configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")
	// flag.Parse()

	if err := config.Load("config.yaml"); err != nil {
		panic(err)
	}

	// starts the core
	if err := core.Start(); err != nil {
		log.Println("failed to start the core package")
		panic(err)
	}

	if err := router.Start(); err != nil {
		log.Println("failed to start/run the router")
		panic(err)
	}
}
