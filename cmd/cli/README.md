# Sonic Speedtest CLI

A command-line speedtest client for testing network performance against the Sonic Speedtest server.

## Features

- **Ping Test**: Measures latency with 10 requests
- **Jitter Test**: Calculates variance in latency
- **Download Test**: Multi-threaded download speed measurement
- **Upload Test**: Multi-threaded upload speed measurement
- **Color Output**: Beautiful terminal display with ANSI colors
- **JSON Output**: Machine-readable output for automation

## Installation

Build the CLI client:

```bash
go build -o sonic-cli.exe ./cmd/client
```

## Usage

### Basic Usage

Run a speedtest against the default server (localhost:8080):

```bash
./sonic-cli.exe
```

### Custom Server

Test against a custom server:

```bash
./sonic-cli.exe --server http://speedtest.example.com:8080
```

### JSON Output

Get results in JSON format (useful for scripts/automation):

```bash
./sonic-cli.exe --json
```

Example JSON output:
```json
{
  "ping": 12,
  "jitter": 2.3,
  "download": 145.67,
  "upload": 89.34
}
```

## Flags

- `--server`: Server URL (default: `http://localhost:8080`)
- `--json`: Output results in JSON format (default: false)

## Tests

Run the test suite:

```bash
go test ./cmd/client/...
```

## Example Output

```
╔═══════════════════════════════════════╗
║     SONIC SPEEDTEST CLI v1.0          ║
╚═══════════════════════════════════════╝

▶ Running Ping Test...
  Ping: 12 ms
  Jitter: 2.3 ms

▶ Running Download Test...
  Download: 145.67 Mbps

▶ Running Upload Test...
  Upload: 89.34 Mbps


╔═══════════════════════════════════════╗
║           TEST RESULTS                ║
╠═══════════════════════════════════════╣
║  Ping:                       12 ms ║
║  Jitter:                    2.3 ms ║
║  Download:               145.67 Mbps ║
║  Upload:                  89.34 Mbps ║
╚═══════════════════════════════════════╝
```
