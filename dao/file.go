package dao

import (
	"lostfound/model"
	"time"
)

func ListExpiredFiles(before time.Time, limit int) ([]model.File, error) {
	var files []model.File
	err := db.Where("created_at < ?", before).Limit(limit).Find(&files).Error
	return files, err
}

func GetFilesByURLs(urls []string) ([]model.File, error) {
	var files []model.File
	if err := db.Where("url IN ?", urls).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

//	func GetFileByID(id uint) (*model.File, error) {
//		var file model.File
//		if err := db.First(&file, id).Error; err != nil {
//			return nil, err
//		}
//		return &file, nil
//	}
func DeleteFile(id uint) error {
	return db.Delete(&model.File{}, id).Error
}
func GenerateFile(file *model.File) error {
	return db.Create(file).Error
}
func CountItemByImage(url string) (int64, error) {
	var count int64
	err := db.Model(&model.Item{}).
		Where("JSON_CONTAINS(images, JSON_QUOTE(?))", url).
		Count(&count).Error
	return count, err
}
func CountClaimByProofImage(url string) (int64, error) {
	var count int64
	err := db.Model(&model.Claim{}).
		Where("JSON_CONTAINS(proof_images, JSON_QUOTE(?))", url).
		Where("pending_status NOT IN ?", []string{model.ClaimStatusCancelled, model.ClaimStatusRejected}).
		Count(&count).Error
	return count, err
}
