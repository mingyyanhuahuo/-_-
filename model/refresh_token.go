package model

import "time"

type RefreshToken struct {
	ID        uint       `json:"-"`
	UserID    uint       `gorm:"index" json:"-"`
	TokenHash string     `gorm:"size:255;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `json:"-"`
	RevokedAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"-"`
}
