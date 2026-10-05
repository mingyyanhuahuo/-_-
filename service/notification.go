package service

import (
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
)

func ListNotifications(userID uint, req *dto.NotificationListRequest) (*dto.NotificationList, error) {
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	list, total, err := dao.ListNotifications(userID, req.Type, (page-1)*pageSize, pageSize)
	if err != nil {
		logger.Logger.Error("查询通知列表失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if list == nil {
		list = []model.Notification{}
	}
	return &dto.NotificationList{
		PageMeta: dto.PageMeta{Total: total, Page: page, PageSize: pageSize},
		List:     list,
	}, nil
}