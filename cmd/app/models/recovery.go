package models

import "time"

type Recovery struct {
	ID        uint   `gorm:"primary_key"`
	Email     string `gorm:"type:varchar(100);not null"`
	IP        string `gorm:"type:varchar(100);not null"`
	Domain    string `gorm:"type:varchar(100);not null"`
	Code      string `gorm:"type:varchar(100);not null"`
	AccessID  uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
