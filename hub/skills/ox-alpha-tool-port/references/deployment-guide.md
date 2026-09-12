# Frictionless Agent Distribution & Deployment Guide

Practical blueprint for packaging and deploying an autonomous AI operating layer (Go engine + Web dashboard monorepo) based on the proven Hermes Agent distribution pattern.

## 1. Multi-Stage Docker Packaging

Build a minimal production container (~30MB) combining a Node.js web dashboard bundle and a static Go binary:

```dockerfile
# Stage 1: Build Web Frontend
FROM node:22-alpine AS web-builder
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY app/web/package.json ./app/web/
COPY app/i18n ./app/i18n
RUN pnpm install --frozen-lockfile --filter @ox/web...
COPY app/web ./app/web
COPY app/assets ./app/assets
RUN pnpm --filter @ox/web build

# Stage 2: Build Static Go Engine
FROM golang:1.24-alpine AS engine-builder
WORKDIR /app/engine
RUN apk add --no-cache git ca-certificates
COPY engine/go.mod engine/go.sum ./
RUN go mod download
COPY engine/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/ox-engine ./cmd/server

# Stage 3: Minimal Production Alpine Runner
FROM alpine:3.21 AS runner
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S oxgroup && adduser -S oxuser -G oxgroup
COPY --from=engine-builder /app/ox-engine /app/ox-engine
COPY --from=web-builder /app/app/web/dist /app/app/web/dist
COPY ox.yml /app/ox.yml.example
RUN mkdir -p /app/data /app/.ox && chown -R oxuser:oxgroup /app
USER oxuser
EXPOSE 1060
HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:1060/health || exit 1
CMD ["/app/ox-engine"]
```

## 2. Docker Compose with Persistence & AI Forwarding

Ensure data persistence and host-proxy access (e.g. 9Router or Ollama running on host):

```yaml
services:
  ox:
    build:
      context: .
      dockerfile: Dockerfile
    image: theowlops/ox:latest
    container_name: ox-os
    restart: unless-stopped
    ports:
      - "1060:1060"
    volumes:
      - ./data:/app/data
      - ./.ox:/app/.ox
      - ./ox.yml:/app/ox.yml:ro
    environment:
      - OX_ENV=production
      - SERVER_PORT=1060
      - SERVER_HOST=0.0.0.0
      - AI_DEFAULT_PROVIDER=${AI_DEFAULT_PROVIDER:-custom}
      - AI_CUSTOM_BASE_URL=${AI_CUSTOM_BASE_URL:-http://host.docker.internal:20128/v1}
    extra_hosts:
      - "host.docker.internal:host-gateway"
```

## 3. Automated SSL Production Setup (Caddy Reverse Proxy)

Avoid certbot and manual nginx renewal configurations. Use Caddy to handle automatic Let's Encrypt certificates:

```yaml
# docker-compose.prod.yml
services:
  ox:
    image: theowlops/ox:latest
    restart: always
    volumes:
      - ./data:/app/data
      - ./ox.yml:/app/ox.yml:ro

  caddy:
    image: caddy:2-alpine
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - ./caddy_data:/data
      - ./caddy_config:/config
    depends_on:
      - ox
```

`Caddyfile`:
```caddy
yourdomain.com {
    reverse_proxy ox:1060
}
```

## 4. One-Liner Shell Installers

### Linux / macOS (`install.sh`)
```bash
# curl -fsSL https://raw.githubusercontent.com/TheOwlOps/OX/main/install.sh | bash
INSTALL_DIR="${OX_HOME:-$HOME/.ox-alpha}"
git clone --depth 1 https://github.com/TheOwlOps/OX.git "$INSTALL_DIR"
cd "$INSTALL_DIR"
if command -v docker &> /dev/null && docker compose version &> /dev/null; then
    docker compose up -d --build
    echo "OX running at http://localhost:1060"
    exit 0
fi
# local fallback
corepack enable && pnpm install && pnpm --filter @ox/web build
```

### Windows PowerShell (`install.ps1`)
```powershell
# irm https://raw.githubusercontent.com/TheOwlOps/OX/main/install.ps1 | iex
$InstallDir = if ($env:OX_HOME) { $env:OX_HOME } else { "$env:USERPROFILE\.ox-alpha" }
git clone --depth 1 https://github.com/TheOwlOps/OX.git $InstallDir
Set-Location $InstallDir
if ((Get-Command docker -ErrorAction SilentlyContinue) -and (docker compose version 2>$null)) {
    docker compose up -d --build
    Write-Host "[SUCCESS] Running at http://localhost:1060"
}
```

## 5. Automated SemVer Bumping & Binary Release CI

Do not rely on manual tagging or release drafting. Use GitHub Actions on `push: branches: [main]` to automate Semantic Versioning:
1. Fetch previous tags: `LATEST_TAG=$(git tag -l "v*" --sort=-v:refname | head -n 1)`.
2. Inspect commit history since latest tag:
   - Contains `breaking change|major:` → Bump Major (`v1.0.0`).
   - Contains `feat:|feature:` → Bump Minor (`v0.2.0`).
   - Default / `fix:|chore:|refactor:` → Bump Patch (`v0.1.1`).
3. Push new git tag: `git tag "$NEW_TAG" && git push origin "$NEW_TAG"`.
4. Compile standalone cross-platform binaries:
   - Windows desktop Wails `.exe`.
   - Windows Go engine `.exe`.
   - Linux Go engine static binary (`.tar.gz`).
5. Publish via `softprops/action-gh-release@v2` with `generate_release_notes: true` attaching all compiled binaries so users have instant one-click downloads.

## 6. Managing Repository About & Topics via Git Credentials

When `gh` CLI is absent on host, extract the user's authenticated GitHub token via git credential helper:
```bash
TOKEN=$(printf "protocol=https\nhost=github.com\n" | git credential fill | grep "password=" | cut -d= -f2)
```
Update repository metadata via GitHub REST API:
- `PATCH https://api.github.com/repos/{owner}/{repo}` with `{"description": "...", "homepage": "..."}` (ensure UTF-8 encoding for unicode emoji).
- `PUT https://api.github.com/repos/{owner}/{repo}/topics` with `{"names": ["ai-agents", "golang", "wails", ...]}`.

## 7. Decoupling Deprecated Node.js Monorepo Packages

When transitioning backend execution to a Go Engine (`engine/`):
1. Delete legacy TypeScript packages (`packages/ai`, `core`, `db`, `protocol`) to eliminate security alerts and dead dependencies.
2. Decouple CLI (`app/cli`) from `@ox/*` workspace packages by directing chat queries to Go engine's HTTP SSE endpoint (`POST http://localhost:1060/api/chat/stream`) and replacing custom YAML serializers with standard `yaml` npm package.
3. Update `pnpm-workspace.yaml` to only include active packages (`app/*`), ensuring `pnpm build` executes cleanly with zero orphaned workspace resolutions.
