package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// JSON is a custom type for JSON fields
type JSON map[string]interface{}

// Pagination represents pagination parameters
type Pagination struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// GetPage returns the current page (1-indexed)
func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetLimit returns the limit per page
func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

// GetOffset returns the offset for pagination
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}
