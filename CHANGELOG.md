# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-14

### Added
- Initial release
- M3U8 playlist proxy with URL rewriting
- Video segment proxy (TS, MP4 fragments)
- Direct video file proxy with streaming
- Health check endpoint
- CORS support for all endpoints
- Makefile for build automation
- Live reload support with Air
- Modular architecture with clean separation of concerns
- Custom HTTP client with timeout configuration
- Structured logging
- Error handling with custom error types
- Graceful shutdown support

### Features
- **Handler Layer**: Clean HTTP request/response handling
- **Service Layer**: Business logic for proxy operations
- **HTTP Client**: Configurable fasthttp wrapper
- **Configuration**: Centralized config management
- **Logging**: Structured error and info logging
- **Routes**: RESTful endpoint organization

### Performance
- Connection pooling via fasthttp
- Streaming support for large files
- Efficient memory usage
- Low latency proxy operations

### Security
- CORS enabled (configurable)
- User-Agent spoofing for compatibility
- Referer header forwarding

## [Unreleased]

### Planned
- Redis caching layer
- Rate limiting middleware
- Authentication/Authorization