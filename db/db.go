package db

import (
	"CortexMCP/db/entity"
	"CortexMCP/db/repository"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Database holds the database connection and repositories
type Database struct {
	DB           *gorm.DB
	Repositories *repository.Repositories
}

// NewDatabase creates a new database connection and initializes repositories
func NewDatabase(config Config) (*Database, error) {
	// Create DSN string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Register entities
	if err := entity.RegisterEntities(db); err != nil {
		return nil, fmt.Errorf("failed to register entities: %w", err)
	}

	// Run migrations
	if err := RunMigrations(db.DB()); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create repositories
	repos := repository.NewRepositories(db)

	return &Database{
		DB:           db,
		Repositories: repos,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
