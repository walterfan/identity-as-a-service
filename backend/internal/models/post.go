package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string     `gorm:"size:255;not null" json:"title"`
	Slug        string     `gorm:"size:255;not null;unique" json:"slug"`
	Content     string     `gorm:"type:text" json:"content"`
	Excerpt     string     `gorm:"size:500" json:"excerpt"`
	ImageURL    string     `gorm:"size:255" json:"image_url"`
	Published   bool       `gorm:"default:false" json:"published"`
	PublishedAt *time.Time `json:"published_at"`
	Tags        string     `gorm:"size:255" json:"tags"` // Comma-separated tags
	UserID      uint       `json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"-"`
}
