package handler

import (
	"raproxy-streaming/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// ProxyHandler handles proxy HTTP requests
type ProxyHandler struct {
	service *service.ProxyService
}

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(service *service.ProxyService) *ProxyHandler {
	return &ProxyHandler{
		service: service,
	}
}

// ProxyM3U8 handles M3U8 playlist proxy requests
// @Summary Proxy M3U8 playlists
// @Description Proxy M3U8 playlists, rewriting segment URLs to also go through proxy
// @Tags Proxy
// @Param url query string true "M3U8 URL to proxy"
// @Param referer query string false "Referer header"
// @Success 200 {string} string "M3U8 playlist content"
// @Failure 502 {object} map[string]string
// @Failure 504 {object} map[string]string
// @Router /proxy/m3u8 [get]
func (h *ProxyHandler) ProxyM3U8(c *fiber.Ctx) error {
	targetURL := c.Query("url")
	if targetURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "url parameter is required",
		})
	}

	referer := c.Query("referer", "")

	content, contentType, err := h.service.ProxyM3U8(targetURL, referer)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Cache-Control", "no-cache")
	c.Set("Content-Type", contentType)

	return c.Send(content)
}

// ProxySegment handles video segment proxy requests
// @Summary Proxy video segments
// @Description Proxy video segments (ts, mp4 fragments)
// @Tags Proxy
// @Param url query string true "Segment URL to proxy"
// @Param referer query string false "Referer header"
// @Success 200 {string} binary "Video segment data"
// @Failure 502 {object} map[string]string
// @Failure 504 {object} map[string]string
// @Router /proxy/segment [get]
func (h *ProxyHandler) ProxySegment(c *fiber.Ctx) error {
	targetURL := c.Query("url")
	if targetURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "url parameter is required",
		})
	}

	referer := c.Query("referer", "")

	content, contentType, err := h.service.ProxySegment(targetURL, referer)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Cache-Control", "max-age=3600")
	c.Set("Content-Type", contentType)

	return c.Send(content)
}

// ProxyVideo handles direct video proxy requests with streaming
// @Summary Proxy MP4/direct video URLs
// @Description Proxy MP4/direct video URLs with streaming
// @Tags Proxy
// @Param url query string true "Video URL to proxy"
// @Param referer query string false "Referer header"
// @Success 200 {string} binary "Video data"
// @Failure 502 {object} map[string]string
// @Failure 504 {object} map[string]string
// @Router /proxy/video [get]
func (h *ProxyHandler) ProxyVideo(c *fiber.Ctx) error {
	targetURL := c.Query("url")
	if targetURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "url parameter is required",
		})
	}

	referer := c.Query("referer", "")

	resp, err := h.service.ProxyVideo(targetURL, referer)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer fasthttp.ReleaseResponse(resp)

	contentType := string(resp.Header.Peek("Content-Type"))
	if contentType == "" {
		contentType = "video/mp4"
	}

	contentLength := string(resp.Header.Peek("Content-Length"))

	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Content-Type", contentType)
	if contentLength != "" {
		c.Set("Content-Length", contentLength)
	}

	// Copy body to avoid issues with response release
	body := make([]byte, len(resp.Body()))
	copy(body, resp.Body())

	return c.Send(body)
}
