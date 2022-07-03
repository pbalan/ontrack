package models

import (
	"gorm.io/gorm"
	"time"
)

type TokenDetail struct {
	gorm.Model
	AccessToken  string    `gorm:"size:255;unique;not null"`
	RefreshToken string    `gorm:"size:255;unique;not null"`
	AccessUuid   string    `gorm:"size:255;not null"`
	RefreshUuid  string    `gorm:"size:255;not null"`
	AtValidUpto  int64     `gorm:"not null"`
	RtValidUpto  int64     `gorm:"not null"`
	CreatedAt    time.Time `sql:"default:ON UPDATE CURRENT_TIMESTAMP" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `sql:"default:NULL" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
