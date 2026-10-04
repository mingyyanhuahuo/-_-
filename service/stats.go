package service

import (
	"lostfound/dao"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
)

func Overview() (model.OverviewResponse, error) {
	var OR model.OverviewResponse
	if err := dao.TotalClaims(&OR); err != nil {
		logger.Logger.Error("查询概览统计失败", zap.Error(err))
		return OR, errcode.ErrInternalServer
	}
	return OR, nil
}

func Category(itemType string, includeEmpty bool) ([]model.CategoryResponse, error) {
	var result []model.CategoryResponse
	if itemType != "" {
		if err := dao.Category(&result, itemType, includeEmpty); err != nil {
			logger.Logger.Error("查询分类统计失败", zap.Error(err))
			return result, errcode.ErrInternalServer
		}
	} else {
		if err := dao.Category(&result, model.ItemTypeLost, includeEmpty); err != nil {
			logger.Logger.Error("查询分类统计失败", zap.Error(err))
			return result, errcode.ErrInternalServer
		}
		if err := dao.Category(&result, model.ItemTypeFound, includeEmpty); err != nil {
			logger.Logger.Error("查询分类统计失败", zap.Error(err))
			return result, errcode.ErrInternalServer
		}
	}
	return result, nil
}

func Trend(days int) ([]model.TrendResponse, error) {
	var TR []model.TrendResponse
	if err := dao.Trend(&TR, days); err != nil {
		logger.Logger.Error("查询趋势统计失败", zap.Error(err))
		return TR, errcode.ErrInternalServer
	}
	return TR, nil
}

func ClaimRate(startTime string, endTime string) (model.ClaimRateResponse, error) {
	var CR model.ClaimRateResponse
	if err := dao.ClaimRate(&CR, startTime, endTime); err != nil {
		logger.Logger.Error("查询认领率失败", zap.Error(err))
		return CR, errcode.ErrInternalServer
	}
	return CR, nil
}