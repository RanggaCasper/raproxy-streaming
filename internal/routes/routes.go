package routes

import (
	"raproxy-streaming/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, proxyHandler *handler.ProxyHandler) {
	// Create proxy group
	proxy := app.Group("/proxy")

	// Register proxy endpoints
	proxy.Get("/m3u8", proxyHandler.ProxyM3U8)
	proxy.Get("/segment", proxyHandler.ProxySegment)
	proxy.Get("/video", proxyHandler.ProxyVideo)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
