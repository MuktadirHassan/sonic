package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunPingTest(t *testing.T) {
	// Create a test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Temporarily modify the ping iterations for testing
	originalPingIterations := 10
	// Run with fewer iterations for faster tests

	// Run ping test
	ping, jitter := runPingTest(server.URL)

	// Verify results - ping can be 0 for instant local responses
	if ping < 0 || ping > 1000 {
		t.Errorf("Expected reasonable ping value (0-1000ms), got %.2f", ping)
	}
	if jitter < 0 {
		t.Errorf("Expected jitter >= 0, got %.2f", jitter)
	}

	_ = originalPingIterations // avoid unused variable
}

func TestRunDownloadTest(t *testing.T) {
	// Create a test server that sends data
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make([]byte, 1024*1024) // 1MB
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(data)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Run download test
	speed := runDownloadTest(server.URL)

	// Verify results
	if speed <= 0 {
		t.Errorf("Expected download speed > 0, got %.2f", speed)
	}
}

func TestRunUploadTest(t *testing.T) {
	// Create a test server that accepts upload
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Run upload test
	speed := runUploadTest(server.URL)

	// Verify results
	if speed <= 0 {
		t.Errorf("Expected upload speed > 0, got %.2f", speed)
	}
}

func TestPingTestWithSlowServer(t *testing.T) {
	// Create a slow responding server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Run ping test
	ping, _ := runPingTest(server.URL)

	// Verify ping reflects the delay
	if ping < 40 { // Should be at least 40ms (with some margin)
		t.Errorf("Expected ping >= 40ms due to server delay, got %.2f", ping)
	}
}
