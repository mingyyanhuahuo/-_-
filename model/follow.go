package model

import "time"

type Follow struct {
	ID        uint      `json:"-"`
	UserID    uint      `gorm:"uniqueIndex:uk_fav_user_item" json:"-"`
	ItemID    uint      `gorm:"uniqueIndex:uk_fav_user_item;index" json:"itemId"`
	CreatedAt time.Time `json:"createTime"`
}
