package service

import (
	"errors"
	"lostfound/dao"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ListCategories(role string, includeDisabled bool) ([]model.ItemType, error) {
	if role != model.RoleSysAdmin {
		includeDisabled = false
	}
	list, err := dao.ListItemTypes(includeDisabled)
	if err != nil {
		logger.Logger.Error("查询分类列表失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return list, nil
}

type CategoryCreateRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=20"`
	Icon    string `json:"icon" binding:"max=200"`
	Sort    uint   `json:"sort"`
	Enabled *bool  `json:"enabled"`
}

func CreateCategory(req *CategoryCreateRequest) (*model.ItemType, error) {
	exist, err := dao.GetItemTypeByName(req.Name)
	if err != nil {
		logger.Logger.Error("查询分类失败", zap.String("name", req.Name), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if exist != nil {
		return nil, errcode.ErrCategoryNameExist
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	itemType := &model.ItemType{
		Name:    req.Name,
		Icon:    req.Icon,
		Sort:    req.Sort,
		Enabled: enabled,
	}
	if err := dao.CreateItemType(itemType); err != nil {
		logger.Logger.Error("创建分类失败", zap.String("name", req.Name), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return itemType, nil
}

func UpdateCategory(categoryID uint, req *CategoryCreateRequest) (*model.ItemType, error) {
	if _, err := dao.GetItemTypeByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryInvalid
		}
		logger.Logger.Error("查询分类失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	exist, err := dao.GetItemTypeByName(req.Name)
	if err != nil {
		logger.Logger.Error("查询分类失败", zap.String("name", req.Name), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if exist != nil && exist.ID != categoryID {
		return nil, errcode.ErrCategoryNameExist
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	fields := map[string]any{
		"name":    req.Name,
		"icon":    req.Icon,
		"sort":    req.Sort,
		"enabled": enabled,
	}
	if err := dao.UpdateItemType(categoryID, fields); err != nil {
		logger.Logger.Error("更新分类失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	count, err := dao.CountItemsByCategory(categoryID)
	if err != nil {
		logger.Logger.Error("查询分类物品数失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	return &model.ItemType{
		ID:        categoryID,
		Name:      req.Name,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Enabled:   enabled,
		ItemCount: count,
	}, nil
}

func DeleteCategory(categoryID uint) error {
	if _, err := dao.GetItemTypeByID(categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrCategoryInvalid
		}
		logger.Logger.Error("查询分类失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return errcode.ErrInternalServer
	}

	count, err := dao.CountItemsByCategory(categoryID)
	if err != nil {
		logger.Logger.Error("查询分类物品数失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if count > 0 {
		return errcode.ErrCategoryHasItem
	}

	if err := dao.DeleteItemType(categoryID); err != nil {
		logger.Logger.Error("删除分类失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
