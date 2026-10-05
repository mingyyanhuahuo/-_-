package dao

import (
	"lostfound/model"
)

func AddFavorite(follow *model.Follow) error {
	return db.Where("item_id = ? AND user_id = ?", follow.ItemID, follow.UserID).
		FirstOrCreate(follow).Error
}

func RemoveFavorite(itemID, userID uint) error {
	return db.Where("item_id = ? AND user_id = ?", itemID, userID).
		Delete(&model.Follow{}).Error
}

func ListFavorites(userID uint, offset, limit int) ([]model.Follow, int64, error) {
	var total int64
	if err := db.Model(&model.Follow{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var follows []model.Follow
	if err := db.Model(&model.Follow{}).
		Where("user_id = ?", userID).
		Preload("Item.Category").Preload("Item.Author").
		Order("created_at DESC, id DESC").
		Offset(offset).Limit(limit).
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}
	return follows, total, nil
}