package config

// These are runtime-set values used for configuration.

// DatabaseFilePath is the path to the file ot be used as the global database for this run of the application.
var DatabaseFilePath = "data/owncast.db"

// EnableDebugFeatures will print additional data to help in debugging.
var EnableDebugFeatures = false

// VersionInfo  is a string for displaying the version with a v prefix.
var VersionInfo = "v" + CurrentBuildString

// VersionNumber is the current version string.
var VersionNumber = CurrentBuildString

// WebServerPort is the port for Owncast's webserver that is used for this execution of the service.
var WebServerPort = 8080

// InternalHLSListenerPort is the port for HLS writes that is used for this execution of the service.
var InternalHLSListenerPort = "8927"

// ConfigFilePath is the path to the config file for migration.
var ConfigFilePath = "config.yaml"
