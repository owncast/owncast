package data

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
)

const (
	datastoreValuesVersion   = 2
	datastoreValueVersionKey = "DATA_STORE_VERSION"
)

func (s *Service) migrateDatastoreValues(datastore *Datastore) {
	currentVersion, _ := s.Store.GetNumber(datastoreValueVersionKey)
	if currentVersion == 0 {
		currentVersion = datastoreValuesVersion
	}

	for v := currentVersion; v < datastoreValuesVersion; v++ {
		log.Infof("Migration Store values from %d to %d\n", int(v), int(v+1))
		switch v {
		case 0:
			s.migrateToDatastoreValues1(datastore)
		case 1:
			s.migrateToDatastoreValues2(datastore)
		default:
			log.Fatalln("missing Store values migration step")
		}
	}
	if err := s.Store.SetNumber(datastoreValueVersionKey, datastoreValuesVersion); err != nil {
		log.Errorln("error setting Store value version:", err)
	}
}

func (s *Service) migrateToDatastoreValues1(datastore *Datastore) {
	// Migrate the forbidden usernames to be a slice instead of a string.
	forbiddenUsernamesString, _ := datastore.GetString(blockedUsernamesKey)
	if forbiddenUsernamesString != "" {
		forbiddenUsernamesSlice := strings.Split(forbiddenUsernamesString, ",")
		if err := datastore.SetStringSlice(blockedUsernamesKey, forbiddenUsernamesSlice); err != nil {
			log.Errorln("error migrating blocked username list:", err)
		}
	}

	// Migrate the suggested usernames to be a slice instead of a string.
	suggestedUsernamesString, _ := datastore.GetString(suggestedUsernamesKey)
	if suggestedUsernamesString != "" {
		suggestedUsernamesSlice := strings.Split(suggestedUsernamesString, ",")
		if err := datastore.SetStringSlice(suggestedUsernamesKey, suggestedUsernamesSlice); err != nil {
			log.Errorln("error migrating suggested username list:", err)
		}
	}
}

func (s *Service) migrateToDatastoreValues2(datastore *Datastore) {
	oldAdminPassword, _ := datastore.GetString("stream_key")
	_ = s.SetAdminPassword(oldAdminPassword)
	_ = s.SetStreamKeys([]models.StreamKey{
		{Key: oldAdminPassword, Comment: "Default stream key"},
	})
}
