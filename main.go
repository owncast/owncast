package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/router"
)

// the following are injected at build-time
var (
	//GitCommit is the commit which this version of owncast is running
	GitCommit = "unknown"
	//BuildVersion is the version
	BuildVersion = "0.0.0"
	//BuildType is the type of build
	BuildType = "localdev"
)

func main() {
	// logrus.SetReportCaller(true)
	log.Println(getVersion())

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

//getVersion gets the version string
func getVersion() string {
	return fmt.Sprintf("Owncast v%s-%s (%s)", BuildVersion, BuildType, GitCommit)
}
