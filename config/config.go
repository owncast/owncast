package config

// Runtime settings that are not stored to disk
// or provided via cli flags.

// DatabaseFilePath is the path to the file ot be used as the global database for this run of the application.
var DatabaseFilePath = "data/owncast.db"

// EnableDebugFeatures will print additional data to help in debugging.
var EnableDebugFeatures = false

//
var VersionInfo = "v" + CurrentBuildString
var VersionNumber = CurrentBuildString

var WebServerPortOverride = 8080
var RTMPServerPortOverride = 1935

var HighestQualityStreamIndex = 0
