package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetHashtagsFromText(t *testing.T) {
	text := `Some text with a #hashtag goes here.\n\n
	Another #secondhashtag, goes here.\n\n
	#thirdhashtag`

	hashtags := GetHashtagsFromText(text)

	if hashtags[0] != "#hashtag" || hashtags[1] != "#secondhashtag" || hashtags[2] != "#thirdhashtag" {
		t.Error("Incorrect hashtags fetched from text.")
	}
}

func TestPercentageUtilsTest(t *testing.T) {
	total := 42
	number := 18

	percent := IntPercentage(number, total)

	if percent != 42 {
		t.Error("Incorrect percentage calculation.")
	}
}

func TestValidatedFfmpegPathPrefersLocalBinary(t *testing.T) {
	// Skip on Windows as executable permissions work differently.
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows")
	}

	// Create a temporary directory to act as our working directory.
	tempDir, err := os.MkdirTemp("", "ffmpeg-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a fake ffmpeg executable in the temp directory.
	fakeFfmpegPath := filepath.Join(tempDir, "ffmpeg")
	if err := os.WriteFile(fakeFfmpegPath, []byte("#!/bin/sh\necho fake"), 0o755); err != nil {
		t.Fatalf("Failed to create fake ffmpeg: %v", err)
	}

	// Save the current working directory and change to the temp directory.
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Call ValidatedFfmpegPath with an empty string to trigger auto-detection.
	result := ValidatedFfmpegPath("")

	// The result should be the absolute path to our local fake ffmpeg,
	// not the system ffmpeg in /usr/bin or similar.
	expectedPath, err := filepath.Abs("./ffmpeg")
	if err != nil {
		t.Fatalf("Failed to get absolute path for fake ffmpeg: %v", err)
	}

	if result != expectedPath {
		t.Errorf("Expected local ffmpeg path %q, got %q", expectedPath, result)
	}
}

func TestDecodeBase64Image(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantExtension string
		wantErr       bool
	}{
		{
			name:          "valid PNG",
			input:         "data:image/png;base64,iVBORw0KGgo=",
			wantExtension: ".png",
		},
		{
			name:          "valid JPEG",
			input:         "data:image/jpeg;base64,/9j/4AAQ",
			wantExtension: ".jpeg",
		},
		{
			name:          "valid GIF",
			input:         "data:image/gif;base64,R0lGODlh",
			wantExtension: ".gif",
		},
		{
			name:          "valid SVG",
			input:         "data:image/svg+xml;base64,PHN2Zz4=",
			wantExtension: ".svg",
		},
		{
			name:          "valid ICO via image/x-icon",
			input:         "data:image/x-icon;base64,AAABAA==",
			wantExtension: ".ico",
		},
		{
			name:          "valid ICO via image/vnd.microsoft.icon",
			input:         "data:image/vnd.microsoft.icon;base64,AAABAA==",
			wantExtension: ".ico",
		},
		{
			name:    "unsupported content type",
			input:   "data:image/webp;base64,AAABAA==",
			wantErr: true,
		},
		{
			name:    "missing comma separator",
			input:   "data:image/png;base64AAAA",
			wantErr: true,
		},
		{
			name:    "missing header",
			input:   ",AAABAA==",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, err := DecodeBase64Image(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ext != tt.wantExtension {
				t.Errorf("expected extension %q, got %q", tt.wantExtension, ext)
			}
		})
	}
}

func TestVerifyFFMpegPath(t *testing.T) {
	// Test with non-existent path.
	err := VerifyFFMpegPath("/nonexistent/path/to/ffmpeg")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Test with a directory instead of a file.
	tempDir, err := os.MkdirTemp("", "ffmpeg-dir-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = VerifyFFMpegPath(tempDir)
	if err == nil {
		t.Error("Expected error when path is a directory")
	}

	// Test with a non-executable file (Unix only).
	if runtime.GOOS != "windows" {
		nonExecFile := filepath.Join(tempDir, "nonexec")
		if err := os.WriteFile(nonExecFile, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create non-executable file: %v", err)
		}

		err = VerifyFFMpegPath(nonExecFile)
		if err == nil {
			t.Error("Expected error for non-executable file")
		}
	}

	// Test with a valid executable file.
	if runtime.GOOS != "windows" {
		execFile := filepath.Join(tempDir, "executable")
		if err := os.WriteFile(execFile, []byte("#!/bin/sh\necho test"), 0o755); err != nil {
			t.Fatalf("Failed to create executable file: %v", err)
		}

		err = VerifyFFMpegPath(execFile)
		if err != nil {
			t.Errorf("Expected no error for valid executable, got: %v", err)
		}
	}
}
