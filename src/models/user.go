package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email     *string   `gorm:"size:255;unique;not null"`
	Username  string    `gorm:"size:255;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	FirstName string    `gorm:"size:255;not null"`
	LastName  string    `gorm:"size:255;not null"`
	NickName  string    `gorm:"size:512;default:NULL"`
	CreatedAt time.Time `sql:"default:ON UPDATE CURRENT_TIMESTAMP" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `sql:"default:NULL" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
