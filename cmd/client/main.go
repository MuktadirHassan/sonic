package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

type SpeedTestResult struct {
	Ping     float64 `json:"ping"`
	Jitter   float64 `json:"jitter"`
	Download float64 `json:"download"`
	Upload   float64 `json:"upload"`
}

var (
	serverURL  string
	jsonOutput bool
)

func init() {
	flag.StringVar(&serverURL, "server", "http://localhost:8080", "Speedtest server URL")
	flag.BoolVar(&jsonOutput, "json", false, "Output results in JSON format")
}

func main() {
	flag.Parse()

	if !jsonOutput {
		printBanner()
	}

	result := SpeedTestResult{}

	// Run Ping Test
	if !jsonOutput {
		fmt.Printf("\n%s%s▶ Running Ping Test...%s\n", colorBold, colorCyan, colorReset)
	}
	ping, jitter := runPingTest(serverURL)
	result.Ping = ping
	result.Jitter = jitter
	if !jsonOutput {
		fmt.Printf("  %sPing:%s %.0f ms\n", colorGreen, colorReset, ping)
		fmt.Printf("  %sJitter:%s %.1f ms\n", colorYellow, colorReset, jitter)
	}

	// Run Download Test
	if !jsonOutput {
		fmt.Printf("\n%s%s▶ Running Download Test...%s\n", colorBold, colorCyan, colorReset)
	}
	download := runDownloadTest(serverURL)
	result.Download = download
	if !jsonOutput {
		fmt.Printf("  %sDownload:%s %.2f Mbps\n", colorGreen, colorReset, download)
	}

	// Run Upload Test
	if !jsonOutput {
		fmt.Printf("\n%s%s▶ Running Upload Test...%s\n", colorBold, colorPurple, colorReset)
	}
	upload := runUploadTest(serverURL)
	result.Upload = upload
	if !jsonOutput {
		fmt.Printf("  %sUpload:%s %.2f Mbps\n", colorGreen, colorReset, upload)
	}

	// Output results
	if jsonOutput {
		jsonData, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonData))
	} else {
		printSummary(result)
	}
}

func printBanner() {
	fmt.Printf("%s%s", colorCyan, colorBold)
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║     SONIC SPEEDTEST CLI v1.0          ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Printf("%s", colorReset)
}

func printSummary(result SpeedTestResult) {
	fmt.Printf("\n%s%s", colorBold, colorWhite)
	fmt.Println("\n╔═══════════════════════════════════════╗")
	fmt.Println("║           TEST RESULTS                ║")
	fmt.Println("╠═══════════════════════════════════════╣")
	fmt.Printf("║  %-15s %15.0f ms ║\n", "Ping:", result.Ping)
	fmt.Printf("║  %-15s %15.1f ms ║\n", "Jitter:", result.Jitter)
	fmt.Printf("║  %-15s %12.2f Mbps ║\n", "Download:", result.Download)
	fmt.Printf("║  %-15s %12.2f Mbps ║\n", "Upload:", result.Upload)
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Printf("%s\n", colorReset)
}

func runPingTest(baseURL string) (float64, float64) {
	client := &http.Client{Timeout: 5 * time.Second}
	pings := []float64{}

	for i := 0; i < 10; i++ {
		start := time.Now()
		resp, err := client.Get(fmt.Sprintf("%s/api/ping", baseURL))
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error during ping: %v\n", err)
			continue
		}
		resp.Body.Close()

		pings = append(pings, float64(elapsed))
		time.Sleep(50 * time.Millisecond)
	}

	if len(pings) == 0 {
		return 0, 0
	}

	// Calculate min ping
	minPing := pings[0]
	sum := 0.0
	for _, p := range pings {
		if p < minPing {
			minPing = p
		}
		sum += p
	}
	avg := sum / float64(len(pings))

	// Calculate jitter (average absolute deviation)
	jitterSum := 0.0
	for _, p := range pings {
		jitterSum += math.Abs(p - avg)
	}
	jitter := jitterSum / float64(len(pings))

	return minPing, jitter
}

func runDownloadTest(baseURL string) float64 {
	client := &http.Client{Timeout: 30 * time.Second}
	threads := 4
	targetDuration := 5 * time.Second
	bufferSize := 4 * 1024 * 1024 // 4MB

	var totalBytes int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	startTime := time.Now()
	running := true

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for running && time.Since(startTime) < targetDuration {
				resp, err := client.Get(fmt.Sprintf("%s/api/download?size=%d", baseURL, bufferSize))
				if err != nil {
					break
				}

				buf := make([]byte, 32*1024)
				for {
					n, err := resp.Body.Read(buf)
					if n > 0 {
						mu.Lock()
						totalBytes += int64(n)
						mu.Unlock()
					}
					if err != nil {
						break
					}
					if !running {
						break
					}
				}
				resp.Body.Close()

				if !running {
					break
				}
			}
		}()
	}

	wg.Wait()
	running = false

	duration := time.Since(startTime).Seconds()
	speedBps := float64(totalBytes*8) / duration
	speedMbps := speedBps / 1000000

	return speedMbps
}

func runUploadTest(baseURL string) float64 {
	client := &http.Client{Timeout: 30 * time.Second}
	threads := 3
	targetDuration := 5 * time.Second
	payloadSize := 10 * 1024 * 1024 // 10MB

	payload := make([]byte, payloadSize)

	var totalBytes int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	startTime := time.Now()
	running := true

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for running && time.Since(startTime) < targetDuration {
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/upload", baseURL), bytes.NewReader(payload))
				if err != nil {
					break
				}
				req.Header.Set("Content-Type", "application/octet-stream")

				resp, err := client.Do(req)
				if err != nil {
					break
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				if !running {
					break
				}

				mu.Lock()
				totalBytes += int64(payloadSize)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	running = false

	duration := time.Since(startTime).Seconds()
	speedBps := float64(totalBytes*8) / duration
	speedMbps := speedBps / 1000000

	return speedMbps
}
