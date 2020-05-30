package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var directory = "hls"
var directoryHash string

func createIPFSDirectory() {
	sh := shell.NewShell("localhost:5001")
	newlyCreatedDirectoryHash, error := sh.AddDir(directory)
	verifyError((error))
	directoryHash = newlyCreatedDirectoryHash
}

func save(filePath string) string {
	file, err := os.Open(filePath)
	payload := bufio.NewReader(file)
	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(payload, shell.Pin(true))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	filename := filepath.Base(filePath)
	newDirectoryHash := addFileToDirectory(directoryHash, cid, filename)
	newFilePath := fmt.Sprintf("/ipfs/%s/%s", newDirectoryHash, filename)
	// fmt.Printf("added %s -> %s\n", filePath, newFilePath)

	return newFilePath
}

func saveData(stringData string, filename string) string {
	payload := strings.NewReader(stringData)

	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(payload)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	newDirectoryHash := addFileToDirectory(directoryHash, cid, filename)
	newFilePath := fmt.Sprintf("/ipfs/%s/%s", newDirectoryHash, filename)
	// fmt.Printf("added %s -> %s\n", filename, newFilePath)

	return newFilePath
}

func addFileToDirectory(directoryHash string, fileHash string, filename string) string {
	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")

	newDirectoryHash, err := sh.Patch("/ipfs/"+directoryHash, "add-link", filename, fileHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	return newDirectoryHash
}
