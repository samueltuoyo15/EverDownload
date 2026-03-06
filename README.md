# EverDownload: Universal Video Downloader

EverDownload is a Go web app that downloads videos from 1000+ platforms using yt-dlp and ffmpeg. Metadata is Redis-cached for fast repeat lookups. No external APIs.

<img src="/public/demo.png" alt="EverDownload Demo">

## Stack

- **Go**: concurrent backend
- **yt-dlp + ffmpeg**: video extraction and format merging
- **Redis**: metadata caching (10 min TTL)
- **Docker + Docker Compose**: containerised deployment
- **Vanilla JS**: frontend (no framework)

## Supported Platforms

YouTube, TikTok, Instagram, Twitter/X, Reddit, Twitch, Vimeo, Facebook, Bilibili, Rumble, Kick, LinkedIn, Snapchat, Pinterest, Threads, Discord, and 1000+ more via yt-dlp.

## Getting Started

```bash
git clone https://github.com/samueltuoyo15/EverDownload.git
cd EverDownload
docker compose up --build -d
```

Open `http://localhost:8080`.

## Running Without Docker

Requires Go 1.24+, yt-dlp, ffmpeg, and a running Redis instance.

```bash
go run ./cmd/server/
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `REDIS_URL` | `redis://localhost:6379` | Redis connection URL |
| `PORT` | `5000` | Server port |

## Project Structure

```
cmd/server/main.go          entry point
internal/
  formats/formats.go        format parsing and labeling
  cache/cache.go            Redis wrapper
  ytdlp/ytdlp.go            yt-dlp process runner
  handlers/
    info.go                 POST /api/info
    download.go             GET /download
utils/allowed.hosts.go      URL allowlist
templates/index.html        frontend
static/style.css            styles
```

## Contributing

1. Fork the repo
2. Create a branch: `git checkout -b feature/your-feature`
3. Commit: `git commit -m "feat: your change"`
4. Push and open a Pull Request

## License

MIT: see [LICENSE](LICENSE).

## Author

**Samuel Tuoyo**
- LinkedIn: [samuel-tuoyo](https://www.linkedin.com/in/samuel-tuoyo-8568b62b6/)
- X: [@TuoyoS26091](https://x.com/TuoyoS26091)

---
[![Built with Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Powered by Redis](https://img.shields.io/badge/Redis-FF4438?style=flat-square&logo=redis&logoColor=white)](https://redis.io/)
[![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
