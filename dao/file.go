package dao

import (
	"lostfound/model"
)

func GetFileByID(id uint) (*model.File, error) {
	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
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
		Count(&count).Error
	return count, err
}
