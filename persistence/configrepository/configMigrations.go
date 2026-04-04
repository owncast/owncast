package configrepository

import (
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/webserver/handlers/generated"
	log "github.com/sirupsen/logrus"
)

const (
	datastoreValuesVersion   = 5
	datastoreValueVersionKey = "DATA_STORE_VERSION"
)

func migrateDatastoreValues(datastore *data.Datastore, configRepository ConfigRepository) {
	currentVersion, _ := datastore.GetNumber(datastoreValueVersionKey)
	if currentVersion == 0 {
		currentVersion = datastoreValuesVersion
	}

	for v := currentVersion; v < datastoreValuesVersion; v++ {
		log.Infof("Migration datastore values from %d to %d\n", int(v), int(v+1))
		switch v {
		case 0:
			migrateToDatastoreValues1(datastore)
		case 1:
			migrateToDatastoreValues2(datastore, configRepository)
		case 2:
			migrateToDatastoreValues3ServingEndpoint3(configRepository)
		case 3:
			migrateToDatastoreValues4(datastore, configRepository)
		case 4:
			migrateToDatastoreValues5(datastore, configRepository)
		default:
			log.Fatalln("missing datastore values migration step")
		}
	}
	if err := datastore.SetNumber(datastoreValueVersionKey, datastoreValuesVersion); err != nil {
		log.Errorln("error setting datastore value version:", err)
	}
}

func migrateToDatastoreValues1(datastore *data.Datastore) {
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

func migrateToDatastoreValues2(datastore *data.Datastore, configRepository ConfigRepository) {
	oldAdminPassword, _ := datastore.GetString("stream_key")
	// Avoids double hashing the password
	_ = datastore.SetString("admin_password_key", oldAdminPassword)
	comment := "Default stream key"
	_ = configRepository.SetStreamKeys([]generated.StreamKey{
		{Key: &oldAdminPassword, Comment: &comment},
	})
}

func migrateToDatastoreValues3ServingEndpoint3(configRepository ConfigRepository) {
	s3Config := configRepository.GetS3Config()

	if !s3Config.Enabled {
		return
	}

	_ = configRepository.SetVideoServingEndpoint(s3Config.ServingEndpoint)
}

func migrateToDatastoreValues4(datastore *data.Datastore, configRepository ConfigRepository) {
	unhashed_pass, _ := datastore.GetString("admin_password_key")
	err := configRepository.SetAdminPassword(unhashed_pass)
	if err != nil {
		log.Fatalln("error migrating admin password:", err)
	}
}

func migrateToDatastoreValues5(datastore *data.Datastore, configRepository ConfigRepository) {
	// Migrate the old video_codec value (which was really an encoder name like "libx264",
	// "h264_nvenc") to the new video_encoder key with the encoder type ("software", "nvenc").
	// Inlined from transcoder.OldCodecNameToEncoderType to avoid an import cycle
	// (core/transcoder imports persistence/configrepository).
	oldCodecToEncoderType := map[string]string{
		"libx264":           "software",
		"h264_nvenc":        "nvenc",
		"h264_vaapi":        "vaapi",
		"h264_qsv":          "qsv",
		"h264_omx":          "omx",
		"h264_v4l2m2m":      "v4l2m2m",
		"h264_videotoolbox": "videotoolbox",
	}

	oldCodecValue, _ := datastore.GetString("video_codec")
	if oldCodecValue != "" {
		encoderType, ok := oldCodecToEncoderType[oldCodecValue]
		if !ok {
			encoderType = "software"
		}
		if err := configRepository.SetVideoEncoder(encoderType); err != nil {
			log.Errorln("error migrating video encoder setting:", err)
		}
	}

	// Add default video and audio codec values to existing stream output variants.
	variants := configRepository.GetStreamOutputVariants()
	for i := range variants {
		if variants[i].VideoCodec == "" {
			variants[i].VideoCodec = "h264"
		}
		if variants[i].AudioCodec == "" {
			variants[i].AudioCodec = "aac"
		}
	}
	if err := configRepository.SetStreamOutputVariants(variants); err != nil {
		log.Errorln("error migrating stream output variants with codec defaults:", err)
	}
}
