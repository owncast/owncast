package migrations

import (
	"strings"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/storage/datastore"

	log "github.com/sirupsen/logrus"
)

const (
	datastoreValuesVersion   = 3
	datastoreValueVersionKey = "DATA_STORE_VERSION"
)

func migrateDatastoreValues(ds *datastore.Datastore, cr *configrepository.ConfigRepository) {
	currentVersion, _ := ds.GetNumber(datastoreValueVersionKey)
	if currentVersion == 0 {
		currentVersion = datastoreValuesVersion
	}

	for v := currentVersion; v < datastoreValuesVersion; v++ {
		log.Infof("Migration datastore values from %d to %d\n", int(v), int(v+1))
		switch v {
		case 0:
			migrateToDatastoreValues1(ds)
		case 1:
			migrateToDatastoreValues2(ds)
		case 2:
			migrateToDatastoreValues3ServingEndpoint3(ds)
		default:
			log.Fatalln("missing datastore values migration step")
		}
	}
	if err := ds.SetNumber(datastoreValueVersionKey, datastoreValuesVersion); err != nil {
		log.Errorln("error setting datastore value version:", err)
	}
}

func migrateToDatastoreValues1(ds *datastore.Datastore, cr *configrepository.ConfigRepository) {
	// Migrate the forbidden usernames to be a slice instead of a string.

	forbiddenUsernamesString, _ := ds.GetString(blockedUsernamesKey)
	if forbiddenUsernamesString != "" {
		forbiddenUsernamesSlice := strings.Split(forbiddenUsernamesString, ",")
		if err := ds.SetStringSlice(blockedUsernamesKey, forbiddenUsernamesSlice); err != nil {
			log.Errorln("error migrating blocked username list:", err)
		}
	}

	// Migrate the suggested usernames to be a slice instead of a string.
	suggestedUsernamesString, _ := ds.GetString(suggestedUsernamesKey)
	if suggestedUsernamesString != "" {
		suggestedUsernamesSlice := strings.Split(suggestedUsernamesString, ",")
		if err := ds.SetStringSlice(suggestedUsernamesKey, suggestedUsernamesSlice); err != nil {
			log.Errorln("error migrating suggested username list:", err)
		}
	}
}

func migrateToDatastoreValues2(ds *datastore.Datastore) {
	oldAdminPassword, _ := datastore.GetString("stream_key")
	_ = SetAdminPassword(oldAdminPassword)
	_ = SetStreamKeys([]models.StreamKey{
		{Key: oldAdminPassword, Comment: "Default stream key"},
	})
}

func migrateToDatastoreValues3ServingEndpoint3(_ *Datastore) {
	s3Config := GetS3Config()

	if !s3Config.Enabled {
		return
	}

	_ = SetVideoServingEndpoint(s3Config.ServingEndpoint)
}
