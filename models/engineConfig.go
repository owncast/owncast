package models

// EngineConfig is the read-only configuration surface consumed by the
// streaming engine packages (rtmp, transcoder, storage). It is the narrow
// seam that lets those packages run outside the core process without linking
// the database-backed configuration repository: the core's ConfigRepository
// satisfies it structurally, while a remote engine satisfies it from a pulled
// config snapshot.
//
// Keep this list to exactly the values the engine packages read. Widening it
// re-couples the engine to core internals.
type EngineConfig interface {
	GetRTMPPortNumber() int
	GetRTMPBindAddress() string
	GetFfMpegPath() string
	GetVideoCodec() string
	GetS3Config() S3
	GetStreamLatencyLevel() LatencyLevel
	GetStreamOutputVariants() []StreamOutputVariant
	GetVideoServingEndpoint() string
	GetStreamKeys() []StreamKey
}
