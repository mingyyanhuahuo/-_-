package dao

import (
	"lostfound/model"
	"time"

	"gorm.io/gorm"
)

func GenerateAnnouncement(announcement *model.Announcement) error {
	return db.Create(announcement).Error
}
func UpdateAnnouncement(id uint, fields map[string]any) error {
	return db.Model(&model.Announcement{}).Where("id = ?", id).Updates(fields).Error
}
func DeleteAnnouncement(id uint) error {
	return db.Delete(&model.Announcement{}, id).Error
}
func PublishSchedulerAnnouncement() (int64, error) {
	res := db.Model(&model.Announcement{}).
		Where("announcement_status = ?", model.AnnouncementStatusDraft).
		Where("publish_at <= ? AND publish_at IS NOT NULL", time.Now()).
		Update("announcement_status", model.AnnouncementStatusPublished)
	return res.RowsAffected, res.Error
}
func GetAnnouncementByID(id uint) (*model.Announcement, error) {
	var announcement model.Announcement
	err := db.First(&announcement, id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}
func ListAnnouncements(status string, publishedOnly bool, offset, limit int) ([]model.Announcement, int64, error) {
	builder := func() *gorm.DB {
		q := db.Model(&model.Announcement{})
		if status != "" {
			q = q.Where("announcement_status = ?", status)
		}
		if publishedOnly {
			q = q.Where("publish_at <= ? AND publish_at IS NOT NULL", time.Now())
		}
		return q
	}
	var total int64
	if err := builder().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var announcements []model.Announcement
	if err := builder().Order("is_top DESC, publish_at DESC,id DESC").
		Offset(offset).Limit(limit).Find(&announcements).Error; err != nil {
		return nil, 0, err
	}
	return announcements, total, nil
}
