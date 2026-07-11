package stream

import (
	"os"
	"testing"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

func newStreamStatsTestRepo(t *testing.T) configrepository.ConfigRepository {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("failed to set up datastore: %v", err)
	}
	return configrepository.New(ds)
}

func TestSaveStatsPersistsRecentLastLiveHeartbeatWhileOnline(t *testing.T) {
	repo := newStreamStatsTestRepo(t)
	connectedAt := time.Now().Add(-3 * time.Hour)
	s := &Service{
		configRepository: repo,
		stats: &models.Stats{
			StreamConnected: true,
			LastConnectTime: &utils.NullTime{Time: connectedAt, Valid: true},
		},
	}

	beforeSave := time.Now().Add(-2 * time.Second)
	s.saveStats()

	got, err := repo.GetLastDisconnectTime()
	if err != nil {
		t.Fatalf("GetLastDisconnectTime: %v", err)
	}
	if got == nil || !got.Valid {
		t.Fatal("expected an online heartbeat timestamp to be persisted")
	}
	if got.Time.Before(beforeSave) {
		t.Fatalf("persisted timestamp %v is older than save window start %v", got.Time, beforeSave)
	}
	if got.Time.Before(connectedAt) {
		t.Fatalf("persisted timestamp %v should not predate the live session start %v", got.Time, connectedAt)
	}

	restored := (&Service{configRepository: repo}).getSavedStats()
	if restored.LastDisconnectTime == nil || !restored.LastDisconnectTime.Valid {
		t.Fatal("expected restored stats to include the persisted last-live timestamp")
	}
	if restored.LastDisconnectTime.Time.Before(beforeSave) {
		t.Fatalf("restored timestamp %v is older than save window start %v", restored.LastDisconnectTime.Time, beforeSave)
	}
}

func TestSaveStatsPersistsExplicitDisconnectTimeWhenOffline(t *testing.T) {
	repo := newStreamStatsTestRepo(t)
	disconnectedAt := time.Now().Add(-45 * time.Minute).Round(time.Second)
	s := &Service{
		configRepository: repo,
		stats: &models.Stats{
			StreamConnected:    false,
			LastDisconnectTime: &utils.NullTime{Time: disconnectedAt, Valid: true},
		},
	}

	s.saveStats()

	got, err := repo.GetLastDisconnectTime()
	if err != nil {
		t.Fatalf("GetLastDisconnectTime: %v", err)
	}
	if got == nil || !got.Valid {
		t.Fatal("expected a disconnect timestamp to be persisted")
	}
	if !got.Time.Equal(disconnectedAt) {
		t.Fatalf("persisted disconnect time = %v, want %v", got.Time, disconnectedAt)
	}
}
