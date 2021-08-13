// +build !freebsd

package chat

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

func setSystemConcurrentConnectionLimit(limit int64) {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	originalLimit := rLimit.Cur
	rLimit.Cur = uint64(limit)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Traceln("Max process connection count changed from system limit of", originalLimit, "to", limit)
}
