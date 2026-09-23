package dao

import (
	"lostfound/model"

	"gorm.io/gorm"
)

// DeleteUser 软删除用户
func DeleteUser(userID uint) error {
	return db.Delete(&model.User{}, userID).Error
}

// CancelUserPendingClaims 将用户所有待处理认领置为 cancelled
func CancelUserPendingClaims(userID uint) error {
	return db.Model(&model.Claim{}).
		Where("user_id = ? AND pending_status = ?", userID, model.ClaimStatusPending).
		Update("pending_status", model.ClaimStatusCancelled).Error
}

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

// CountUserItems 统计用户发布数(未删除的物品)
func CountUserItems(userID uint) (int64, error) {
	var count int64
	if err := db.Model(&model.Item{}).Where("author_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountUserClaims 统计用户认领数(未删除的认领)
func CountUserClaims(userID uint) (int64, error) {
	var count int64
	if err := db.Model(&model.Claim{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
