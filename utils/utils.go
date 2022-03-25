package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"mvdan.cc/xurls"
)

// DoesFileExists checks if the file exists.
func DoesFileExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Errorln(err)
		return false
	}
}

// GetRelativePathFromAbsolutePath gets the relative path from the provided absolute path.
func GetRelativePathFromAbsolutePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]
	file := pathComponents[len(pathComponents)-1]

	return filepath.Join(variant, file)
}

// GetIndexFromFilePath is a utility that will return the index/key/variant name in a full path.
func GetIndexFromFilePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]

	return variant
}

// Copy copies the file to destination.
func Copy(source, destination string) error {
	input, err := os.ReadFile(source) // nolint
	if err != nil {
		return err
	}

	return os.WriteFile(destination, input, 0o600)
}

// Move moves the file at source to destination.
func Move(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil {
		log.Warnln("Moving with os.Rename failed, falling back to copy and delete!", err)
		return moveFallback(source, destination)
	}
	return nil
}

// moveFallback moves a file using a copy followed by a delete, which works across file systems.
// source: https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
func moveFallback(source, destination string) error {
	inputFile, err := os.Open(source) // nolint: gosec
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destination) // nolint: gosec
	if err != nil {
		_ = inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	_ = inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

// IsUserAgentABot returns if a web client user-agent is seen as a bot.
func IsUserAgentABot(userAgent string) bool {
	if userAgent == "" {
		return false
	}

	botStrings := []string{
		"mastodon",
		"pleroma",
		"applebot",
		"whatsapp",
		"matrix",
		"synapse",
		"element",
		"rocket.chat",
		"duckduckbot",
	}

	for _, botString := range botStrings {
		if strings.Contains(strings.ToLower(userAgent), botString) {
			return true
		}
	}

	ua := user_agent.New(userAgent)
	return ua.Bot()
}

// IsUserAgentAPlayer returns if a web client user-agent is seen as a media player.
func IsUserAgentAPlayer(userAgent string) bool {
	if userAgent == "" {
		return false
	}

	playerStrings := []string{
		"mpv",
		"player",
		"vlc",
		"applecoremedia",
	}

	for _, playerString := range playerStrings {
		if strings.Contains(strings.ToLower(userAgent), playerString) {
			return true
		}
	}

	return false
}

// RenderSimpleMarkdown will return HTML without sanitization or specific formatting rules.
func RenderSimpleMarkdown(raw string) string {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
				extension.WithLinkifyURLRegexp(
					xurls.Strict,
				),
			),
		),
	)

	trimmed := strings.TrimSpace(raw)
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(trimmed), &buf); err != nil {
		log.Fatalln(err)
	}

	return buf.String()
}

// RenderPageContentMarkdown will return HTML specifically handled for the user-specified page content.
func RenderPageContentMarkdown(raw string) string {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
				extension.WithLinkifyURLRegexp(
					xurls.Strict,
				),
			),
		),
	)

	trimmed := strings.TrimSpace(raw)
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(trimmed), &buf); err != nil {
		log.Fatalln(err)
	}

	return strings.TrimSpace(buf.String())
}

// GetCacheDurationSecondsForPath will return the number of seconds to cache an item.
func GetCacheDurationSecondsForPath(filePath string) int {
	filename := path.Base(filePath)
	fileExtension := path.Ext(filePath)

	if filename == "thumbnail.jpg" || filename == "preview.gif" {
		// Thumbnails & preview gif re-generate during live
		return 20
	} else if fileExtension == ".js" || fileExtension == ".css" {
		// Cache javascript & CSS
		return 60 * 60 * 3
	} else if fileExtension == ".ts" {
		// Cache video segments as long as you want. They can't change.
		// This matters most for local hosting of segments for recordings
		// and not for live or 3rd party storage.
		return 31557600
	} else if fileExtension == ".m3u8" {
		return 0
	} else if fileExtension == ".jpg" || fileExtension == ".png" || fileExtension == ".gif" || fileExtension == ".svg" {
		return 60 * 60 * 24 * 7
	}

	// Default cache length in seconds
	return 60 * 60 * 2
}

// IsValidURL will return if a URL string is a valid URL or not.
func IsValidURL(urlToTest string) bool {
	if _, err := url.ParseRequestURI(urlToTest); err != nil {
		return false
	}

	u, err := url.Parse(urlToTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// ValidatedFfmpegPath will take a proposed path to ffmpeg and return a validated path.
func ValidatedFfmpegPath(ffmpegPath string) string {
	if ffmpegPath != "" {
		if err := VerifyFFMpegPath(ffmpegPath); err == nil {
			return ffmpegPath
		}
		log.Warnln(ffmpegPath, "is an invalid path to ffmpeg will try to use a copy in your path, if possible")
	}

	// First look to see if ffmpeg is in the current working directory
	localCopy := "./ffmpeg"
	hasLocalCopyError := VerifyFFMpegPath(localCopy)
	if hasLocalCopyError == nil {
		// No error, so all is good.  Use the local copy.
		return localCopy
	}

	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Unable to determine path to ffmpeg. Please specify it in the admin or place a copy in the Owncast directory.")
	}

	path := strings.TrimSpace(string(out))
	return path
}

// VerifyFFMpegPath verifies that the path exists, is a file, and is executable.
func VerifyFFMpegPath(path string) error {
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		return errors.New("ffmpeg path does not exist")
	}

	if err != nil {
		return fmt.Errorf("error while verifying the ffmpeg path: %s", err.Error())
	}

	if stat.IsDir() {
		return errors.New("ffmpeg path can not be a folder")
	}

	mode := stat.Mode()
	// source: https://stackoverflow.com/a/60128480
	if mode&0o111 == 0 {
		return errors.New("ffmpeg path is not executable")
	}

	return nil
}

// CleanupDirectory removes the directory and makes it fresh again. Throws fatal error on failure.
func CleanupDirectory(path string) {
	log.Traceln("Cleaning", path)
	if err := os.RemoveAll(path); err != nil {
		log.Fatalln("Unable to remove directory. Please check the ownership and permissions", err)
	}
	if err := os.MkdirAll(path, 0o750); err != nil {
		log.Fatalln("Unable to create directory. Please check the ownership and permissions", err)
	}
}

// FindInSlice will return if a string is in a slice, and the index of that string.
func FindInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// StringSliceToMap is a convenience function to convert a slice of strings into
// a map using the string as the key.
func StringSliceToMap(stringSlice []string) map[string]interface{} {
	stringMap := map[string]interface{}{}

	for _, str := range stringSlice {
		stringMap[str] = true
	}

	return stringMap
}

// Float64MapToSlice is a convenience function to convert a map of floats into.
func Float64MapToSlice(float64Map map[string]float64) []float64 {
	float64Slice := []float64{}

	for _, val := range float64Map {
		float64Slice = append(float64Slice, val)
	}

	return float64Slice
}

// StringMapKeys returns a slice of string keys from a map.
func StringMapKeys(stringMap map[string]interface{}) []string {
	stringSlice := []string{}
	for k := range stringMap {
		stringSlice = append(stringSlice, k)
	}
	return stringSlice
}

// GenerateRandomDisplayColor will return a random _hue_ to be used when displaying a user.
// The UI should determine the right saturation and lightness in order to make it look right.
func GenerateRandomDisplayColor() int {
	rangeLower := 0
	rangeUpper := 360
	return rangeLower + rand.Intn(rangeUpper-rangeLower+1) //nolint
}

// GetHostnameFromURL will return the hostname component from a URL string.
func GetHostnameFromURL(u url.URL) string {
	return u.Host
}

// GetHostnameFromURLString will return the hostname component from a URL object.
func GetHostnameFromURLString(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}

	return u.Host
}

// GetHashtagsFromText returns all the #Hashtags from a string.
func GetHashtagsFromText(text string) []string {
	re := regexp.MustCompile(`#[a-zA-Z0-9_]+`)
	return re.FindAllString(text, -1)
}

// ShuffleStringSlice will shuffle a slice of strings.
func ShuffleStringSlice(s []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}
