package db

// func ConnectAndMigrateWithDSN(dsn string) (*gorm.DB, error) {
// 	logrus.Infof("Connecting to database")
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		DisableAutomaticPing: true,
// 	})
// 	if err != nil {
// 		logrus.Fatalf("DB connection failed: %v", err)
// 		return nil, err
// 	}
// 	return db, nil
// }
