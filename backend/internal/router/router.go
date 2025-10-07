package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Setup creates and configures the Echo router
func Setup() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Add your API routes here
	api.GET("/users", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Users endpoint",
		})
	})

	return e
}
