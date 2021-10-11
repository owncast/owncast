package main

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.Warnln("Explicitly building with pkged.go is no longer required and the file will be removed in the future. See https://github.com/owncast/owncast/pull/1464.")
}
