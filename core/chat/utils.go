package chat

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Set the soft file handler limit as 70% of
// the max as the client connection limit.
func handleMaxConnectionCount() uint {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	originalLimit := rLimit.Cur
	// Set the limit to 70% of max so the machine doesn't die even if it's maxed out for some reason.
	proposedLimit := int(float32(rLimit.Max) * 0.7)

	rLimit.Cur = uint64(proposedLimit)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Traceln("Max process connection count increased from", originalLimit, "to", proposedLimit)

	return uint(float32(rLimit.Cur))
}
