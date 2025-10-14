package router

import (
	"net/http"

	"workbench/internal/core/handlers"
	"workbench/internal/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Setup creates and configures the Echo router
func Setup() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := e.Group("/api/v1")

	getDB := database.GetDB()

	// Initialize handlers here
	userHandler := handlers.NewUserHandler(getDB)
	petrographyClasticHandler := handlers.NewPetrographyClasticHandler(getDB)
	petrographyCarbonateHandler := handlers.NewPetrographyCarbonateHandler(getDB)
	extractionHandler := handlers.NewExtractionHandler(getDB)

	// Add Routes here
	userHandler.UserRoutes(api)
	petrographyClasticHandler.PetrographyClasticRoutes(api)
	petrographyCarbonateHandler.PetrographyCarbonateRoutes(api)
	extractionHandler.ExtractionRoutes(api)

	return e
}
