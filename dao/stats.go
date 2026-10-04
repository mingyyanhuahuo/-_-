package dao

import (
	"lostfound/model"
	"time"
)

type counts struct {
	Total        int64
	LostCount    int64
	FoundCount   int64
	PendingAudit int64
	ClosedCount  int64
	ClaimedCount int64
}

func TotalClaims(OR *model.OverviewResponse) error {
	var c counts
	err := db.Model(&model.Item{}).
		Select(`
            COUNT(*) AS total,
            SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) AS lost_count,
            SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) AS found_count,
            SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS pending_audit,
            SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS closed_count,
            SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) AS claimed_count
        `,
			model.ItemTypeLost, model.ItemTypeFound,
			model.ItemStatusPending, model.ItemStatusClosed, model.ItemStatusClaimed).
		Scan(&c).Error
	if err != nil {
		return err
	}

	OR.TotalClaims = c.Total
	OR.LostCount = c.LostCount
	OR.FoundCount = c.FoundCount
	OR.PendingAudit = c.PendingAudit
	OR.ClosedCount = c.ClosedCount
	OR.ClaimedCount = c.ClaimedCount

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if err := db.Model(&model.Item{}).
		Where("created_at >= ?", start).
		Count(&OR.TodayNewItems).Error; err != nil {
		return err
	}

	if err := db.Model(&model.User{}).
		Count(&OR.UserCount).Error; err != nil {
		return err
	}

	return nil
}

// Category 统计各分类下的物品数量；includeEmpty 为 true 时返回零物品分类
func Category(OR *[]model.CategoryResponse, itemtype string, includeEmpty bool) error {
	q := db.Model(&model.ItemType{}).
		Select("item_types.id AS category_id, item_types.name AS category_name, COUNT(items.id) AS count").
		Joins("LEFT JOIN items ON items.category_id = item_types.id AND items.type = ? AND items.deleted_at IS NULL", itemtype).
		Group("item_types.id, item_types.name")
	if !includeEmpty {
		q = q.Having("COUNT(items.id) > 0")
	}

	var categoryResponses []model.CategoryResponse
	if err := q.Scan(&categoryResponses).Error; err != nil {
		return err
	}
	for _, cr := range categoryResponses {
		*OR = append(*OR, cr)
	}
	return nil
}

func Trend(TR *[]model.TrendResponse, days int) error {
	var trendResponses []model.TrendResponse
	err := db.Model(&model.Item{}).
		Select("DATE(created_at) AS date, SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) AS lost_count, SUM(CASE WHEN type = ? THEN 1 ELSE 0 END) AS found_count",
			model.ItemTypeLost, model.ItemTypeFound).
		Where("created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Scan(&trendResponses).Error
	if err != nil {
		return err
	}

	for _, tr := range trendResponses {
		*TR = append(*TR, tr)
	}
	return nil
}

func ClaimRate(CR *model.ClaimRateResponse, startTime string, endTime string) error {
	var totalClaims int64
	var approvedClaims int64
	var rejectedClaims int64
	var cancelledClaims int64

	err := db.Model(&model.Claim{}).
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Count(&totalClaims).Error
	if err != nil {
		return err
	}

	err = db.Model(&model.Claim{}).
		Where("pending_status = ? AND created_at BETWEEN ? AND ?",
			model.ClaimStatusApproved, startTime, endTime).
		Count(&approvedClaims).Error
	if err != nil {
		return err
	}

	err = db.Model(&model.Claim{}).
		Where("pending_status = ? AND created_at BETWEEN ? AND ?",
			model.ClaimStatusRejected, startTime, endTime).
		Count(&rejectedClaims).Error
	if err != nil {
		return err
	}

	err = db.Model(&model.Claim{}).
		Where("pending_status = ? AND created_at BETWEEN ? AND ?",
			model.ClaimStatusCancelled, startTime, endTime).
		Count(&cancelledClaims).Error
	if err != nil {
		return err
	}

	CR.TotalClaims = totalClaims
	CR.ApprovedClaims = approvedClaims
	CR.RejectedClaims = rejectedClaims
	CR.CancelledClaims = cancelledClaims
	if totalClaims > 0 {
		CR.ClaimRate = float64(approvedClaims) / float64(totalClaims)
	} else {
		CR.ClaimRate = 0
	}
	return nil
}