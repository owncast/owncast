package core

import (
	"fmt"
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

//GetVersion gets the version string
func GetVersion() string {
	return fmt.Sprintf("Owncast v%s-%s (%s)", BuildVersion, BuildType, GitCommit)
}
