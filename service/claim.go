package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GenerateClaim(userID uint, req *dto.ClaimCreateRequest) (*dto.ClaimStatusResponse, error) {
	item, err := getItemOrFail(req.ItemID)
	if err != nil {
		return nil, err
	}
	if item.AuthorID == userID {
		return nil, errcode.ErrCannotClaimSelf
	}
	if item.Status != model.ItemStatusApproved {
		return nil, errcode.ErrInfoStatusNotAllow
	}
	if err := checkImages(userID, req.ProofImages); err != nil {
		return nil, err
	}
	count, err := dao.CountClaim(req.ItemID, userID, model.ClaimStatusPending)
	if err != nil {
		logger.Logger.Error("统计认领申请失败", zap.Uint("itemID", req.ItemID), zap.Uint("userID", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if count > 0 {
		return nil, errcode.ErrHasPendingClaim
	}
	recentCount, err := dao.CountRecentClaims(userID, req.ItemID, time.Now().Add(-24*time.Hour))
	if err != nil {
		logger.Logger.Error("统计最近认领申请失败", zap.Uint("itemID", req.ItemID), zap.Uint("userID", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if recentCount >= maxClaimPerItemPerDay {
		return nil, errcode.ErrClaimTooFrequent
	}
	proofImages, err := json.Marshal(req.ProofImages)
	if err != nil {
		logger.Logger.Error("序列化认领凭证图片失败", zap.Uint("itemID", req.ItemID), zap.Uint("userID", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	contactValue := req.ContactValue
	if contactValue == "" {
		applicant, err := dao.GetUserByID(userID)
		if err != nil {
			logger.Logger.Error("获取用户信息失败", zap.Uint("userID", userID), zap.Error(err))
			return nil, errcode.ErrInternalServer
		}
		contactValue = applicant.Phone
	}
	claim := &model.Claim{
		ItemID:        req.ItemID,
		UserID:        userID,
		ClaimReason:   req.ClaimReason,
		ProofImages:   proofImages,
		ContactType:   model.ContactTypePhone,
		ContactValue:  contactValue,
		PendingStatus: model.ClaimStatusPending,
	}
	if err := dao.GenerateClaimWithCount(claim); err != nil {
		logger.Logger.Error("创建认领申请失败", zap.Uint("itemID", req.ItemID), zap.Uint("userID", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if err := dao.GenerateNotification([]model.Notification{{
		UserID:      item.AuthorID,
		Type:        model.NotifyClaimCreated,
		Content:     fmt.Sprintf("您的物品《%s》有新的认领申请，请及时处理", item.Title),
		RelatedID:   &claim.ID,
		RelatedType: model.TargetClaim,
	}}); err != nil {
		logger.Logger.Error("发送认领申请通知失败", zap.Uint("itemID", req.ItemID), zap.Uint("userID", userID), zap.Error(err))
	}
	return &dto.ClaimStatusResponse{ClaimId: claim.ID, Status: model.ClaimStatusPending}, nil
}

func claimProofImageURLs(claim *model.Claim) []string {
	urls := []string{}
	if len(claim.ProofImages) > 0 {
		if err := json.Unmarshal(claim.ProofImages, &urls); err != nil {
			logger.Logger.Error("反序列化认领凭证图片失败", zap.Uint("claimID", claim.ID), zap.Error(err))
			return []string{}
		}
		if urls == nil {
			urls = []string{}
		}
	}
	return urls
}

func toClaimBrief(claim *model.Claim) dto.ClaimBrief {
	coverImage := ""
	if urls := itemImageURLs(&claim.Item); len(urls) > 0 {
		coverImage = urls[0]
	}
	return dto.ClaimBrief{
		ClaimId:        claim.ID,
		ItemId:         claim.ItemID,
		ItemTitle:      claim.Item.Title,
		ItemCoverImage: coverImage,
		ClaimReason:    claim.ClaimReason,
		Status:         claim.PendingStatus,
		Applicant: dto.UserBrief{
			UserId:   claim.UserID,
			Nickname: claim.User.NickName,
			Avatar:   claim.User.Avatar,
		},
		CreateTime: claim.CreatedAt,
	}
}
func toClaimDetail(claim *model.Claim, viewerID uint, role string) dto.ClaimDetail {
	contactValue := claim.ContactValue
	if claim.UserID != viewerID && claim.PendingStatus != model.ClaimStatusApproved && !isAdmin(role) {
		contactValue = maskContactValue(model.ContactTypePhone, claim.ContactValue)
	}
	var remark *string
	if claim.Remark != "" {
		remark = &claim.Remark
	}
	claimBrief := toClaimBrief(claim)
	return dto.ClaimDetail{
		ClaimBrief:   claimBrief,
		ProofImages:  claimProofImageURLs(claim),
		ContactValue: contactValue,
		Remark:       remark,
		ReviewerID:   claim.ReviewerID,
		ReviewTime:   claim.ReviewTime,
	}
}
func GetClaimDetail(userID, claimID uint, role string) (*dto.ClaimDetail, error) {
	claim, err := GetClaimOrFail(claimID)
	if err != nil {
		return nil, err
	}
	if claim.UserID != userID && claim.Item.AuthorID != userID && !isAdmin(role) {
		return nil, errcode.ErrClaimReqNotFound
	}
	detail := toClaimDetail(claim, userID, role)
	return &detail, nil
}

func GetClaimOrFail(claimID uint) (*model.Claim, error) {
	claim, err := dao.GetClaimByID(claimID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrClaimReqNotFound
		}
		logger.Logger.Error("查询认领申请失败", zap.Uint("claimID", claimID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return claim, nil
}

func CancelClaim(userID, claimID uint) error {
	claim, err := GetClaimOrFail(claimID)
	if err != nil {
		return err
	}
	if claim.UserID != userID {
		return errcode.ErrClaimReqNotFound
	}
	if claim.PendingStatus != model.ClaimStatusPending {
		return errcode.ErrClaimReqHandled
	}
	if err := dao.CancelClaim(claimID, claim.ItemID); err != nil {
		logger.Logger.Error("取消认领申请失败", zap.Uint("claimID", claimID), zap.Uint("itemID", claim.ItemID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
func toClaimList(claims []*model.Claim, total int64, page, pageSize int) *dto.ClaimsListRep {
	list := make([]dto.ClaimBrief, 0, len(claims))
	for _, claim := range claims {
		list = append(list, toClaimBrief(claim))
	}
	return &dto.ClaimsListRep{
		PageMeta: dto.PageMeta{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		List: list,
	}
}
func ListMyClaims(userID uint, req *dto.ClaimListRequest) (*dto.ClaimsListRep, error) {
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	q := dao.ClaimListQuery{
		UserID: userID,
		Offset: (page - 1) * pageSize,
		Limit:  pageSize,
	}
	if len(req.Status) > 0 {
		q.Statuses = req.Status
	} else {
		q.Statuses = []string{model.ClaimStatusPending, model.ClaimStatusApproved}
	}
	claims, total, err := dao.ListClaim(q)
	if err != nil {
		logger.Logger.Error("查询认领申请列表失败", zap.Uint("userID", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return toClaimList(claims, total, page, pageSize), nil
}
func ListReceivedClaims(userID, itemID uint, role string, req *dto.ClaimListRequest) (*dto.ClaimsListRep, error) {
	item, err := dao.GetItemByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrLostInfoNotFound
		}
		logger.Logger.Error("查询物品信息失败", zap.Uint("itemID", itemID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if item.AuthorID != userID && !isAdmin(role) {
		return nil, errcode.ErrLostInfoNotFound
	}
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	q := dao.ClaimListQuery{
		ItemID: itemID,
		Offset: (page - 1) * pageSize,
		Limit:  pageSize,
	}
	if len(req.Status) > 0 {
		q.Statuses = req.Status
	}
	claims, total, err := dao.ListClaim(q)
	if err != nil {
		logger.Logger.Error("查询收到的认领申请列表失败", zap.Uint("userID", userID), zap.Uint("itemID", itemID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return toClaimList(claims, total, page, pageSize), nil
}
