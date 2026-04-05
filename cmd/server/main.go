package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/auth"
	"github.com/Prayas-35/Finance-Data-Processing/internal/config"
	"github.com/Prayas-35/Finance-Data-Processing/internal/database"
	"github.com/Prayas-35/Finance-Data-Processing/internal/handlers"
	"github.com/Prayas-35/Finance-Data-Processing/internal/middleware"
	"github.com/Prayas-35/Finance-Data-Processing/internal/repositories"
	"github.com/Prayas-35/Finance-Data-Processing/internal/routes"
	"github.com/Prayas-35/Finance-Data-Processing/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect db", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTokenTTL)

	userRepo := repositories.NewUserRepository(pool)
	recordRepo := repositories.NewRecordRepository(pool)
	dashboardRepo := repositories.NewDashboardRepository(pool)

	authService := services.NewAuthService(userRepo, jwtManager)
	userService := services.NewUserService(userRepo)
	recordService := services.NewRecordService(recordRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	seedDefaultAdmin(userService)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(logger.New())
	app.Use(cors.New())

	routes.Register(app, routes.Dependencies{
		AuthHandler:      handlers.NewAuthHandler(authService),
		UserHandler:      handlers.NewUserHandler(userService),
		RecordHandler:    handlers.NewRecordHandler(recordService),
		DashboardHandler: handlers.NewDashboardHandler(dashboardService),
		JWTManager:       jwtManager,
	})

	go func() {
		if err = app.Listen(":" + cfg.AppPort); err != nil {
			slog.Error("server stopped", "error", err)
		}
	}()

	shutdown(app)
}

func seedDefaultAdmin(users *services.UserService) {
	_ = users.EnsureSeedAdmin(
		context.Background(),
		os.Getenv("DEFAULT_ADMIN_EMAIL"),
		os.Getenv("DEFAULT_ADMIN_NAME"),
		os.Getenv("DEFAULT_ADMIN_PASSWORD"),
	)
}

func shutdown(app *fiber.App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
