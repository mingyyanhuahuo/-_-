package dao

import (
	"lostfound/model"
)

func GenerateNotification(notisfictions []model.Notification) error {
	if len(notisfictions) == 0 {
		return nil
	}
	return db.Create(&notisfictions).Error
}
