package model

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID           uint           `json:"itemId"`
	AuthorID     uint           `gorm:"index" json:"-"`
	Author       User           `gorm:"foreignKey:AuthorID" json:"publisher"`
	CategoryID   uint           `gorm:"index" json:"-"`
	Category     ItemType       `gorm:"foreignKey:CategoryID" json:"category"`
	AdminID      *uint          `json:"-"`
	Admin        User           `gorm:"foreignKey:AdminID" json:"-"`
	Type         string         `gorm:"size:8" json:"type"` // lost | found
	Title        string         `gorm:"size:50" json:"title"`
	Description  string         `gorm:"size:1000" json:"description"`
	Location     string         `gorm:"size:100" json:"location"`
	Images       string         `gorm:"type:text" json:"images"`
	LostTime     time.Time      `json:"lostTime"`
	ContactType  string         `gorm:"size:16" json:"contactType"` // phone | wechat | qq | email
	ContactValue string         `gorm:"size:100" json:"contactValue"`
	Status       string         `gorm:"size:16;default:pending;index" json:"status"` // pending | approved | rejected | claimed | closed
	RejectReason string         `gorm:"size:200" json:"rejectedReason"`
	ClaimCount   uint           `json:"claimCount"`
	ViewCount    uint           `json:"viewCount"`
	ReviewTime   *time.Time     `json:"-"`
	CreatedAt    time.Time      `json:"createTime"`
	UpdatedAt    time.Time      `json:"updateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
