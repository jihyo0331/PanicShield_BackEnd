package db

import (
	"ps_backend/model"
	"ps_backend/pkg/utils"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	logrus.Info("Starting database seeding...")

	// Seed default admin user
	hashedPassword, _ := utils.HashPassword("admin123") // default admin password
	adminUser := model.User{
		Username:      "admin",
		PasswordHash:  hashedPassword,
		PhoneNumber:   "0000000000",
		Verified:      true,
		SpeakingStyle: "존댓말",
		Tone:          "진지함",
		CreatedAt:     time.Now(),
	}
	if err := db.FirstOrCreate(&adminUser, model.User{Username: "admin"}).Error; err != nil {
		logrus.Errorf("Failed to seed admin user: %v", err)
	} else {
		logrus.Info("Seeded admin user")
	}

	// Seed interests
	interests := []string{"운동", "음악", "독서"}
	for _, interestName := range interests {
		interest := model.Interest{Name: interestName}
		if err := db.FirstOrCreate(&interest, model.Interest{Name: interestName}).Error; err != nil {
			logrus.Errorf("Failed to seed interest %s: %v", interestName, err)
		} else {
			logrus.Infof("Seeded interest %s", interestName)
		}
	}

	// Seed sub-interests
	subInterestsMap := map[string][]string{
		"운동": {"축구", "달리기"},
		"음악": {"피아노", "기타"},
		"독서": {"소설", "시"},
	}

	for parentName, subNames := range subInterestsMap {
		var parentInterest model.Interest
		if err := db.Where("name = ?", parentName).First(&parentInterest).Error; err != nil {
			logrus.Errorf("Failed to find parent interest %s: %v", parentName, err)
			continue
		}
		for _, subName := range subNames {
			subInterest := model.SubInterest{
				Name:       subName,
				InterestID: parentInterest.ID,
			}
			if err := db.FirstOrCreate(&subInterest, model.SubInterest{Name: subName, InterestID: parentInterest.ID}).Error; err != nil {
				logrus.Errorf("Failed to seed sub-interest %s under %s: %v", subName, parentName, err)
			} else {
				logrus.Infof("Seeded sub-interest %s under %s", subName, parentName)
			}
		}
	}

	// Seed panic guides
	panicGuides := []model.PanicGuide{
		{Title: "심호흡", Description: "천천히 깊게 숨을 쉬세요."},
		{Title: "자리 이동", Description: "안전한 곳으로 이동하세요."},
		{Title: "도움 요청", Description: "주변 사람에게 도움을 요청하세요."},
	}

	for _, guide := range panicGuides {
		if err := db.FirstOrCreate(&guide, model.PanicGuide{Title: guide.Title}).Error; err != nil {
			logrus.Errorf("Failed to seed panic guide %s: %v", guide.Title, err)
		} else {
			logrus.Infof("Seeded panic guide %s", guide.Title)
		}
	}

	logrus.Info("Database seeding completed.")
}
