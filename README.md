# Sonic Speedtest

[![Tests](https://github.com/MuktadirHassan/sonic/actions/workflows/test.yml/badge.svg)](https://github.com/MuktadirHassan/sonic/actions/workflows/test.yml)
[![Docker](https://github.com/MuktadirHassan/sonic/actions/workflows/docker.yml/badge.svg)](https://github.com/MuktadirHassan/sonic/actions/workflows/docker.yml)

A high-performance, production-ready speedtest server built with Go and a modern web interface.

## Features

- 🚀 **Fast & Accurate**: Download, Upload, Ping, and Jitter measurements
- 🎨 **Beautiful UI**: Modern TailwindCSS interface with animated gauges
- 📊 **History Tracking**: Local storage of last 10 test results (IndexedDB)
- 💻 **CLI Client**: Command-line tool for headless environments
- 🐳 **Docker Ready**: Production-optimized container images
- 🔄 **Auto-Releases**: Automated versioning and changelogs with release-please

## Quick Start

### Server (Docker Recommended)

The easiest way to run the server is with Docker:

```bash
docker pull ghcr.io/muktadirhassan/sonic:latest
docker run -p 8080:8080 ghcr.io/muktadirhassan/sonic:latest
```

Visit http://localhost:8080 to run a test.

**Alternative: Run from binary**

If you prefer not to use Docker, download the server binary from [releases](https://github.com/MuktadirHassan/sonic/releases/latest):

```bash
# Linux/macOS
wget https://github.com/MuktadirHassan/sonic/releases/latest/download/sonic-server_<version>_linux_amd64.tar.gz
tar -xzf sonic-server_<version>_linux_amd64.tar.gz
./sonic-server

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/MuktadirHassan/sonic/releases/latest/download/sonic-server_<version>_windows_amd64.zip" -OutFile sonic-server.zip
Expand-Archive sonic-server.zip
.\sonic-server\sonic-server.exe
```

### CLI Client

#### Option 1: Go Install (Easiest)

```bash
go install github.com/muktadirhassan/sonic/cmd/cli@latest
sonic-cli --server http://localhost:8080
```

#### Option 2: Download Binary

Download from [releases page](https://github.com/MuktadirHassan/sonic/releases/latest):

```bash
# Linux/macOS
wget https://github.com/MuktadirHassan/sonic/releases/latest/download/sonic-cli_<version>_linux_amd64.tar.gz
tar -xzf sonic-cli_<version>_linux_amd64.tar.gz
./sonic-cli

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/MuktadirHassan/sonic/releases/latest/download/sonic-cli_<version>_windows_amd64.zip" -OutFile sonic-cli.zip
Expand-Archive sonic-cli.zip
.\sonic-cli\sonic-cli.exe
```

**Available platforms:** Linux, macOS (Intel/Apple Silicon), Windows (amd64, arm64)

### From Source

```bash
# Clone the repository
git clone https://github.com/MuktadirHassan/sonic.git
cd sonic

# Run the server
go run cmd/server/main.go

# Or build and run
go build -o sonic-server ./cmd/server
./sonic-server
```

### CLI Client

```bash
# Build the CLI
go build -o sonic-cli ./cmd/cli

# Run a test
./sonic-cli

# Custom server
./sonic-cli --server http://speedtest.example.com:8080

# JSON output
./sonic-cli --json
```

## Development

### Prerequisites

- Go 1.24 or higher
- Docker (optional)

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Server
go build -o sonic-server ./cmd/server

# CLI
go build -o sonic-cli ./cmd/cli

# Docker
docker build -t sonic .
```

## API Endpoints

- `GET /` - Web interface
- `GET /api/ping` - Latency test endpoint
- `GET /api/download?size=<bytes>` - Download speed test
- `POST /api/upload` - Upload speed test

## Deployment

### Docker Compose

```yaml
version: '3.8'
services:
  sonic:
    image: ghcr.io/muktadirhassan/sonic:latest
    ports:
      - "8080:8080"
    restart: unless-stopped
```

## Contributing

Contributions are welcome! Please follow conventional commits for your PR titles:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `perf:` - Performance improvements
- `test:` - Test additions or changes

## License

MIT License - see LICENSE file for details