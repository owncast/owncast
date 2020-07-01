package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
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
	log.Println(getVersion())

	configFile := flag.String("configFile", "config.yaml", "Config File full path. Defaults to current folder")
	enableDebugOptions := flag.Bool("enableDebugFeatures", false, "Enable additional debugging options.")

	flag.Parse()

	if *enableDebugOptions {
		logrus.SetReportCaller(true)
	}

	if err := config.Load(*configFile, getVersion()); err != nil {
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
