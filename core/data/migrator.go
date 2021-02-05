package data

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// RunMigrations will start the migration process from the config file.
func RunMigrations() {
	if !utils.DoesFileExists(config.BackupDirectory) {
		if err := os.Mkdir(config.BackupDirectory, 0700); err != nil {
			log.Errorln("Unable to create backup directory", err)
			return
		}
	}

	migrateConfigFile()
	migrateStatsFile()
	migrateYPKey()
}

func migrateStatsFile() {
	oldStats := models.Stats{
		Clients: make(map[string]models.Client),
	}

	if !utils.DoesFileExists(config.StatsFile) {
		return
	}

	log.Infoln("Migrating", config.StatsFile, "to new datastore")

	jsonFile, err := ioutil.ReadFile(config.StatsFile)
	if err != nil {
		log.Errorln(err)
		return
	}

	if err := json.Unmarshal(jsonFile, &oldStats); err != nil {
		log.Errorln(err)
		return
	}

	_ = SetPeakSessionViewerCount(oldStats.SessionMaxViewerCount)
	_ = SetPeakOverallViewerCount(oldStats.OverallMaxViewerCount)

	if err := utils.Move(config.StatsFile, "backup/stats.old"); err != nil {
		log.Warnln(err)
	}
}

func migrateYPKey() {
	filePath := ".yp.key"

	if !utils.DoesFileExists(filePath) {
		return
	}

	log.Infoln("Migrating", filePath, "to new datastore")

	keyFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorln("Unable to migrate", keyFile, "there may be issues registering with the directory")
	}

	if err := SetDirectoryRegistrationKey(string(keyFile)); err != nil {
		log.Errorln("Unable to migrate", keyFile, "there may be issues registering with the directory")
		return
	}

	if err := utils.Move(filePath, "backup/yp.key.old"); err != nil {
		log.Warnln(err)
	}
}

func migrateConfigFile() {
	filePath := "config.yaml"

	if !utils.DoesFileExists(filePath) {
		return
	}

	log.Infoln("Migrating", filePath, "to new datastore")

	var oldConfig configFile

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorln("config file err", err)
		return
	}

	if err := yaml.Unmarshal(yamlFile, &oldConfig); err != nil {
		log.Errorln("Error reading the config file.", err)
		return
	}

	_ = SetServerName(oldConfig.InstanceDetails.Name)
	_ = SetServerSummary(oldConfig.InstanceDetails.Summary)
	_ = SetServerMetadataTags(oldConfig.InstanceDetails.Tags)
	_ = SetStreamKey(oldConfig.VideoSettings.StreamingKey)
	_ = SetLogoPath(oldConfig.InstanceDetails.Logo)
	_ = SetNSFW(oldConfig.InstanceDetails.NSFW)
	_ = SetServerURL(oldConfig.YP.InstanceURL)
	_ = SetDirectoryEnabled(oldConfig.YP.Enabled)
	_ = SetSocialHandles(oldConfig.InstanceDetails.SocialHandles)
	_ = SetFfmpegPath(oldConfig.FFMpegPath)
	_ = SetHTTPPortNumber(float64(oldConfig.WebServerPort))
	_ = SetRTMPPortNumber(float64(oldConfig.RTMPServerPort))

	// Migrate video variants
	variants := []models.StreamOutputVariant{}
	for _, variant := range oldConfig.VideoSettings.StreamQualities {
		migratedVariant := models.StreamOutputVariant{}
		migratedVariant.IsAudioPassthrough = true
		migratedVariant.IsVideoPassthrough = variant.IsVideoPassthrough
		migratedVariant.Framerate = variant.Framerate
		migratedVariant.VideoBitrate = variant.VideoBitrate
		migratedVariant.ScaledHeight = variant.ScaledHeight
		migratedVariant.ScaledWidth = variant.ScaledWidth

		presetMapping := map[string]int{
			"ultrafast": 1,
			"superfast": 2,
			"veryfast":  3,
			"faster":    4,
			"fast":      5,
		}
		migratedVariant.CpuUsageLevel = presetMapping[variant.EncoderPreset]
		variants = append(variants, migratedVariant)
	}
	_ = SetStreamOutputVariants(variants)

	// Migrate latency level
	level := 4
	oldSegmentLength := oldConfig.VideoSettings.ChunkLengthInSeconds
	oldNumberOfSegments := oldConfig.Files.MaxNumberInPlaylist
	latencyLevels := models.GetLatencyConfigs()

	if oldSegmentLength == latencyLevels[1].SecondsPerSegment && oldNumberOfSegments == latencyLevels[1].SegmentCount {
		level = 1
	} else if oldSegmentLength == latencyLevels[2].SecondsPerSegment && oldNumberOfSegments == latencyLevels[2].SegmentCount {
		level = 2
	} else if oldSegmentLength == latencyLevels[3].SecondsPerSegment && oldNumberOfSegments == latencyLevels[3].SegmentCount {
		level = 3
	} else if oldSegmentLength == latencyLevels[5].SecondsPerSegment && oldNumberOfSegments == latencyLevels[5].SegmentCount {
		level = 5
	} else if oldSegmentLength >= latencyLevels[6].SecondsPerSegment && oldNumberOfSegments >= latencyLevels[6].SegmentCount {
		level = 6
	}

	_ = SetStreamLatencyLevel(float64(level))

	// Migrate storage config
	_ = SetS3Config(models.S3(oldConfig.Storage))

	// Migrate the old content.md file
	content, err := ioutil.ReadFile(config.ExtraInfoFile)
	if err == nil && len(content) > 0 {
		_ = SetExtraPageBodyContent(string(content))
	}

	if err := utils.Move(filePath, "backup/config.old"); err != nil {
		log.Warnln(err)
	}

	log.Infoln("Your old config file can be found in the backup directory for reference. For all future configuration use the web admin.")
}

type configFile struct {
	DatabaseFilePath    string `yaml:"databaseFile"`
	EnableDebugFeatures bool   `yaml:"-"`
	FFMpegPath          string
	Files               files           `yaml:"files"`
	InstanceDetails     instanceDetails `yaml:"instanceDetails"`
	VersionInfo         string          `yaml:"-"` // For storing the version/build number
	VersionNumber       string          `yaml:"-"`
	VideoSettings       videoSettings   `yaml:"videoSettings"`
	WebServerPort       int
	RTMPServerPort      int
	YP                  yp `yaml:"yp"`
	Storage             s3 `yaml:"s3"`
}

// instanceDetails defines the user-visible information about this particular instance.
type instanceDetails struct {
	Name             string                `yaml:"name"`
	Title            string                `yaml:"title"`
	Summary          string                `yaml:"summary"`
	Logo             string                `yaml:"logo"`
	Tags             []string              `yaml:"tags"`
	Version          string                `yaml:"version"`
	NSFW             bool                  `yaml:"nsfw"`
	ExtraPageContent string                `yaml:"extraPageContent"`
	StreamTitle      string                `yaml:"streamTitle"`
	SocialHandles    []models.SocialHandle `yaml:"socialHandles"`
}

type videoSettings struct {
	ChunkLengthInSeconds      int             `yaml:"chunkLengthInSeconds"`
	StreamingKey              string          `yaml:"streamingKey"`
	StreamQualities           []streamQuality `yaml:"streamQualities"`
	HighestQualityStreamIndex int             `yaml:"-"`
}

// yp allows registration to the central Owncast yp (Yellow pages) service operating as a directory.
type yp struct {
	Enabled     bool   `yaml:"enabled"`
	InstanceURL string `yaml:"instanceUrl"` // The public URL the directory should link to
}

// streamQuality defines the specifics of a single HLS stream variant.
type streamQuality struct {
	// Enable passthrough to copy the video and/or audio directly from the
	// incoming stream and disable any transcoding.  It will ignore any of
	// the below settings.
	IsVideoPassthrough bool `yaml:"videoPassthrough" json:"videoPassthrough"`
	IsAudioPassthrough bool `yaml:"audioPassthrough" json:"audioPassthrough"`

	VideoBitrate int `yaml:"videoBitrate" json:"videoBitrate"`
	AudioBitrate int `yaml:"audioBitrate" json:"audioBitrate"`

	// Set only one of these in order to keep your current aspect ratio.
	// Or set neither to not scale the video.
	ScaledWidth  int `yaml:"scaledWidth" json:"scaledWidth,omitempty"`
	ScaledHeight int `yaml:"scaledHeight" json:"scaledHeight,omitempty"`

	Framerate     int    `yaml:"framerate" json:"framerate"`
	EncoderPreset string `yaml:"encoderPreset" json:"encoderPreset"`
}

type files struct {
	MaxNumberInPlaylist int `yaml:"maxNumberInPlaylist"`
}

// s3 is for configuring the s3 integration.
type s3 struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	Endpoint        string `yaml:"endpoint" json:"endpoint,omitempty"`
	ServingEndpoint string `yaml:"servingEndpoint" json:"servingEndpoint,omitempty"`
	AccessKey       string `yaml:"accessKey" json:"accessKey,omitempty"`
	Secret          string `yaml:"secret" json:"secret,omitempty"`
	Bucket          string `yaml:"bucket" json:"bucket,omitempty"`
	Region          string `yaml:"region" json:"region,omitempty"`
	ACL             string `yaml:"acl" json:"acl,omitempty"`
}
