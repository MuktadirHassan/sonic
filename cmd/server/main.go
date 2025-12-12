package main

import (
	"log"
	"net/http"
	"speedtest-server/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	// API endpoints
	mux.HandleFunc("/api/ping", corsMiddleware(handlers.PingHandler))
	mux.HandleFunc("/api/download", corsMiddleware(handlers.DownloadHandler))
	mux.HandleFunc("/api/upload", corsMiddleware(handlers.UploadHandler))

	log.Println("Speedtest server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// corsMiddleware adds CORS headers to the response
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
