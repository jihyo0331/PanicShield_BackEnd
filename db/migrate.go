package db

import (
	"os"
	"ps_backend/model"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrateWithDSN(dsn string) (*gorm.DB, error) {
	logrus.Infof("Starting database migration")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("DB connection failed: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.ChatbotLog{},
		&model.Interest{},
		&model.SubInterest{},
		&model.UserInterest{},
		&model.UserSubInterest{},
		&model.VitalSign{},
		&model.PanicGuide{},
		&model.UserPanicGuide{},
	)
	if err != nil {
		logrus.Fatalf("Migration failed: %v", err)
		return nil, err
	}

	logrus.Infof("Database migration completed")
	return db, nil
}

func GetDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		panic("DATABASE_DSN environment variable not set")
	}
	db, err := ConnectAndMigrateWithDSN(dsn)
	if err != nil {
		panic(err)
	}
	return db
}
