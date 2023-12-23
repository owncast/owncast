package data

import (
	"strings"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

const (
	datastoreValuesVersion   = 4
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
		case 2:
			migrateToDatastoreValues3ServingEndpoint3(datastore)
		case 3:
			migrateToDatastoreVaapiCodecSettingValue4(datastore)
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

func migrateToDatastoreVaapiCodecSettingValue4(_ *Datastore) {
	// If the currently selected codec is "vaapi" then we need
	// to migrate the name to the updated name.
	currentCodec := GetVideoCodec()
	if currentCodec != "h264_vaapi" {
		return
	}

	// The updated name for the old vaapi codec is "h264_vaapi_legacy"
	// so we update it. This is assuming existing users will be using older
	// versions of ffmpeg.
	_ = SetVideoCodec("h264_vaapi_legacy")

	log.Println("An update to the vaapi video codec has been made. It will now be selected as vaapi (Legacy) in your video settings. However, if you are running a version of ffmpeg greater than 5.0 you should update your video codec settings to use the 5.0+ codec setting instead.")
}
