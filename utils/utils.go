package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

func GetIndexFromFilePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]

	return variant
}

// Copy copies the file to destination.
func Copy(source, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destination, input, 0600)
}

// Move moves the file to destination.
func Move(source, destination string) error {
	return os.Rename(source, destination)
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
	}

	for _, botString := range botStrings {
		if strings.Contains(strings.ToLower(userAgent), botString) {
			return true
		}
	}

	ua := user_agent.New(userAgent)
	return ua.Bot()
}

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
		panic(err)
	}

	return buf.String()
}

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
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// GetCacheDurationSecondsForPath will return the number of seconds to cache an item.
func GetCacheDurationSecondsForPath(filePath string) int {
	if path.Base(filePath) == "thumbnail.jpg" {
		// Thumbnails re-generate during live
		return 20
	} else if path.Ext(filePath) == ".js" || path.Ext(filePath) == ".css" {
		// Cache javascript & CSS
		return 60
	} else if path.Ext(filePath) == ".ts" {
		// Cache video segments as long as you want. They can't change.
		// This matters most for local hosting of segments for recordings
		// and not for live or 3rd party storage.
		return 31557600
	} else if path.Ext(filePath) == ".m3u8" {
		return 0
	}

	// Default cache length in seconds
	return 30
}

func IsValidUrl(urlToTest string) bool {
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
		} else {
			log.Warnln(ffmpegPath, "is an invalid path to ffmpeg will try to use a copy in your path, if possible")
		}
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
		log.Fatalln("Unable to determine path to ffmpeg. Please specify it in the config file.")
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
	//source: https://stackoverflow.com/a/60128480
	if mode&0111 == 0 {
		return errors.New("ffmpeg path is not executable")
	}

	return nil
}

// Removes the directory and makes it again. Throws fatal error on failure.
func CleanupDirectory(path string) {
	log.Traceln("Cleaning", path)
	if err := os.RemoveAll(path); err != nil {
		log.Fatalln("Unable to remove directory. Please check the ownership and permissions", err)
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalln("Unable to create directory. Please check the ownership and permissions", err)
	}
}
