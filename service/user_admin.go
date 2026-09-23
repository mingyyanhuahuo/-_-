package service

import (
	"errors"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
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

// GetUser 获取用户详情(含发布数/认领数),仅 sys_admin 调用
func GetUser(userID uint) (*dto.UserAdminDetail, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrResourceNotFound
		}
		logger.Logger.Error("查询用户失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	itemCount, err := dao.CountUserItems(userID)
	if err != nil {
		logger.Logger.Error("统计用户发布数失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	claimCount, err := dao.CountUserClaims(userID)
	if err != nil {
		logger.Logger.Error("统计用户认领数失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}

	return &dto.UserAdminDetail{
		User:       *user,
		ItemCount:  itemCount,
		ClaimCount: claimCount,
	}, nil
}

// DeleteUser 删除用户,不可删除自己。软删除 + 取消其待处理认领
func DeleteUser(operatorID, targetID uint) error {
	if operatorID == targetID {
		return errcode.ErrForbidden
	}

	if _, err := dao.GetUserByID(targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrResourceNotFound
		}
		logger.Logger.Error("查询用户失败", zap.Uint("userId", targetID), zap.Error(err))
		return errcode.ErrInternalServer
	}

	if err := dao.CancelUserPendingClaims(targetID); err != nil {
		logger.Logger.Error("取消用户待处理认领失败", zap.Uint("userId", targetID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if err := dao.DeleteUser(targetID); err != nil {
		logger.Logger.Error("删除用户失败", zap.Uint("userId", targetID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
