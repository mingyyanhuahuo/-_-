package model

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID                 uint           `json:"announcementId"`
	PublisherID        uint           `gorm:"index" json:"publisherId"`
	Title              string         `gorm:"size:100" json:"title"`
	Content            string         `gorm:"type:text" json:"content"`
	IsTop              bool           `json:"isTop"`
	AnnouncementStatus string         `gorm:"size:16;default:draft;index" json:"status"` // draft | published | offline
	PublishAt          *time.Time     `json:"publishAt"`
	CreatedAt          time.Time      `json:"createTime"`
	UpdatedAt          time.Time      `json:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
