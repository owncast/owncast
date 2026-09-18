package storage

import (
	"crypto/md5" //nolint:gosec // S3 uses Content-MD5 for request integrity.
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestDeleteObjectsUsesContentMD5(t *testing.T) {
	type observedRequest struct {
		body   []byte
		header http.Header
		err    error
	}

	requests := make(chan observedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		requests <- observedRequest{body: body, header: r.Header.Clone(), err: err}
		w.Header().Set("Content-Type", "application/xml")
		_, _ = io.WriteString(w, `<DeleteResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"/>`)
	}))
	defer server.Close()

	client := s3.New(s3.Options{
		Region:       "us-east-1",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("access-key", "secret", "")),
		BaseEndpoint: aws.String(server.URL),
		UsePathStyle: true,
	})
	storage := &S3Storage{s3Client: client, s3Bucket: "bucket"}
	storage.deleteObjects([]s3object{{key: "hls/0/segment.ts"}})

	select {
	case request := <-requests:
		if request.err != nil {
			t.Fatalf("read request body: %v", request.err)
		}
		digest := md5.Sum(request.body) //nolint:gosec // S3 uses Content-MD5 for request integrity.
		want := base64.StdEncoding.EncodeToString(digest[:])
		if got := request.header.Get("Content-MD5"); got != want {
			t.Errorf("Content-MD5 = %q, want %q", got, want)
		}
		for name := range request.header {
			if strings.HasPrefix(strings.ToLower(name), "x-amz-checksum-") {
				t.Errorf("unexpected flexible checksum header %q", name)
			}
		}
		if got := request.header.Get("X-Amz-Sdk-Checksum-Algorithm"); got != "" {
			t.Errorf("X-Amz-Sdk-Checksum-Algorithm = %q, want empty", got)
		}
	default:
		t.Fatal("DeleteObjects request was not received")
	}
}
