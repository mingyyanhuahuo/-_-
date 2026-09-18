package service

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"lostfound/config"
	"lostfound/dao"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func UpLoadDir() string {
	dir := config.GetConfig().Upload.Dir
	if dir == "" {
		dir = "./uploads"
	}
	return dir
}

func UploadFile(userID uint, fileHeader *multipart.FileHeader) (*model.File, error) {
	if fileHeader.Size > maxUploadFileSize {
		return nil, errcode.ErrFileTooLarge
	}
	src, err := fileHeader.Open()
	if err != nil {
		logger.Logger.Error("打开上传文件失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	defer src.Close()
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)

	mimeType := http.DetectContentType(head[:n])
	ext, ok := allowedImageTypes[mimeType]
	if !ok {
		return nil, errcode.ErrFileTypeNotSupport
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		logger.Logger.Error("重置上传文件失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		logger.Logger.Error("生成随机文件名失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	name := hex.EncodeToString(buf) + ext
	dir := UpLoadDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Logger.Error("创建上传目录失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	dst, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		logger.Logger.Error("创建上传文件失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		logger.Logger.Error("写入文件内容失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	file := &model.File{
		UserID:   userID,
		Url:      strings.TrimSuffix(config.GetConfig().Upload.BaseUrl, "/") + "/uploads/" + name,
		Size:     uint(fileHeader.Size),
		FileType: mimeType,
	}
	if err := dao.GenerateFile(file); err != nil {
		logger.Logger.Error("保存文件记录失败", zap.Error(err))
		os.Remove(filepath.Join(dir, name))
		return nil, errcode.ErrInternalServer
	}
	return file, nil
}

// func DeleteFile(userID uint, role string, fileID uint) error {
// 	file, err := dao.GetFileByID(fileID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return errcode.ErrFileNotFound
// 		}
// 		logger.Logger.Error("查询文件记录失败", zap.Uint("userID", userID), zap.Uint("fileID", fileID), zap.Error(err))
// 		return errcode.ErrInternalServer
// 	}
// 	if file.UserID != userID && role != model.RoleSysAdmin {
// 		return errcode.ErrFileNotFound
// 	}
// 	itemCount, err := dao.CountItemByImage(file.Url)
// 	if err != nil {
// 		logger.Logger.Error("统计物品引用失败", zap.Uint("userID", userID), zap.Uint("fileID", fileID), zap.Error(err))
// 		return errcode.ErrInternalServer
// 	}
// 	if itemCount > 0 {
// 		return errcode.ErrInfoStatusNotAllow
// 	}
// 	claimCount, err := dao.CountClaimByProofImage(file.Url)
// 	if err != nil {
// 		logger.Logger.Error("统计认领申请引用失败", zap.Uint("userID", userID), zap.Uint("fileID", fileID), zap.Error(err))
// 		return errcode.ErrInternalServer
// 	}
// 	if claimCount > 0 {
// 		return errcode.ErrInfoStatusNotAllow
// 	}
// 	if err := dao.DeleteFile(fileID); err != nil {
// 		logger.Logger.Error("删除文件记录失败", zap.Uint("userID", userID), zap.Uint("fileID", fileID), zap.Error(err))
// 		return errcode.ErrInternalServer
// 	}
// 	os.Remove(filepath.Join(UpLoadDir(), filepath.Base(file.Url)))
// 	return nil
// }

func CleanOldFiles() (int, error) {
	files, err := dao.ListExpiredFiles(time.Now().Add(-fileCleanGrace), 200)
	if err != nil {
		logger.Logger.Error("查询过期文件失败", zap.Error(err))
		return 0, err
	}
	cleaned := 0
	for _, file := range files {
		itemCount, err := dao.CountItemByImage(file.Url)
		if err != nil {
			logger.Logger.Error("统计物品引用失败", zap.Uint("fileID", file.ID), zap.Error(err))
			continue
		}
		claimCount, err := dao.CountClaimByProofImage(file.Url)
		if err != nil {
			logger.Logger.Error("统计认领申请引用失败", zap.Uint("fileID", file.ID), zap.Error(err))
			continue
		}
		if itemCount > 0 || claimCount > 0 {
			continue
		}
		if err := dao.DeleteFile(file.ID); err != nil {
			logger.Logger.Error("删除文件记录失败", zap.Uint("fileID", file.ID), zap.Error(err))
			continue
		}
		os.Remove(filepath.Join(UpLoadDir(), filepath.Base(file.Url)))
		cleaned++
	}
	return cleaned, nil

}
