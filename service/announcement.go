package service

import (
	"errors"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func normalizePageNum(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}
func ListAnnouncements(role, status string, page, pageSize int) (*dto.AnnouncementsList, error) {
	publishedOnly := (role != model.RoleSysAdmin)
	if publishedOnly {
		status = model.AnnouncementStatusPublished
	}

	page, pageSize = normalizePageNum(page, pageSize)

	offset := (page - 1) * pageSize
	announcements, total, err := dao.ListAnnouncements(status, publishedOnly, offset, pageSize)
	if err != nil {
		logger.Logger.Error("查询公告列表失败", zap.Error(err))
		return nil, errcode.ErrAnnouncementListFailed
	}
	if announcements == nil {
		announcements = []model.Announcement{}
	}
	pageMeta := dto.PageMeta{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return &dto.AnnouncementsList{
		PageMeta:      pageMeta,
		Announcements: announcements,
	}, nil
}
func GenerateAnnouncement(userID int, req *dto.AnnouncementRequest) (*model.Announcement, error) {
	now := time.Now()
	announcement := &model.Announcement{
		PublisherID:        uint(userID),
		Title:              req.Title,
		Content:            req.Content,
		IsTop:              req.IsTop,
		PublishAt:          req.PublishAt,
		AnnouncementStatus: model.AnnouncementStatusPublished,
	}
	if announcement.Title == "" || announcement.Content == "" {
		return nil, errcode.ErrBadRequest
	}
	if announcement.PublishAt == nil {
		announcement.PublishAt = &now
	} else if announcement.PublishAt.After(now) {
		announcement.AnnouncementStatus = model.AnnouncementStatusDraft
	}

	if err := dao.GenerateAnnouncement(announcement); err != nil {
		logger.Logger.Error("发布公告失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return announcement, nil
}

func GetAnnouncement(id uint, role string) (*model.Announcement, error) {
	announcement, err := dao.GetAnnouncementByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAnnouncementNotFound
		}
		logger.Logger.Error("查询公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if role != model.RoleSysAdmin && announcement.AnnouncementStatus != model.AnnouncementStatusPublished {
		return nil, errcode.ErrAnnouncementNotFound
	}
	return announcement, nil
}

func UpdateAnnouncement(id uint, req *dto.AnnouncementRequest) (*model.Announcement, error) {
	announcement, err := dao.GetAnnouncementByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrAnnouncementNotFound
		}
		logger.Logger.Error("查询公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	now := time.Now()
	publishAt := req.PublishAt
	if publishAt == nil {
		publishAt = &now
	} else if publishAt.After(now) {
		announcement.AnnouncementStatus = model.AnnouncementStatusDraft
	}
	fields := map[string]any{
		"title":      req.Title,
		"content":    req.Content,
		"is_top":     req.IsTop,
		"publish_at": publishAt,
	}
	if announcement.AnnouncementStatus == model.AnnouncementStatusDraft {
		fields["announcement_status"] = model.AnnouncementStatusDraft
	}
	if err := dao.UpdateAnnouncement(id, fields); err != nil {
		logger.Logger.Error("更新公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	announcement.Title = req.Title
	announcement.Content = req.Content
	announcement.IsTop = req.IsTop
	announcement.PublishAt = publishAt
	return announcement, nil
}
func ChangeAnnouncementStatus(id uint, action string) error {
	announcement, err := dao.GetAnnouncementByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrAnnouncementNotFound
		}
		logger.Logger.Error("查询公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return errcode.ErrInternalServer
	}
	var newStatus string
	switch action {
	case "publish":
		if announcement.AnnouncementStatus == model.AnnouncementStatusPublished {
			return errcode.ErrAnnouncementStatusNotAllow
		}
		newStatus = model.AnnouncementStatusPublished
	case "offline":
		if announcement.AnnouncementStatus != model.AnnouncementStatusPublished {
			return errcode.ErrAnnouncementStatusNotAllow
		}
		newStatus = model.AnnouncementStatusOffline
	default:
		return errcode.ErrBadRequest
	}

	if err := dao.UpdateAnnouncement(id, map[string]any{"announcement_status": newStatus}); err != nil {
		logger.Logger.Error("更新公告状态失败", zap.Uint("announcementId", id), zap.Error(err))
		return errcode.ErrInternalServer
	}

	return nil
}
func DeleteAnnouncement(id uint) error {
	_, err := dao.GetAnnouncementByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrAnnouncementNotFound
		}
		logger.Logger.Error("查询公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if err := dao.DeleteAnnouncement(id); err != nil {
		logger.Logger.Error("删除公告失败", zap.Uint("announcementId", id), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
