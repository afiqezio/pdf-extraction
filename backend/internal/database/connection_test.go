package database

import (
	"testing"
	"workbench/internal/config"
	"workbench/internal/core/models"

	"github.com/google/uuid"
)

func TestDatabaseConnection(t *testing.T) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := Initialize(&cfg.Database)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test creating a tenant
	tenant := &models.Tenant{
		Name:      "Test Tenant",
		Subdomain: "test-" + uuid.New().String()[:8],
		Settings:  models.JSON{"theme": "light"},
	}

	if err := db.Create(tenant).Error; err != nil {
		t.Errorf("Failed to create tenant: %v", err)
	}

	// Test creating a user
	user := &models.User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	if err := db.Create(user).Error; err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Clean up
	db.Delete(user)
	db.Delete(tenant)
}
