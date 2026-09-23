package dao

import (
	"lostfound/model"
	"lostfound/pkg/errcode"
	"time"

	"gorm.io/gorm"
)

type ClaimListQuery struct {
	UserID   uint
	ItemID   uint
	Statuses []string
	Offset   int
	Limit    int
}

func GenerateClaim(claim *model.Claim) error {
	return db.Create(claim).Error
}

func GetClaimByID(id uint) (*model.Claim, error) {
	var claim model.Claim
	err := db.Preload("Item").Preload("User").First(&claim, id).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func CountClaim(itemID, userID uint, status string) (int64, error) {
	var count int64
	err := db.Model(&model.Claim{}).
		Where("item_id = ? AND user_id = ? AND pending_status = ?", itemID, userID, status).
		Count(&count).Error
	return count, err
}

func UpdateClaim(claimID uint, field map[string]any) error {
	return db.Model(&model.Claim{}).Where("id = ?", claimID).Updates(field).Error
}

func DelClaimCount(itemID uint) error {
	return db.Model(&model.Item{}).
		Where("id = ? AND claim_count > 0", itemID).
		UpdateColumn("claim_count", gorm.Expr("claim_count - ?", 1)).Error
}
func IncrItemClaimCount(itemID uint) error {
	return db.Model(&model.Item{}).
		Where("id = ?", itemID).
		UpdateColumn("claim_count", gorm.Expr("claim_count + ?", 1)).Error
}

func ListClaim(query ClaimListQuery) ([]*model.Claim, int64, error) {
	builder := func(db *gorm.DB) *gorm.DB {
		tx := db.Model(&model.Claim{})
		if query.UserID != 0 {
			tx = tx.Where("user_id = ?", query.UserID)
		}
		if query.ItemID != 0 {
			tx = tx.Where("item_id = ?", query.ItemID)
		}
		if len(query.Statuses) > 0 {
			tx = tx.Where("pending_status IN ?", query.Statuses)
		}
		return tx
	}
	var total int64
	if err := builder(db).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var claims []*model.Claim
	if err := builder(db).Preload("Item").
		Preload("User").Order("created_at DESC,id DESC").
		Offset(query.Offset).Limit(query.Limit).
		Find(&claims).Error; err != nil {
		return nil, 0, err
	}
	return claims, total, nil
}

func ApproveClaim(claimID, itemID, reviewerID uint, remark string) ([]model.Claim, error) {
	var others []model.Claim
	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		if err := tx.Model(&model.Claim{}).
			Where("id = ? AND pending_status = ?", claimID, model.ClaimStatusPending).
			Updates(map[string]any{
				"pending_status": model.ClaimStatusApproved,
				"reviewer_id":    reviewerID,
				"remark":         remark,
				"review_time":    now,
			}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Item{}).
			Where("id = ? AND pending_status = ?", itemID, model.ItemStatusPending).
			Update("status", model.ItemStatusClaimed).Error; err != nil {
			return err
		}
		if err := tx.Where("item_id = ? AND pending_status = ? AND id != ?",
			itemID, model.ClaimStatusPending, claimID).Find(&others).Error; err != nil {
			return err
		}
		if len(others) == 0 {
			return nil
		}
		ids := make([]uint, 0, len(others))
		for i := range others {
			ids = append(ids, others[i].ID)
		}
		return tx.Model(&model.Claim{}).
			Where("id IN ?", ids).
			Updates(map[string]any{
				"pending_status": model.ClaimStatusRejected,
				"reviewer_id":    reviewerID,
				"review_time":    now,
				"remark":         "该物品已被其他申请认领",
			}).Error
	})
	if err != nil {
		return nil, err
	}
	return others, nil
}

func CountRecentClaims(userID, itemID uint, since time.Time) (int64, error) {
	var count int64
	err := db.Model(&model.Claim{}).
		Where("user_id = ? AND item_id = ? AND created_at >= ?", userID, itemID, since).
		Count(&count).Error
	return count, err
}
func CancelClaim(claimID, itemID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Claim{}).
			Where("id = ? AND pending_status = ?", claimID, model.ClaimStatusPending).
			Update("pending_status", model.ClaimStatusCancelled)
		if res.RowsAffected == 0 {
			return errcode.ErrClaimReqHandled
		}
		if err := tx.Model(&model.Claim{}).
			Where("id = ?", claimID).
			Update("pending_status", model.ClaimStatusCancelled).Error; err != nil {
			return err
		}
		return tx.Model(&model.Item{}).
			Where("id = ? AND claim_count > 0", itemID).
			UpdateColumn("claim_count", gorm.Expr("claim_count - ?", 1)).Error
	})
}

func GenerateClaimWithCount(claim *model.Claim) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(claim).Error; err != nil {
			return err
		}
		return tx.Model(&model.Item{}).
			Where("id = ?", claim.ItemID).
			UpdateColumn("claim_count", gorm.Expr("claim_count + ?", 1)).Error
	})
}
