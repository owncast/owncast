package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getTempPipePath() string {
	return filepath.Join(os.TempDir(), "streampipe.flv")
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getRelativePathFromAbsolutePath(path string) string {
	pathComponents := strings.Split(path, "/")
	variant := pathComponents[len(pathComponents)-2]
	file := pathComponents[len(pathComponents)-1]
	return filepath.Join(variant, file)
}

func verifyError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func copy(src, dst string) {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println("Error creating", dst)
		fmt.Println(err)
		return
	}
}

func resetDirectories(configuration Config) {
	// Wipe the public, web-accessible hls data directory
	os.RemoveAll(configuration.PublicHLSPath)
	os.RemoveAll(configuration.PrivateHLSPath)
	os.MkdirAll(configuration.PublicHLSPath, 0777)
	os.MkdirAll(configuration.PrivateHLSPath, 0777)

	// Create private hls data dirs
	if !configuration.VideoSettings.EnablePassthrough || len(configuration.VideoSettings.StreamQualities) == 0 {
		for index := range configuration.VideoSettings.StreamQualities {
			os.MkdirAll(path.Join(configuration.PrivateHLSPath, strconv.Itoa(index)), 0777)
			os.MkdirAll(path.Join(configuration.PublicHLSPath, strconv.Itoa(index)), 0777)
		}
	} else {
		os.MkdirAll(path.Join(configuration.PrivateHLSPath, strconv.Itoa(0)), 0777)
		os.MkdirAll(path.Join(configuration.PublicHLSPath, strconv.Itoa(0)), 0777)
	}

	// Provide default files
	showStreamOfflineState(configuration)
	if !fileExists("webroot/thumbnail.png") {
		copy("static/logo.png", "webroot/thumbnail.png")
	}
}

func getClientIDFromRequest(req *http.Request) string {
	var clientID string
	xForwardedFor := req.Header.Get("X-FORWARDED-FOR")
	if xForwardedFor != "" {
		clientID = xForwardedFor
	} else {
		ipAddressString := req.RemoteAddr
		ipAddressComponents := strings.Split(ipAddressString, ":")
		ipAddressComponents[len(ipAddressComponents)-1] = ""
		clientID = strings.Join(ipAddressComponents, ":")
	}

	// fmt.Println("IP address determined to be", ipAddress)

	return clientID + req.UserAgent()
}
