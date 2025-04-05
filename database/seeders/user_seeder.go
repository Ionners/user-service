package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunUserSeeder(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	user := models.User{
		UUID:        uuid.New(),
		Name:        "Administrator",
		UserName:    "admin",
		Password:    string(password),
		PhoneNumber: "081234567890",
		Email:       "admin.gmail.com",
		RoleID:      constants.Admin,
	}

	err := db.FirstOrCreate(&user, models.User{UserName: user.UserName}).Error
	if err != nil {
		logrus.Errorf("failed to seed user: %v", err)
		panic(err)
	}
	logrus.Infof("user %s seeded successfully", user.UserName)
}
