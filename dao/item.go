package dao

import (
	"encoding/json"
	"lostfound/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ItemListQuery struct {
	ItemType   string
	CategoryID uint
	Keyword    string
	Location   string
	Statuses   []string
	AuthorID   uint
	StartTime  *time.Time
	EndTime    *time.Time
	SortBy     string
	SortOrder  string
	Offset     int
	Limit      int
}

func GenerateItem(item *model.Item) error {
	return db.Create(item).Error
}
func GetItemByID(id uint) (*model.Item, error) {
	var item model.Item
	err := db.Preload("Category").Preload("Author").First(&item, id).Error
	if err != nil {

		return nil, err
	}
	return &item, nil
}
func UpdateItem(id uint, fields map[string]any) error {
	return db.Model(&model.Item{}).Where("id = ?", id).Updates(fields).Error
}
func DeleteItem(id uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var item model.Item
		if err := tx.First(&item, id).Error; err != nil {
			return err
		}
		var urls []string
		if len(item.Images) > 0 {
			if err := json.Unmarshal(item.Images, &urls); err != nil {
				return err
			}
		}
		if err := tx.Delete(&model.Item{}, id).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Claim{}).
			Where("item_id = ? AND pending_status = ?", id, model.ClaimStatusPending).
			Update("pending_status", model.ClaimStatusCancelled).Error; err != nil {
			return err
		}
		if err := tx.Where("item_id = ?", id).Delete(&model.Follow{}).Error; err != nil {
			return err
		}
		if err := tx.Where("item_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if len(urls) > 0 {
			if err := tx.Where("url IN ?", urls).Delete(&model.File{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func GetItemTypeByID(id uint) (*model.ItemType, error) {
	var itemType model.ItemType
	err := db.First(&itemType, id).Error
	if err != nil {
		return nil, err
	}
	return &itemType, nil
}
func CountItemsByCategory(categoryID uint) (int64, error) {
	var count int64
	err := db.Model(&model.Item{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CountApprovedClaim(itemID, userID uint) (int64, error) {
	var count int64
	err := db.Model(&model.Claim{}).
		Where("item_id = ? AND user_id = ? AND pending_status = ?", itemID, userID, model.ClaimStatusApproved).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func IsFavorited(itemID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&model.Follow{}).
		Where("item_id = ? AND user_id = ?", itemID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func IncViewCount(itemID uint) error {
	return db.Model(&model.Item{}).Where("id = ?", itemID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func itemOrderClause(sortBy, sortOrder string) string {
	column := "created_at"
	if sortBy == "lostTime" {
		column = "lost_time"
	}
	direction := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		direction = "ASC"
	}
	return column + " " + direction + ", id " + direction
}

func ListItems(q ItemListQuery) ([]model.Item, int64, error) {
	builder := func(db *gorm.DB) *gorm.DB {
		tx := db.Model(&model.Item{})
		if q.ItemType != "" {
			tx = tx.Where("type = ?", q.ItemType)
		}
		if q.CategoryID != 0 {
			tx = tx.Where("category_id = ?", q.CategoryID)
		}
		if q.Keyword != "" {
			like := "%" + q.Keyword + "%"
			tx = tx.Where("title LIKE ? OR description LIKE ?", like, like)
		}
		if q.Location != "" {
			tx = tx.Where("location LIKE ?", "%"+q.Location+"%")
		}
		if len(q.Statuses) > 0 {
			tx = tx.Where("status IN ?", q.Statuses)
		}
		if q.StartTime != nil && !q.StartTime.IsZero() {
			tx = tx.Where("lost_time >= ?", *q.StartTime)
		}
		if q.EndTime != nil && !q.EndTime.IsZero() {
			tx = tx.Where("lost_time <= ?", *q.EndTime)
		}
		return tx
	}
	var total int64
	err := builder(db).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	var items []model.Item
	if err := builder(db).Preload("Category").Preload("Author").
		Order(itemOrderClause(q.SortBy, q.SortOrder)).
		Offset(q.Offset).Limit(q.Limit).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

//__________________________________________pend-model
