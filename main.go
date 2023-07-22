package main

import "github.com/owncast/owncast/cmd"

func main() {
	app := &cmd.Application{}
	app.Start()
}

// var configRepository = configrepository.Get()

// // nolint:cyclop
// func main() {
// 	flag.Parse()

// 	config = configservice.NewConfig()

// 	// Otherwise save the default emoji to the data directory.
// 	if err := data.SetupEmojiDirectory(); err != nil {
// 		log.Fatalln("Cannot set up emoji directory", err)
// 	}

// 	if err := data.SetupPersistence(config.DatabaseFilePath); err != nil {
// 		log.Fatalln("failed to open database", err)
// 	}

// 	handleCommandLineFlags()

// 	// starts the core
// 	if err := core.Start(); err != nil {
// 		log.Fatalln("failed to start the core package", err)
// 	}

// 	go metrics.Start(core.GetStatus)

// 	webserver := webserver.New()
// 	if err := webserver.Start(config.WebServerIP, config.WebServerPort); err != nil {
// 		log.Fatalln("failed to start/run the web server", err)
// 	}
// }

// func handleCommandLineFlags() {
// }
