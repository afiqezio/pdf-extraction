package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"workbench/internal/config"
	"workbench/internal/core/models"
)

var (
	DB *gorm.DB
)

// Initialize creates database connection and runs migrations
func Initialize(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	// Configure GORM
	gormConfig := &gorm.Config{
		// Set log mode based on environment
		Logger: logger.Default.LogMode(logger.Info),
		// Disable foreign key constraint when migrating
		DisableForeignKeyConstraintWhenMigrating: false,
		// Use singular table names
		NamingStrategy: nil,
		// Current time function
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open database connection
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established")

	// Run migrations
	if err := Migrate(DB); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return DB, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Create UUID extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.EPBEBase{},
		&models.DepthInfo{},
		&models.MetadataInfo{},
		&models.EPBEPetrographyCarbonate{},
		&models.EPBEPetrographyClastic{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Verify tables were created
	var tableCount int64
	err = db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tableCount).Error
	if err != nil {
		return fmt.Errorf("failed to count tables: %w", err)
	}
	log.Printf("🔍 Tables created: %d", tableCount)

	log.Println("✅ Database migrations completed")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck checks if database is accessible
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}
