package service

import (
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
)

func ListUsers(role, keyword string, page, pageSize int) (*dto.UserAdminList, error) {
	page, pageSize = normalizePageNum(page, pageSize)
	offset := (page - 1) * pageSize

	users, total, err := dao.ListUsers(role, keyword, offset, pageSize)
	if err != nil {
		logger.Logger.Error("查询用户列表失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if users == nil {
		users = []model.User{}
	}
	return &dto.UserAdminList{
		PageMeta: dto.PageMeta{Total: total, Page: page, PageSize: pageSize},
		List:     users,
	}, nil
}
