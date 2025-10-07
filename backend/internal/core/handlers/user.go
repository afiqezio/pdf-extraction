package handlers

import (
	"net/http"
	"strconv"

	"workbench/internal/core/models"
	"workbench/internal/database"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func (userHandler *UserHandler) UserRoutes(g *echo.Group) {
	users := g.Group("/users")
	users.POST("", userHandler.CreateUser)
	users.GET("", userHandler.GetUsers)
	users.GET("/search", userHandler.SearchUsers)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user.Email == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email, password, first name, and last name are required",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "User with this email already exists",
		})
	}

	// Create user
	if err := h.db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetUsers retrieves all users with pagination
func (h *UserHandler) GetUsers(c echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	pagination := &models.Pagination{
		Page:  page,
		Limit: limit,
	}

	var users []models.User
	var total int64

	// Count total records
	h.db.Model(&models.User{}).Count(&total)

	// Get users with pagination
	if err := h.db.Scopes(database.Paginate(pagination)).Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":        pagination.GetPage(),
			"limit":       pagination.GetLimit(),
			"total":       total,
			"total_pages": (total + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit()),
		},
	})
}

// UpdateUser updates a user by ID
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user",
		})
	}

	// Bind update data
	var updateData models.User
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update user fields (exclude ID and timestamps)
	user.Email = updateData.Email
	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName
	user.IsActive = updateData.IsActive

	// Only update password if provided
	if updateData.Password != "" {
		user.Password = updateData.Password
	}

	if err := h.db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser soft deletes a user by ID
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user",
		})
	}

	// Soft delete
	if err := h.db.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

// SearchUsers searches users by email, first name, or last name
func (h *UserHandler) SearchUsers(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	pagination := &models.Pagination{
		Page:  page,
		Limit: limit,
	}

	var users []models.User
	var total int64

	// Count total records with search
	searchQuery := h.db.Model(&models.User{}).Scopes(
		database.Search(query, "email", "first_name", "last_name"),
	)
	searchQuery.Count(&total)

	// Get users with search and pagination
	if err := searchQuery.Scopes(database.Paginate(pagination)).Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to search users",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":        pagination.GetPage(),
			"limit":       pagination.GetLimit(),
			"total":       total,
			"total_pages": (total + int64(pagination.GetLimit()) - 1) / int64(pagination.GetLimit()),
		},
		"query": query,
	})
}
