package db

import (
	"fmt"
	"os"

	"ps_backend/model"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect loads environment variables and returns a GORM DB connection.
// Falls back to default local credentials if environment variables are unset.
func Connect() *gorm.DB {
	// Load .env file if present
	_ = godotenv.Load()

	// Read individual DB connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "psdb"
	}
	if password == "" {
		password = "jmung002"
	}
	if dbname == "" {
		dbname = "ps_db"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	// Open connection with automatic ping disabled
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		logrus.Fatalf("DB connection failed: %v", err)
	}
	return db
}

// Connect is an alias for Connect.
func GetDB() *gorm.DB {
	return Connect()
}

// ConnectAndMigrate connects to the database and runs auto-migrations.
func ConnectAndMigrate() (*gorm.DB, error) {
	db := Connect()
	logrus.Info("Running migrations...")
	if err := db.AutoMigrate(
		&model.User{},
		&model.Interest{},
		&model.SubInterest{},
		&model.UserInterest{},
		&model.UserSubInterest{},
		&model.ChatbotLog{},
		&model.VitalSign{},
		&model.PanicGuide{},
		&model.UserPanicGuide{},
	); err != nil {
		return nil, err
	}
	logrus.Info("Migrations completed.")
	return db, nil
}
