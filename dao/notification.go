package dao

import (
	"lostfound/model"

	"gorm.io/gorm"
)

func GenerateNotification(notisfictions []model.Notification) error {
	if len(notisfictions) == 0 {
		return nil
	}
	return db.Create(&notisfictions).Error
}

func ListNotifications(userID uint, notifType string, offset, limit int) ([]model.Notification, int64, error) {
	builder := func(tx *gorm.DB) *gorm.DB {
		q := tx.Model(&model.Notification{}).Where("user_id = ?", userID)
		if notifType != "" {
			q = q.Where("type = ?", notifType)
		}
		return q
	}
	var total int64
	if err := builder(db).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.Notification
	if err := builder(db).Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}