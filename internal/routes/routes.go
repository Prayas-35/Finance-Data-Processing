package routes

import (
	"github.com/Prayas-35/Finance-Data-Processing/internal/auth"
	"github.com/Prayas-35/Finance-Data-Processing/internal/handlers"
	"github.com/Prayas-35/Finance-Data-Processing/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	AuthHandler      *handlers.AuthHandler
	UserHandler      *handlers.UserHandler
	RecordHandler    *handlers.RecordHandler
	DashboardHandler *handlers.DashboardHandler
	JWTManager       *auth.JWTManager
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", deps.AuthHandler.Login)

	protected := api.Group("", middleware.RequireAuth(deps.JWTManager))

	users := protected.Group("/users", middleware.RequireRole("admin"))
	users.Post("", deps.UserHandler.Create)
	users.Get("", deps.UserHandler.List)
	users.Patch("/:id", deps.UserHandler.Update)
	users.Patch("/:id/active", deps.UserHandler.SetActive)

	recordsRead := protected.Group("/records", middleware.RequireRole("viewer", "analyst", "admin"))
	recordsRead.Get("", deps.RecordHandler.List)
	recordsRead.Get("/:id", deps.RecordHandler.GetByID)

	recordsWrite := protected.Group("/records", middleware.RequireRole("analyst", "admin"))
	recordsWrite.Post("", deps.RecordHandler.Create)
	recordsWrite.Patch("/:id", deps.RecordHandler.Update)
	recordsWrite.Delete("/:id", deps.RecordHandler.Delete)

	dashboard := protected.Group("/dashboard", middleware.RequireRole("viewer", "analyst", "admin"))
	dashboard.Get("/summary", deps.DashboardHandler.Summary)
	dashboard.Get("/categories", deps.DashboardHandler.Categories)
	dashboard.Get("/trends", deps.DashboardHandler.Trends)
	dashboard.Get("/recent", deps.DashboardHandler.Recent)
}
