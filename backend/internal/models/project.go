package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"size:255" json:"image_url"`
	ProjectURL  string `gorm:"size:255" json:"project_url"`
	GithubURL   string `gorm:"size:255" json:"github_url"`
	Featured    bool   `gorm:"default:false" json:"featured"`
	Tags        string `gorm:"size:255" json:"tags"` // Comma-separated tags
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
}
