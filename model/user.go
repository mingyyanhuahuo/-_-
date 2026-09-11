package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"userId"`
	UserName  string         `gorm:"size:20;uniqueIndex" json:"username"`
	PassHash  string         `gorm:"size:255" json:"-"`
	NickName  string         `gorm:"size:20" json:"nickname"`
	StudentNo string         `gorm:"size:12;uniqueIndex" json:"studentNo"`
	Phone     string         `gorm:"size:11" json:"phone"`
	Email     string         `gorm:"size:64" json:"email"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:16;default:student" json:"role"` // student | lf_admin | sys_admin
	CreatedAt time.Time      `json:"createTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
