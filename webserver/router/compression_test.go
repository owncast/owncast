package router

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPCompressionAdapterStreamsSSEWithoutCompression(t *testing.T) {
	compress, err := newHTTPCompressionAdapter()
	if err != nil {
		t.Fatal(err)
	}

	flushed := make(chan struct{})
	release := make(chan struct{})

	handler := compress(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, ": connected\n\n"); err != nil {
			t.Error(err)
			return
		}
		w.(http.Flusher).Flush()
		close(flushed)
		<-release
	}))
	server := httptest.NewServer(handler)
	defer server.Close()
	defer close(release)

	client := &http.Client{Transport: &http.Transport{DisableCompression: true}}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	response := make(chan *http.Response, 1)
	go func() {
		resp, requestErr := client.Do(req)
		if requestErr != nil {
			t.Errorf("request: %v", requestErr)
			return
		}
		response <- resp
	}()

	select {
	case <-flushed:
	case <-time.After(time.Second):
		t.Fatal("SSE handler did not flush the initial comment")
	}

	var resp *http.Response
	select {
	case resp = <-response:
	case <-time.After(time.Second):
		t.Fatal("SSE response headers were not delivered")
	}
	defer resp.Body.Close()

	if got := resp.Header.Get("Content-Encoding"); got != "" {
		t.Fatalf("Content-Encoding = %q, want empty", got)
	}
	line, err := bufio.NewReader(resp.Body).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if line != ": connected\n" {
		t.Fatalf("first SSE line = %q", line)
	}
}

func TestHTTPCompressionAdapterStillCompressesJSON(t *testing.T) {
	compress, err := newHTTPCompressionAdapter()
	if err != nil {
		t.Fatal(err)
	}

	handler := compress(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, strings.Repeat("x", 1024))
	}))
	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{Transport: &http.Transport{DisableCompression: true}}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if got := resp.Header.Get("Content-Encoding"); got != "gzip" {
		t.Fatalf("Content-Encoding = %q, want gzip", got)
	}
}
