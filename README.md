# RaProxy-Streaming

A high-performance Go proxy server for streaming video content, built with Fiber framework. This server bypasses CORS restrictions for video streaming by proxying M3U8 playlists and video segments.

## Features

- 🚀 High-performance proxy using Fiber and fasthttp
- 📺 M3U8 playlist proxying with URL rewriting
- 🎬 Video segment streaming (TS, MP4 fragments)
- 🎥 Direct video file proxying
- 🔒 CORS-enabled endpoints
- 📦 Modular architecture with clean separation of concerns
- 🛡️ Error handling and logging
- ⚡ Fast and efficient streaming

## Project Structure

```
raproxy-streaming/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handler/
│   │   └── proxy.go             # HTTP handlers
│   ├── httpclient/
│   │   └── client.go            # HTTP client wrapper
│   ├── logger/
│   │   └── logger.go            # Logging utility
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   └── service/
│       └── proxy.go             # Business logic
├── go.mod                        # Go modules
└── README.md                     # Documentation
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/RanggaCasper/raproxy-streaming.git
cd raproxy-streaming
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o bin/server cmd/server/main.go
```

## Usage

### Running the Server

```bash
# Using go run
go run cmd/server/main.go

# Or using the built binary
./bin/server
```

The server will start on port 3000 by default. You can change the port using the `PORT` environment variable:

```bash
PORT=8080 #example
go run cmd/server/main.go
```

### API Endpoints

#### 1. Proxy M3U8 Playlist
```
GET /proxy/m3u8?url=<m3u8_url>&referer=<referer>
```

**Parameters:**
- `url` (required): The M3U8 playlist URL to proxy
- `referer` (optional): Referer header to send with the request

**Example:**
```bash
curl "http://localhost:3000/proxy/m3u8?url=https://example.com/playlist.m3u8&referer=https://example.com"
```

#### 2. Proxy Video Segment
```
GET /proxy/segment?url=<segment_url>&referer=<referer>
```

**Parameters:**
- `url` (required): The video segment URL to proxy
- `referer` (optional): Referer header to send with the request

**Example:**
```bash
curl "http://localhost:3000/proxy/segment?url=https://example.com/segment.ts&referer=https://example.com"
```

#### 3. Proxy Video File
```
GET /proxy/video?url=<video_url>&referer=<referer>
```

**Parameters:**
- `url` (required): The video file URL to proxy
- `referer` (optional): Referer header to send with the request

**Example:**
```bash
curl "http://localhost:3000/proxy/video?url=https://example.com/video.mp4&referer=https://example.com"
```

#### 4. Health Check
```
GET /health
```

Returns server health status.

## Configuration

Configuration can be modified in `internal/config/config.go`:

- `Timeout`: Overall request timeout (default: 60s)
- `ConnectTimeout`: Connection timeout (default: 10s)
- `MaxRedirects`: Maximum number of redirects to follow (default: 10)

## Architecture

### Layers

1. **Handler Layer** (`internal/handler`): Handles HTTP requests and responses
2. **Service Layer** (`internal/service`): Contains business logic
3. **HTTP Client** (`internal/httpclient`): Manages external HTTP requests
4. **Config** (`internal/config`): Application configuration
5. **Logger** (`internal/logger`): Logging functionality
6. **Routes** (`internal/routes`): Route registration

### Flow

```
HTTP Request → Routes → Handler → Service → HTTP Client → External API
                  ↓         ↓         ↓          ↓
               Fiber    Response  Business   fasthttp
                                  Logic
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Building for Production

```bash
# Build for current platform
go build -o bin/server cmd/server/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe cmd/server/main.go
```

## Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]
```

Build and run:
```bash
docker build -t raproxy-streaming .
docker run -p 3000:3000 raproxy-streaming
```

## Environment Variables

- `PORT`: Server port (default: 3000)

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
