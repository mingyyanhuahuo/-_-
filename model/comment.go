package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `json:"commentId"`
	ItemID    uint           `gorm:"index" json:"itemId"`
	AuthorID  uint           `json:"-"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	Content   string         `gorm:"size:500" json:"content"`
	CreatedAt time.Time      `json:"createTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
