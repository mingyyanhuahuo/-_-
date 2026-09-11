package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint           `json:"fileId"`
	UserID    uint           `gorm:"index" json:"-"`
	Url       string         `gorm:"size:255" json:"url"`
	Size      uint           `json:"size"`
	FileType  string         `gorm:"size:64" json:"mimeType"`
	CreatedAt time.Time      `json:"createTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
