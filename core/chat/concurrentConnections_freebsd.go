//go:build freebsd
// +build freebsd

package chat

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

func setSystemConcurrentConnectionLimit(limit uint64) {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalln(err)
	}

	originalLimit := rLimit.Cur
	rLimit.Cur = int64(limit)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		log.Fatalln(err)
	}

	log.Traceln("Max process connection count changed from system limit of", originalLimit, "to", limit)
}
