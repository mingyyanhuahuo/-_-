package dao

import (
	"lostfound/model"
	"time"
)

func CreateRefreshToken(userID uint, tokenHash string, expiresAt time.Time) error {
	rt := model.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return db.Create(&rt).Error
}

// GetRefreshTokenByToken 按 token 查 DB;不存在返回 (nil, nil)
func GetRefreshTokenByToken(tokenHash string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := db.Where("token_hash = ?", tokenHash).First(&rt).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken 标记单个 refresh token 为已吊销(用于 rotation / logout)
func RevokeRefreshToken(tokenHash string) error {
	now := time.Now()
	result := db.Model(&model.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RevokeAllRefreshTokensByUser 吊销某用户所有未吊销的 refresh token(改密码时调用)
func RevokeAllRefreshTokensByUser(userID uint) error {
	now := time.Now()
	result := db.Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CleanExpiredRefreshTokens 清理已过期记录,可选调用
func CleanExpiredRefreshTokens() (int64, error) {
	result := db.Where("expires_at < ?", time.Now()).Delete(&model.RefreshToken{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
