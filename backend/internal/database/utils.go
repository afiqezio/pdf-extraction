package database

import (
	"fmt"
	"workbench/internal/core/models"

	"gorm.io/gorm"
)

// Paginate returns a function that paginates the query
func Paginate(pagination *models.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// page := pagination.GetPage()
		limit := pagination.GetLimit()
		offset := pagination.GetOffset()

		return db.Offset(offset).Limit(limit)
	}
}

// OrderBy returns a function that orders the query
func OrderBy(sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if sort != "" {
			return db.Order(sort)
		}
		return db.Order("created_at DESC")
	}
}

// FilterByTenant returns a function that filters by tenant
func FilterByTenant(tenantID string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if tenantID != "" {
			return db.Where("tenant_id = ?", tenantID)
		}
		return db
	}
}

// Search returns a function that searches in specified fields
func Search(query string, fields ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query == "" || len(fields) == 0 {
			return db
		}

		condition := ""
		values := []interface{}{}

		for i, field := range fields {
			if i > 0 {
				condition += " OR "
			}
			condition += fmt.Sprintf("%s ILIKE ?", field)
			values = append(values, "%"+query+"%")
		}

		return db.Where(condition, values...)
	}
}

// IsActive returns a function that filters only active records
func IsActive() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true)
	}
}

// WithDeleted includes soft deleted records
func WithDeleted() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}
}
