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

type RegisterBody struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Password  string `json:"password" binding:"required,min=8,max=20"`
	Nickname  string `json:"nickname" binding:"required,min=2,max=20"`
	StudentNo string `json:"studentNo" binding:"required,len=8"`
	Phone     string `json:"phone" binding:"required,len=11"`
	Email     string `json:"email" binding:"required,email"`
}

type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	UserInfo     UserInfo `json:"userInfo"`
}

type UserInfo struct {
	UserID     uint      `json:"userId"`
	UserName   string    `json:"username" binding:"required,min=3,max=20"`
	NickName   string    `json:"nickname" binding:"required,min=2,max=20"`
	Avatar     string    `gorm:"size:255" json:"avatar"`
	StudentNo  string    `gorm:"size:12;uniqueIndex" json:"studentNo"`
	Phone      string    `gorm:"size:11" json:"phone"`
	Email      string    `gorm:"size:64" json:"email"`
	Role       string    `gorm:"size:16;default:student" json:"role"`
	CreateTime time.Time `json:"createTime"`
}
