package dao

import (
	"lostfound/model"

	"gorm.io/gorm"
)

func ListComments(itemID uint, offset, limit int) ([]*model.Comment, int64, error) {
	builder := func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&model.Comment{}).Where("item_id = ?", itemID)
	}
	var total int64
	if err := builder(db).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var comments []*model.Comment
	if err := builder(db).Preload("Author").
		Order("created_at ASC, id ASC").
		Offset(offset).Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func CreateComment(comment *model.Comment) error {
	return db.Create(comment).Error
}

func GetCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func DeleteComment(id uint) error {
	return db.Delete(&model.Comment{}, id).Error
}