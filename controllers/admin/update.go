package admin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/controllers"
	log "github.com/sirupsen/logrus"
)

/*
The auto-update relies on some guesses and hacks to determine how the binary
is being run.

The name of the service run under systemd must be "owncast" for
the guess to be successful. No huge loss if it's not successful, it'll just
not support the auto-update features.

It also determines if the binary is running under a container by figuring out
the container ID as a fallback to refuse an in-place update within a container.

In general it's disabled for everyone and the features are enabled only if
specific conditions are met.

1. Cannot be run inside a container.
2. Cannot be run from source (aka platform is "dev")
3. Must be run under systemd to support auto-restart.
4. Must be named "owncast" under systemd.
*/

// AutoUpdateOptions will return what auto update options are available.
func AutoUpdateOptions(w http.ResponseWriter, r *http.Request) {
	type autoUpdateOptionsResponse struct {
		SupportsUpdate bool `json:"supportsUpdate"`
		CanRestart     bool `json:"canRestart"`
	}

	updateOptions := autoUpdateOptionsResponse{
		SupportsUpdate: false,
		CanRestart:     false,
	}

	// Nothing is supported when running under "dev" or the feature is
	// explicitly disabled.
	if config.BuildPlatform == "dev" || !config.EnableAutoUpdate {
		controllers.WriteResponse(w, updateOptions)
		return
	}

	// If we are not in a container then we can update in place.
	if getContainerID() == "" {
		updateOptions.SupportsUpdate = true
	}

	updateOptions.CanRestart = isRunningUnderSystemD()

	controllers.WriteResponse(w, updateOptions)
}

// AutoUpdateStart will begin the auto update process.
func AutoUpdateStart(w http.ResponseWriter, r *http.Request) {
	// We return the console output directly to the client.
	w.Header().Set("Content-Type", "text/plain")

	// Download the installer and save it to a temp file.
	updater, err := downloadInstaller()
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "failed to download and run installer")
		return
	}

	// Run the installer.
	cmd := exec.Command("bash", updater)
	cmd.Env = append(os.Environ(), "NO_COLOR=true")

	pipeReader, pipeWriter := io.Pipe()
	defer pipeWriter.Close() // nolint: errcheck
	cmd.Stdout = pipeWriter
	cmd.Stderr = pipeWriter
	go func() {
		_, err := io.Copy(w, pipeReader)
		if err != nil {
			log.Debugln(err)
		}
	}()

	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "unable to update: "+err.Error())
		return
	}

	if err := cmd.Run(); err != nil {
		log.Debugln(err)
		controllers.WriteSimpleResponse(w, false, "unable to update: "+err.Error())
		return
	}
}

// AutoUpdateForceQuit will force quit the service.
func AutoUpdateForceQuit(w http.ResponseWriter, r *http.Request) {
	log.Warnln("Owncast is exiting due to request.")
	go func() {
		os.Exit(0)
	}()
	controllers.WriteSimpleResponse(w, true, "forcing quit")
}

func downloadInstaller() (string, error) {
	installer := "https://owncast.online/install.sh"
	out, err := os.CreateTemp(os.TempDir(), "updater.sh")
	if err != nil {
		log.Errorln(err)
		return "", err
	}
	defer out.Close()

	// Get the installer script
	resp, err := http.Get(installer)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the installer to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return out.Name(), nil
}

// Check to see if owncast is listed as a running service under systemd.
// This is a guess and nothing more.
func isRunningUnderSystemD() bool {
	cmd := exec.Command("/bin/sh",
		"-c",
		"systemctl list-units --type=service --state=running | grep owncast",
	)

	out, err := cmd.Output()
	if err != nil {
		log.Errorln(err)
		return false
	}

	if out != nil {
		return true
	}

	return false
}

// Taken from https://stackoverflow.com/questions/23513045/how-to-check-if-a-process-is-running-inside-docker-container
func getContainerID() string {
	pid := os.Getppid()
	cgroupPath := fmt.Sprintf("/proc/%s/cgroup", strconv.Itoa(pid))
	containerID := ""
	content, err := ioutil.ReadFile(cgroupPath) //nolint:gosec
	if err != nil {
		return containerID
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		field := strings.Split(line, ":")
		if len(field) < 3 {
			continue
		}
		cgroupPath := field[2]
		if len(cgroupPath) < 64 {
			continue
		}
		// Non-systemd Docker
		// 5:net_prio,net_cls:/docker/de630f22746b9c06c412858f26ca286c6cdfed086d3b302998aa403d9dcedc42
		// 3:net_cls:/kubepods/burstable/pod5f399c1a-f9fc-11e8-bf65-246e9659ebfc/9170559b8aadd07d99978d9460cf8d1c71552f3c64fefc7e9906ab3fb7e18f69
		pos := strings.LastIndex(cgroupPath, "/")
		if pos > 0 {
			idLen := len(cgroupPath) - pos - 1
			if idLen == 64 {
				// p.InDocker = true
				// docker id
				containerID = cgroupPath[pos+1 : pos+1+64]
				// logs.Debug("pid:%v in docker id:%v", pid, id)
				return containerID
			}
		}
		// systemd Docker
		// 5:net_cls:/system.slice/docker-afd862d2ed48ef5dc0ce8f1863e4475894e331098c9a512789233ca9ca06fc62.scope
		dockerStr := "docker-"
		pos = strings.Index(cgroupPath, dockerStr)
		if pos > 0 {
			posScope := strings.Index(cgroupPath, ".scope")
			idLen := posScope - pos - len(dockerStr)
			if posScope > 0 && idLen == 64 {
				containerID = cgroupPath[pos+len(dockerStr) : pos+len(dockerStr)+64]
				return containerID
			}
		}
	}
	return containerID
}
