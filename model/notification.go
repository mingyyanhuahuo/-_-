package model

import "time"

type Notification struct {
	ID          uint      `json:"notificationId"`
	UserID      uint      `gorm:"index" json:"-"`
	Type        string    `gorm:"size:32" json:"type"` // claim_created | claim_approved | claim_rejected | item_approved | item_rejected | comment_created | system
	Title       string    `gorm:"size:100" json:"title"`
	Content     string    `gorm:"size:500" json:"content"`
	RelatedID   *uint     `json:"relatedId"`
	RelatedType string    `gorm:"size:16" json:"relatedType"` // item | claim | comment
	CreatedAt   time.Time `json:"createTime"`
}
