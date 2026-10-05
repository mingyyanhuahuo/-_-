package service

import (
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
)

func AddFavorite(userID, itemID uint) error {
	if _, err := getItemOrFail(itemID); err != nil {
		return err
	}
	if err := dao.AddFavorite(&model.Follow{UserID: userID, ItemID: itemID}); err != nil {
		logger.Logger.Error("收藏失败", zap.Uint("itemId", itemID), zap.Uint("userId", userID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}

func RemoveFavorite(userID, itemID uint) error {
	if _, err := getItemOrFail(itemID); err != nil {
		return err
	}
	if err := dao.RemoveFavorite(itemID, userID); err != nil {
		logger.Logger.Error("取消收藏失败", zap.Uint("itemId", itemID), zap.Uint("userId", userID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}

func ListFavorites(userID uint, page, pageSize int) (*dto.ItemsList, error) {
	page, pageSize = normalizePageNum(page, pageSize)
	follows, total, err := dao.ListFavorites(userID, (page-1)*pageSize, pageSize)
	if err != nil {
		logger.Logger.Error("查询收藏列表失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	list := make([]dto.ItemBrief, 0, len(follows))
	for i := range follows {
		list = append(list, toItemBrief(&follows[i].Item))
	}
	return &dto.ItemsList{
		PageMeta: dto.PageMeta{Total: total, Page: page, PageSize: pageSize},
		Items:    list,
	}, nil
}