package dao

import "lostfound/model"

func CreateItemType(itemType *model.ItemType) error {
	return db.Create(itemType).Error
}

func UpdateItemType(id uint, fields map[string]any) error {
	return db.Model(&model.ItemType{}).Where("id = ?", id).Updates(fields).Error
}

func DeleteItemType(id uint) error {
	return db.Where("id = ?", id).Delete(&model.ItemType{}).Error
}

func GetItemTypeByName(name string) (*model.ItemType, error) {
	var itemType model.ItemType
	err := db.Where("name = ?", name).First(&itemType).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &itemType, nil
}

func ListItemTypes(includeDisabled bool) ([]model.ItemType, error) {
	var list []model.ItemType
	query := db.Model(&model.ItemType{}).Order("sort ASC")
	if !includeDisabled {
		query = query.Where("enabled = ?", true)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}

	for i := range list {
		var count int64
		if err := db.Model(&model.Item{}).Where("category_id = ?", list[i].ID).Count(&count).Error; err != nil {
			return nil, err
		}
		list[i].ItemCount = count
	}
	return list, nil
}
