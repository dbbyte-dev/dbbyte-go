package dbbyte

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient points a DBByte client at an httptest server running h.
func newTestClient(t *testing.T, h http.HandlerFunc) *DBByte {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)

	dbb, err := New(ClientOptions{APIKey: "dbk_test", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return dbb
}

func decodeBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var got map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	return got
}

func TestUploadURL(t *testing.T) {
	dbb := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/storage/my-bucket/upload-url" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer dbk_test" {
			t.Errorf("Authorization = %q", got)
		}
		body := decodeBody(t, r)
		if body["name"] != "song.mp3" || body["contentType"] != "audio/mpeg" {
			t.Errorf("unexpected body: %v", body)
		}
		if body["size"] != float64(5242880) || body["expiresIn"] != float64(900) {
			t.Errorf("unexpected size/expiresIn: %v", body)
		}
		if _, ok := body["overwrite"]; ok {
			t.Errorf("overwrite should be omitted when false: %v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"url":"https://example.test/upload/tok","method":"PUT","headers":{"Content-Type":"audio/mpeg"},"expiresIn":900}`))
	})

	got, err := dbb.Objects.UploadURL(context.Background(), "my-bucket", UploadURLOptions{
		Name:        "song.mp3",
		ContentType: "audio/mpeg",
		Size:        5242880,
		ExpiresIn:   900,
	})
	if err != nil {
		t.Fatalf("UploadURL: %v", err)
	}
	if got.URL != "https://example.test/upload/tok" || got.Method != "PUT" || got.ExpiresIn != 900 {
		t.Errorf("unexpected response: %+v", got)
	}
	if got.Headers["Content-Type"] != "audio/mpeg" {
		t.Errorf("unexpected headers: %v", got.Headers)
	}
}

func TestUploadURLRequiresName(t *testing.T) {
	dbb := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("no request should be sent when Name is empty")
	})
	if _, err := dbb.Objects.UploadURL(context.Background(), "my-bucket", UploadURLOptions{}); err == nil {
		t.Fatal("expected an error for an empty Name")
	}
}

func TestUploadURLAPIError(t *testing.T) {
	dbb := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		_, _ = w.Write([]byte(`{"error":"File too large. Max allowed size is 50MB."}`))
	})
	_, err := dbb.Objects.UploadURL(context.Background(), "my-bucket", UploadURLOptions{Name: "big.bin", Size: 60 << 20})
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected *Error with status 413, got %v", err)
	}
}

func TestPurge(t *testing.T) {
	dbb := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/purge" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body := decodeBody(t, r)
		if body["slug"] != "my-bucket" || body["scope"] != "file" || body["path"] != "logo.png" {
			t.Errorf("unexpected body: %v", body)
		}
		_, _ = w.Write([]byte(`{"success":true,"purged":4,"failed":[]}`))
	})

	got, err := dbb.CDN.Purge(context.Background(), PurgeOptions{Slug: "my-bucket", Scope: PurgeFile, Path: "logo.png"})
	if err != nil {
		t.Fatalf("Purge: %v", err)
	}
	if !got.Success || got.Purged != 4 || len(got.Failed) != 0 {
		t.Errorf("unexpected response: %+v", got)
	}
}

func TestBucketsCreate(t *testing.T) {
	dbb := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/buckets" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		body := decodeBody(t, r)
		if body["name"] != "Media" || body["region"] != "us-west-1" || body["cdn"] != true {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"success":true,"bucket":{"id":"b1","name":"Media","slug":"media","region":"us-west-1","cdn":true}}`))
	})

	got, err := dbb.Buckets.Create(context.Background(), CreateBucketOptions{Name: "Media", Region: "us-west-1", CDN: true})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !got.Success || got.Bucket.Slug != "media" || !got.Bucket.CDN {
		t.Errorf("unexpected response: %+v", got)
	}
}
