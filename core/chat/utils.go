package chat

import "syscall"

// Set the soft file handler limit and return 80% of
// the max as the max client connection limit.
func getMaxConnectionCount() uint {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	return uint(float32(rLimit.Cur) * 0.8)
}
