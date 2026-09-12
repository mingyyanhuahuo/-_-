package model

import "time"

type ItemType struct {
	ID        uint      `json:"categoryId"`
	Name      string    `gorm:"size:20;uniqueIndex" json:"name"`
	Icon      string    `gorm:"size:200" json:"icon"`
	Sort      uint      `json:"sort"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	ItemCount int64     `gorm:"-" json:"itemCount"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
