package data

import (
	"strings"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	datastoreValuesVersion   = 2
	datastoreValueVersionKey = "DATA_STORE_VERSION"
)

func migrateDatastoreValues(datastore *Datastore) {
	currentVersion, _ := _datastore.GetNumber(datastoreValueVersionKey)
	if currentVersion == 0 {
		currentVersion = datastoreValuesVersion
	}

	for v := currentVersion; v < datastoreValuesVersion; v++ {
		log.Infof("Migration datastore values from %d to %d\n", int(v), int(v+1))
		switch v {
		case 0:
			migrateToDatastoreValues1(datastore)
		case 1:
			migrateToDatastoreValues2(datastore)
		default:
			log.Fatalln("missing datastore values migration step")
		}
	}
	if err := _datastore.SetNumber(datastoreValueVersionKey, datastoreValuesVersion); err != nil {
		log.Errorln("error setting datastore value version:", err)
	}
}

func migrateToDatastoreValues1(datastore *Datastore) {
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

func migrateToDatastoreValues2(datastore *Datastore) {
	oldAdminPassword, _ := datastore.GetString("stream_key")
	oldAdminPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(oldAdminPassword), bcrypt.DefaultCost)
	_ = SetAdminPasswordHashed(oldAdminPasswordHash)
	_ = SetStreamKeysHashed([]models.StreamKeyHashed{
		{Key: oldAdminPasswordHash, Comment: "Default stream key"},
	})
}
