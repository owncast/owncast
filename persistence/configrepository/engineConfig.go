package configrepository

import "github.com/owncast/owncast/models"

// Compile-time guarantee that the production config repository satisfies the
// narrow models.EngineConfig surface the streaming engine packages (rtmp,
// transcoder, storage) consume. This structural coupling is what lets those
// packages run without the database-backed repository. If one of the nine
// engine getters is renamed, removed, or its signature drifts, the build fails
// here with a clear message instead of deep inside service wiring.
var (
	_ models.EngineConfig = ConfigRepository(nil)
	_ models.EngineConfig = (*SqlConfigRepository)(nil)
)
