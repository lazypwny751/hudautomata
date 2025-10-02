package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lazypwny751/hudautomata/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes database connection
func Connect() error {
	var err error
	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite" // Default to SQLite for development
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	if dbDriver == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "huduser"),
			getEnv("DB_PASSWORD", ""),
			getEnv("DB_NAME", "hudautomata"),
			getEnv("DB_PORT", "5432"),
		)
		DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	} else {
		// SQLite for development
		dbPath := getEnv("DB_PATH", "./hudautomata.db")
		DB, err = gorm.Open(sqlite.Open(dbPath), gormConfig)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate models
	if err := AutoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.Transaction{},
		&models.SystemLog{},
	)
}

// SeedData creates initial admin user if not exists
func SeedData() error {
	var count int64
	DB.Model(&models.Admin{}).Count(&count)
	
	if count == 0 {
		// Create default super admin
		hashedPassword, err := hashPassword("admin123")
		if err != nil {
			return err
		}

		admin := models.Admin{
			Username:     "admin",
			Email:        "admin@hudautomata.local",
			PasswordHash: hashedPassword,
			Role:         models.RoleSuperAdmin,
			IsActive:     true,
		}

		if err := DB.Create(&admin).Error; err != nil {
			return err
		}

		log.Println("Default admin user created: username=admin, password=admin123")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
