package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(100);not null"`
	UserName    string    `gorm:"type:varchar(20);not null"`
	Password    string    `gorm:"type:varchar(255);not null"`
	PhoneNumber string    `gorm:"type:varchar(15);not null"`
	Email       string    `gorm:"type:varchar(100);not null"`
	RoleID      uint      `gorm:"type:uint;not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Role        Role `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
