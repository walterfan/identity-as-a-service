package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Name    string `gorm:"size:255;not null" json:"name"`
	Email   string `gorm:"size:255;not null" json:"email"`
	Subject string `gorm:"size:255" json:"subject"`
	Message string `gorm:"type:text;not null" json:"message"`
	Read    bool   `gorm:"default:false" json:"read"`
}
