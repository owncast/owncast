package main

import (
	"os"
	"path/filepath"
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

func verifyError(e error) {
	if e != nil {
		panic(e)
	}
}
