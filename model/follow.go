package model

import "time"

type Follow struct {
	ID        uint      `json:"-"`
	ItemID    uint      `gorm:"uniqueIndex:uk_fav_user_item;index" json:"itemId"`
	UserID    uint      `gorm:"uniqueIndex:uk_fav_user_item" json:"-"`
	CreatedAt time.Time `json:"createTime"`
}
