package service

import (
	"fmt"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/logger"
	"strings"
	"time"

	"lostfound/pkg/errcode"

	"go.uber.org/zap"
)

var auditItemActionStatus = map[string]string{
	model.ReviewActionApprove: model.ItemStatusApproved,
	model.ReviewActionReject:  model.ItemStatusRejected,
}
var auditClaimActionStatus = map[string]string{
	model.ReviewActionApprove: model.ClaimStatusApproved,
	model.ReviewActionReject:  model.ClaimStatusRejected,
}

func toAuditItemBrief(item *model.Item) dto.AuditItemBrief {
	return dto.AuditItemBrief{
		ItemId:       item.ID,
		Type:         item.Type,
		Title:        item.Title,
		Description:  item.Description,
		Images:       itemImageURLs(item),
		CategoryId:   item.CategoryID,
		CategoryName: item.Category.Name,
		Location:     item.Location,
		LostTime:     item.LostTime,
		Status:       item.Status,
		Publisher: dto.UserBrief{
			UserId:   item.AuthorID,
			Nickname: item.Author.NickName,
			Avatar:   item.Author.Avatar,
		},
		CreateTime: item.CreatedAt,
	}
}

func ListAuditItems(role string, req *dto.ItemListRequest) (*dto.AuditItemsList, error) {
	if !isAdmin(role) {
		return nil, errcode.ErrForbidden
	}
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	statuses := []string{model.ItemStatusPending}
	if req.Status != "" {
		statuses = []string{req.Status}
	}
	items, total, err := dao.ListItems(dao.ItemListQuery{
		ItemType:  req.Type,
		Keyword:   req.Keyword,
		Statuses:  statuses,
		SortBy:    "createTime",
		SortOrder: "asc",
		Offset:    (page - 1) * pageSize,
		Limit:     pageSize,
	})
	if err != nil {
		logger.Logger.Error("查询待审核物品列表失败", zap.String("role", role), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	list := make([]dto.AuditItemBrief, 0, len(items))
	for i := range items {
		list = append(list, toAuditItemBrief(&items[i]))
	}
	return &dto.AuditItemsList{
		PageMeta: dto.PageMeta{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		List: list,
	}, nil
}
func AuditItem(role, remark, action string, reviewerID, itemID uint) error {
	if !isAdmin(role) {
		return errcode.ErrForbidden
	}
	status, ok := auditItemActionStatus[action]
	if !ok {
		return errcode.ErrBadRequest
	}
	if action == model.ReviewActionReject && strings.TrimSpace(remark) == "" {
		return errcode.ErrBadRequest
	}
	item, err := getItemOrFail(itemID)
	if err != nil {
		return err
	}
	if item.Status != model.ItemStatusPending {
		return errcode.ErrInfoStatusNotAllow
	}
	if err := UpdateItemStatus(reviewerID, itemID, role, status, remark); err != nil {
		return err
	}
	var notisfyType string
	var content string
	if status == model.ItemStatusRejected {
		notisfyType = model.NotifyItemRejected
		content = fmt.Sprintf("您发布的物品 <%s> 已被管理员审核驳回，原因：%s", item.Title, remark)
	} else {
		notisfyType = model.NotifyItemApproved
		content = fmt.Sprintf("您发布的物品 <%s> 已被管理员审核通过", item.Title)
	}
	if err := dao.GenerateNotification([]model.Notification{{
		UserID:      item.AuthorID,
		Type:        notisfyType,
		Content:     content,
		RelatedID:   &item.ID,
		RelatedType: model.TargetItem,
	}}); err != nil {
		logger.Logger.Error("生成审核通知失败", zap.String("role", role), zap.Uint("itemID", itemID), zap.Error(err))
	}
	return nil
}
func toAuditClaimBrief(claim *model.Claim) dto.AuditClaimBrief {
	coverImage := ""
	imageURLs := itemImageURLs(&claim.Item)
	if len(imageURLs) > 0 {
		coverImage = imageURLs[0]
	}
	var remark *string
	if claim.Remark != "" {
		remark = &claim.Remark
	}
	return dto.AuditClaimBrief{
		ClaimId:        claim.ID,
		ItemId:         claim.ItemID,
		ItemTitle:      claim.Item.Title,
		ItemCoverImage: coverImage,
		ClaimReason:    claim.ClaimReason,
		ProofImages:    claimProofImageURLs(claim),
		ContactValue:   claim.ContactValue,
		Status:         claim.PendingStatus,
		Applicant: dto.UserBrief{
			UserId:   claim.User.ID,
			Nickname: claim.User.NickName,
			Avatar:   claim.User.Avatar,
		},
		Remark:     remark,
		ReviewerID: claim.ReviewerID,
		ReviewTime: claim.ReviewTime,
		CreateTime: claim.CreatedAt,
	}
}
func ListAuditClaims(role string, req *dto.ClaimListRequest) (*dto.AuditClaimsList, error) {
	if !isAdmin(role) {
		return nil, errcode.ErrForbidden
	}
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	statuses := []string{model.ClaimStatusPending}
	if len(req.Status) > 0 {
		statuses = req.Status
	}
	claims, total, err := dao.ListClaim(dao.ClaimListQuery{
		ItemID:   req.ItemID,
		Statuses: statuses,
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	})
	if err != nil {
		logger.Logger.Error("查询待审核申请列表失败", zap.String("role", role), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	list := make([]dto.AuditClaimBrief, 0, len(claims))
	for i := range claims {
		list = append(list, toAuditClaimBrief(claims[i]))
	}
	return &dto.AuditClaimsList{
		PageMeta: dto.PageMeta{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		List: list,
	}, nil
}
func AuditClaim(role, remark, action string, reviewerID, claimID uint) error {

	if !isAdmin(role) {
		return errcode.ErrForbidden
	}
	status, ok := auditClaimActionStatus[action]
	if !ok {
		return errcode.ErrBadRequest
	}
	if action == model.ReviewActionReject && strings.TrimSpace(remark) == "" {
		return errcode.ErrBadRequest
	}
	claim, err := GetClaimOrFail(claimID)
	if err != nil {
		return err
	}
	if claim.PendingStatus != model.ClaimStatusPending {
		return errcode.ErrClaimReqHandled
	}
	switch status {
	case model.ClaimStatusApproved:
		others, err := dao.ApproveClaim(claimID, claim.ItemID, reviewerID, remark)
		if err != nil {
			logger.Logger.Error("审批认领申请失败", zap.String("role", role), zap.Uint("claimID", claimID), zap.Error(err))
			return errcode.ErrInternalServer
		}
		notifications := []model.Notification{{
			UserID:      claim.UserID,
			Type:        model.NotifyClaimApproved,
			Content:     fmt.Sprintf("您对物品 <%s> 的认领申请已被管理员审核通过", claim.Item.Title),
			RelatedID:   &claim.ID,
			RelatedType: model.TargetClaim,
		}}
		for i := range others {
			notifications = append(notifications, model.Notification{
				UserID:      others[i].UserID,
				Type:        model.NotifyClaimRejected,
				Content:     fmt.Sprintf("您对物品 <%s> 的认领申请已被管理员驳回", claim.Item.Title),
				RelatedID:   &others[i].ID,
				RelatedType: model.TargetClaim,
			})
		}
		if err := dao.GenerateNotification(notifications); err != nil {
			logger.Logger.Error("生成认领申请通知失败", zap.String("role", role), zap.Uint("claimID", claimID), zap.Error(err))
		}
	case model.ClaimStatusRejected:
		if err := dao.UpdateClaim(claimID, map[string]any{
			"pending_status": model.ClaimStatusRejected,
			"reviewer_id":    reviewerID,
			"remark":         remark,
			"review_time":    time.Now(),
		}); err != nil {
			logger.Logger.Error("审批认领申请失败", zap.String("role", role), zap.Uint("claimID", claimID), zap.Error(err))
			return errcode.ErrInternalServer
		}
		if err := dao.GenerateNotification([]model.Notification{{
			UserID:      claim.UserID,
			Type:        model.NotifyClaimRejected,
			Content:     fmt.Sprintf("您对物品 <%s> 的认领申请已被管理员驳回", claim.Item.Title),
			RelatedID:   &claim.ID,
			RelatedType: model.TargetClaim,
		}}); err != nil {
			logger.Logger.Error("生成认领申请通知失败", zap.String("role", role), zap.Uint("claimID", claimID), zap.Error(err))
		}
	}
	return nil
}
