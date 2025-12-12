package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PingHandler)

	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check headers
	if rr.Header().Get("Cache-Control") == "" {
		t.Errorf("handler returned missing Cache-Control header")
	}
}

func TestDownloadHandler(t *testing.T) {
	// Test default size (should be 10MB, but let's test a smaller custom size for speed)
	req, err := http.NewRequest("GET", "/api/download?size=1024", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DownloadHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read body to verify size
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	if len(body) != 1024 {
		t.Errorf("handler returned wrong body size: got %v want %v",
			len(body), 1024)
	}

	// content-type check
	if rr.Header().Get("Content-Type") != "application/octet-stream" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			rr.Header().Get("Content-Type"), "application/octet-stream")
	}
}

func TestUploadHandler(t *testing.T) {
	// Create a dummy payload
	payload := []byte("test payload for upload")
	req, err := http.NewRequest("POST", "/api/upload", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UploadHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
