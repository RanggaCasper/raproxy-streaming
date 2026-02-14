package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"raproxy-streaming/internal/config"
	"raproxy-streaming/internal/handler"
	"raproxy-streaming/internal/httpclient"
	"raproxy-streaming/internal/logger"
	"raproxy-streaming/internal/routes"
	"raproxy-streaming/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize logger
	appLogger := logger.New()
	appLogger.Info("Starting raproxy-streaming server...")

	// Initialize HTTP client
	httpClient := httpclient.New(
		cfg.HTTP.Timeout,
		cfg.HTTP.ConnectTimeout,
		cfg.HTTP.MaxRedirects,
	)

	// Initialize service layer
	proxyService := service.NewProxyService(httpClient, appLogger)

	// Initialize handler layer
	proxyHandler := handler.NewProxyHandler(proxyService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "RaProxy Streaming",
		ServerHeader: "raproxy-streaming",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "*",
	}))

	// Setup routes
	routes.SetupRoutes(app, proxyHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start server in goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	appLogger.Info("Server started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	if err := app.Shutdown(); err != nil {
		appLogger.Error("Server forced to shutdown: %v", err)
	}

	appLogger.Info("Server stopped")
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
