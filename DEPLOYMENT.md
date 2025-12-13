# Deployment Guide

## Release Process

This project uses **release-please** for automated versioning and changelog generation based on conventional commits.

### How it Works

1. **Commit with Conventional Commits**:
   ```bash
   git commit -m "feat: add new speedtest feature"
   git commit -m "fix: resolve upload timeout issue"
   git commit -m "docs: update deployment guide"
   ```

2. **Push to main**:
   ```bash
   git push origin main
   ```

3. **Release Please creates a PR**:
   - Automatically bumps version based on commits
   - Generates/updates CHANGELOG.md
   - Creates a release PR

4. **Merge the Release PR**:
   - Release-please creates a GitHub release
   - Tags the commit with version (e.g., `v1.2.3`)

5. **Docker Action Triggers**:
   - Builds Docker image
   - Pushes to `ghcr.io/muktadirhassan/sonic`
   - Tags: `latest`, `v1.2.3`, `v1.2`, `v1`

6. **Binary Release**:
   - GoReleaser builds binaries for all platforms
   - Attaches to GitHub release as downloadable assets

## Conventional Commit Types

- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `feat!:` or `BREAKING CHANGE:` - Breaking change (major version bump)
- `docs:` - Documentation only
- `chore:` - Maintenance tasks
- `test:` - Test changes
- `perf:` - Performance improvements

## Docker Deployment

### Pull and Run

```bash
# Latest version
docker pull ghcr.io/muktadirhassan/sonic:latest
docker run -p 8080:8080 ghcr.io/muktadirhassan/sonic:latest

# Specific version
docker pull ghcr.io/muktadirhassan/sonic:v1.0.0
docker run -p 8080:8080 ghcr.io/muktadirhassan/sonic:v1.0.0
```

### Docker Compose

```yaml
version: '3.8'

services:
  sonic:
    image: ghcr.io/muktadirhassan/sonic:latest
    ports:
      - "8080:8080"
    restart: unless-stopped
    environment:
      - PORT=8080
```

## GitHub Actions Setup

### Required Secrets

No secrets are required! The workflows use `GITHUB_TOKEN` which is automatically provided.

### Workflow Files

- **`.github/workflows/test.yml`**: Runs tests on every PR and push
- **`.github/workflows/release-please.yml`**: Manages releases
- **`.github/workflows/docker.yml`**: Builds and pushes Docker images

### Manual Docker Build

If you want to trigger a Docker build manually:

1. Go to Actions tab in GitHub
2. Select "Docker Build and Push" workflow
3. Click "Run workflow"
4. Enter a tag or leave blank for latest

## Rollback

To rollback to a previous version:

```bash
# Pull specific version
docker pull ghcr.io/muktadirhassan/sonic:v1.0.0
```

