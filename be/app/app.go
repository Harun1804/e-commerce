package app

import (
	"context"
	"harun1804/e-commerce/configs"
	"harun1804/e-commerce/database/migrations"
	"harun1804/e-commerce/pkg/httpresponse"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

func RunApplication() {
	// Initialize Configuration
	cfg := configs.NewConfig()

	// Initialize the logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Initialize the Fiber application
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return httpresponse.InternalServerError(c, err.Error())
		},
	})

	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		zap.L().Info("request",
			zap.String("ip", c.IP()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Duration("latency", duration),
		)
		return err
	})

	zap.L().Info("Starting the e-commerce application...")

	// Initialize database connection
	err := configs.ConnectionDatabase()
	if err != nil {
		zap.L().Fatal("Failed to connect to the database", zap.Error(err))
		return
	}
	zap.L().Info("Database connection established successfully")

	// Running Migrations
	migrations.RunMigrations()
	zap.L().Info("Database migrations completed successfully")

	// Running Seeders

	// Initialize Minio Setup

	// Initialize Swagger Setup

	// Setup Containers
	container := BuildContainers()
	// Setup Routes
	SetupRoutes(app, container)

	port := cfg.App.Port
	if port == 0 {
		port = 3000 // Default port if not specified in the configuration
	}

	app.Get("/", func(c fiber.Ctx) error {
		return httpresponse.Success(c, "E-Commerce API is running", nil)
	})

	srvErr := make(chan error, 1)
	go func() {
		if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			srvErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		if err != nil {
			zap.L().Error("Server stopped with error", zap.Error(err))
		}
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			zap.L().Fatal("Server forced to shutdown", zap.Error(err))
		}

		configs.CloseDatabase()
		<-ctx.Done()
		zap.L().Info("Server shutdown gracefully")
	}
}
