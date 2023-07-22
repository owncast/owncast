// nolint:goimports
//go:build !freebsd && !windows
// +build !freebsd,!windows

package cmd

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

func setSystemConcurrentConnectionLimit(limit int64) {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalln(err)
	}

	originalLimit := rLimit.Cur
	rLimit.Cur = uint64(limit)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalln(err)
	}

	log.Traceln("Max process connection count changed from system limit of", originalLimit, "to", limit)
}

func getMaximumConcurrentConnectionLimit() int64 {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalln(err)
	}

	// Return the limit to 70% of max so the machine doesn't die even if it's maxed out for some reason.
	proposedLimit := int64(float32(rLimit.Max) * 0.7)

	return proposedLimit
}
