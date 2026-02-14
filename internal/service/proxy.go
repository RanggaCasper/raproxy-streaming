package service

import (
	"fmt"
	"net/url"
	"strings"

	"raproxy-streaming/internal/httpclient"
	"raproxy-streaming/internal/logger"

	"github.com/valyala/fasthttp"
)

// ProxyService handles proxy operations
type ProxyService struct {
	client *httpclient.Client
	logger *logger.Logger
}

// NewProxyService creates a new proxy service
func NewProxyService(client *httpclient.Client, logger *logger.Logger) *ProxyService {
	return &ProxyService{
		client: client,
		logger: logger,
	}
}

// ProxyM3U8 fetches and rewrites M3U8 playlist
func (s *ProxyService) ProxyM3U8(targetURL, referer string) ([]byte, string, error) {
	resp, err := s.client.Get(targetURL, referer)
	if err != nil {
		if httpErr, ok := httpclient.IsHTTPError(err); ok {
			s.logger.Error("Proxy m3u8 HTTP error %d: %s", httpErr.StatusCode, targetURL)
			return nil, "", fmt.Errorf("HTTP %d: %s", httpErr.StatusCode, targetURL)
		}
		s.logger.Error("Proxy m3u8 failed for %s: %v", targetURL, err)
		return nil, "", fmt.Errorf("failed to fetch m3u8: %w", err)
	}
	defer fasthttp.ReleaseResponse(resp)

	content := string(resp.Body())
	contentType := string(resp.Header.Peek("Content-Type"))
	if contentType == "" {
		contentType = "application/vnd.apple.mpegurl"
	}

	// Rewrite URLs in M3U8
	if strings.Contains(content, "#EXTINF") || strings.Contains(content, "#EXT-X-STREAM-INF") {
		baseURL := getBaseURL(targetURL)
		content = s.rewriteM3U8Content(content, baseURL, referer)
	}

	return []byte(content), contentType, nil
}

// ProxySegment fetches video segment
func (s *ProxyService) ProxySegment(targetURL, referer string) ([]byte, string, error) {
	resp, err := s.client.Get(targetURL, referer)
	if err != nil {
		if httpErr, ok := httpclient.IsHTTPError(err); ok {
			s.logger.Error("Proxy segment HTTP error %d: %s", httpErr.StatusCode, targetURL)
			return nil, "", fmt.Errorf("HTTP %d", httpErr.StatusCode)
		}
		s.logger.Error("Proxy segment failed for %s: %v", targetURL, err)
		return nil, "", fmt.Errorf("failed to fetch segment: %w", err)
	}
	defer fasthttp.ReleaseResponse(resp)

	contentType := string(resp.Header.Peek("Content-Type"))
	if contentType == "" {
		contentType = "video/mp2t"
	}

	// Copy body to new slice
	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return body, contentType, nil
}

// ProxyVideo fetches video with streaming support
func (s *ProxyService) ProxyVideo(targetURL, referer string) (*fasthttp.Response, error) {
	resp, err := s.client.Get(targetURL, referer)
	if err != nil {
		if httpErr, ok := httpclient.IsHTTPError(err); ok {
			s.logger.Error("Proxy video HTTP error %d: %s", httpErr.StatusCode, targetURL)
			return nil, fmt.Errorf("HTTP %d", httpErr.StatusCode)
		}
		s.logger.Error("Proxy video failed for %s: %v", targetURL, err)
		return nil, fmt.Errorf("failed to fetch video: %w", err)
	}

	return resp, nil
}

// rewriteM3U8Content rewrites segment URLs in M3U8 content
func (s *ProxyService) rewriteM3U8Content(content, baseURL, referer string) string {
	lines := strings.Split(content, "\n")
	rewritten := make([]string, 0, len(lines))

	for _, line := range lines {
		stripped := strings.TrimSpace(line)

		// Check if this is a segment/playlist URL (not a comment)
		if stripped != "" && !strings.HasPrefix(stripped, "#") {
			var segmentURL string

			// Check if URL is absolute or relative
			if strings.HasPrefix(stripped, "http://") || strings.HasPrefix(stripped, "https://") {
				segmentURL = stripped
			} else {
				segmentURL = baseURL + stripped
			}

			// Determine if it's another m3u8 or a segment
			if strings.HasSuffix(stripped, ".m3u8") || strings.Contains(stripped, "m3u8") {
				rewritten = append(rewritten, fmt.Sprintf(
					"/proxy/m3u8?url=%s&referer=%s",
					url.QueryEscape(segmentURL),
					url.QueryEscape(referer),
				))
			} else {
				rewritten = append(rewritten, fmt.Sprintf(
					"/proxy/segment?url=%s&referer=%s",
					url.QueryEscape(segmentURL),
					url.QueryEscape(referer),
				))
			}
		} else {
			rewritten = append(rewritten, line)
		}
	}

	return strings.Join(rewritten, "\n")
}

// getBaseURL extracts base URL from full URL
func getBaseURL(fullURL string) string {
	lastSlash := strings.LastIndex(fullURL, "/")
	if lastSlash == -1 {
		return fullURL + "/"
	}
	return fullURL[:lastSlash+1]
}
