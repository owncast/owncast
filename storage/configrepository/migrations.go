package configrepository

import (
	"strings"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/storage/data"

	log "github.com/sirupsen/logrus"
)

const (
	datastoreValuesVersion   = 3
	datastoreValueVersionKey = "DATA_STORE_VERSION"
)

func migrateDatastoreValues(ds *data.Store, cr *SqlConfigRepository) {
	currentVersion, _ := ds.GetNumber(datastoreValueVersionKey)
	if currentVersion == 0 {
		currentVersion = datastoreValuesVersion
	}

	for v := currentVersion; v < datastoreValuesVersion; v++ {
		log.Infof("Migration datastore values from %d to %d\n", int(v), int(v+1))
		switch v {
		case 0:
			migrateToDatastoreValues1(ds, cr)
		case 1:
			migrateToDatastoreValues2(ds, cr)
		case 2:
			migrateToDatastoreValues3ServingEndpoint3(ds, cr)
		default:
			log.Fatalln("missing datastore values migration step")
		}
	}
	if err := ds.SetNumber(datastoreValueVersionKey, datastoreValuesVersion); err != nil {
		log.Errorln("error setting datastore value version:", err)
	}
}

func migrateToDatastoreValues1(ds *data.Store, cr *SqlConfigRepository) {
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

func migrateToDatastoreValues2(ds *data.Store, cr *SqlConfigRepository) {
	oldAdminPassword, _ := ds.GetString("stream_key")
	_ = cr.SetAdminPassword(oldAdminPassword)
	_ = cr.SetStreamKeys([]models.StreamKey{
		{Key: oldAdminPassword, Comment: "Default stream key"},
	})
}

func migrateToDatastoreValues3ServingEndpoint3(_ *data.Store, cr *SqlConfigRepository) {
	s3Config := cr.GetS3Config()

	if !s3Config.Enabled {
		return
	}

	_ = cr.SetVideoServingEndpoint(s3Config.ServingEndpoint)
}
