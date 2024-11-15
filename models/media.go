package models

import "time"

// Media model
type Media struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"size:100"`
	FileUrl     string `json:"fileUrl" gorm:"not null"`
	FileSize    int64
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `gorm:"many2many:media_tags;"`
}

// MediaTag model (junction table)
type MediaTag struct {
	MediaID   uint   `gorm:"primaryKey;column:media_id"`
	TagID     uint   `gorm:"primaryKey;column:tag_id"`
	Media     *Media `gorm:"foreignKey:MediaID"`
	Tag       *Tag   `gorm:"foreignKey:TagID"`
	CreatedAt time.Time
}
