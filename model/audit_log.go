package model

import "time"

type AuditLog struct {
	ID         uint      `json:"logId"`
	UserID     uint      `gorm:"index" json:"userId"`
	UserName   string    `gorm:"size:20" json:"username"`
	ActionType string    `gorm:"size:16" json:"action"`     // CREATE | UPDATE | DELETE | LOGIN | LOGOUT | AUDIT | EXPORT
	TargetType string    `gorm:"size:16" json:"targetType"` // item | claim | user | category | announcement | file
	TargetID   uint      `json:"targetId"`
	Detail     string    `gorm:"type:text" json:"detail"`
	Ip         string    `gorm:"size:64" json:"ip"`
	ActionTime time.Time `gorm:"autoCreateTime" json:"createTime"`
}
