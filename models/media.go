package models

import (
	"time"

	"github.com/lib/pq"
)

// Media model
type Media struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null;index:idx_media_name"`
	Description string `json:"description" gorm:"size:100"`
	FileUrl     string `json:"fileUrl" gorm:"not null"`
	FileSize    int64
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `gorm:"many2many:media_tags;"`
}

// MediaTag model (junction table)
type MediaTag struct {
	MediaID   uint   `gorm:"primaryKey;column:media_id;index:idx_media_tags_media_id"`
	TagID     uint   `gorm:"primaryKey;column:tag_id;index:idx_media_tags_media_id"`
	Media     *Media `gorm:"foreignKey:MediaID"`
	Tag       *Tag   `gorm:"foreignKey:TagID"`
	CreatedAt time.Time
}

// Custom model to hold media with just tag names
type MediaWithTagNames struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	FileUrl     string         `json:"fileUrl"`
	TagNames    pq.StringArray `json:"tagNames" gorm:"column:tag_names;type:text"`
}
