package handlers

import (
	"crypto/rand"
	"io"
	"net/http"
	"strconv"
)

// Pre-allocate a buffer for download tests to avoid CPU overhead during serving.
// 1MB buffer of random data.
var randomBuffer []byte

func init() {
	randomBuffer = make([]byte, 1024*1024)
	rand.Read(randomBuffer)
}

// PingHandler handles latency checks.
// It returns a minimal response to measure RTT.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers to prevent caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(http.StatusOK)
}

// DownloadHandler serves random data for download speed measurement.
// It accepts a 'size' query parameter (in bytes). Default is 10MB.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for binary stream and no caching
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")

	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size <= 0 {
		size = 10 * 1024 * 1024 // Default 10MB
	}

	// Write random data in chunks
	written := int64(0)
	bufLen := int64(len(randomBuffer))

	for written < size {
		chunk := size - written
		if chunk > bufLen {
			chunk = bufLen
		}
		n, err := w.Write(randomBuffer[:chunk])
		if err != nil {
			return // Client disconnected or network error
		}
		written += int64(n)
	}
}

// UploadHandler consumes data for upload speed measurement.
// It reads the request body and discards it.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Disable caching
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")

	// Consume the body
	io.Copy(io.Discard, r.Body)
	w.WriteHeader(http.StatusOK)
}
