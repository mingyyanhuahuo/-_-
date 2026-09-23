package dao

import (
	"lostfound/model"

	"gorm.io/gorm"
)

func ListUsers(role, keyword string, offset, limit int) ([]model.User, int64, error) {
	builder := func(db *gorm.DB) *gorm.DB {
		tx := db.Model(&model.User{})
		if role != "" {
			tx = tx.Where("role = ?", role)
		}
		if keyword != "" {
			like := "%" + keyword + "%"
			tx = tx.Where("user_name LIKE ? OR nick_name LIKE ? OR student_no LIKE ?", like, like, like)
		}
		return tx
	}

	var total int64
	if err := builder(db).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []model.User
	if err := builder(db).Order("id ASC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
