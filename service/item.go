package service

import (
	"encoding/json"
	"errors"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var minLostTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

var publicItemStatus = []string{
	model.ItemStatusApproved,
	model.ItemStatusClaimed,
	model.ItemStatusClosed,
}
var itemStatusTransition = map[string][]string{
	model.ItemStatusPending:  {model.ItemStatusApproved, model.ItemStatusRejected},
	model.ItemStatusApproved: {model.ItemStatusClosed},
	model.ItemStatusClosed:   {model.ItemStatusApproved},
}

func isAdmin(role string) bool {
	return role == model.RoleLfAdmin || role == model.RoleSysAdmin
}
func itemImageURLs(item *model.Item) []string {
	urls := []string{}
	if len(item.Images) > 0 {
		if err := json.Unmarshal(item.Images, &urls); err != nil {
			return []string{}
		}
		if urls == nil {
			urls = []string{}
		}

	}
	return urls
}
func checkImages(userID uint, urls []string) error {
	if len(urls) == 0 {
		return nil
	}
	uniq := make(map[string]struct{}, len(urls))
	list := make([]string, 0, len(urls))
	for _, url := range urls {
		if _, ok := uniq[url]; ok {
			continue
		}
		uniq[url] = struct{}{}
		list = append(list, url)
	}
	files, err := dao.GetFilesByURLs(list)
	if err != nil {
		logger.Logger.Error("查询文件失败", zap.Uint("userId", userID), zap.Strings("urls", list), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if len(files) != len(list) {
		return errcode.ErrFileNotFound
	}
	for i := range files {
		if files[i].UserID != userID {
			return errcode.ErrFileNotFound
		}
	}
	return nil
}

func checkCategory(categoryID uint) error {
	category, err := dao.GetItemTypeByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrCategoryInvalid
		}
		logger.Logger.Error("查询物品分类失败", zap.Uint("categoryId", categoryID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if !category.Enabled {
		return errcode.ErrCategoryInvalid
	}
	return nil
}
func maskContactValue(contactType, contactValue string) string {
	if contactValue == "" {
		return contactValue
	}
	runes := []rune(contactValue)
	switch contactType {
	case model.ContactTypePhone:
		if len(runes) < 8 {
			return strings.Repeat("*", len(runes))
		}
		return string(runes[:3]) + "****" + string(runes[len(runes)-4:])
	case model.ContactTypeEmail:
		atIndex := strings.Index(contactValue, "@")
		if atIndex <= 1 {
			return strings.Repeat("*", len(runes))
		}
		return string(runes[:1]) + "***" + contactValue[atIndex:]
	default:
		if len(runes) < 4 {
			return strings.Repeat("*", len(runes))
		}
		return string(runes[:2]) + "****" + string(runes[len(runes)-2:])
	}
}
func getItemOrFail(itemID uint) (*model.Item, error) {
	item, err := dao.GetItemByID(itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrLostInfoNotFound
		}
		logger.Logger.Error("查询物品失败", zap.Uint("itemId", itemID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return item, nil
}
func toItemBrief(item *model.Item) dto.ItemBrief {
	coverImage := ""
	if urls := itemImageURLs(item); len(urls) > 0 {
		coverImage = urls[0]
	}
	return dto.ItemBrief{
		ItemId:       item.ID,
		Type:         item.Type,
		Title:        item.Title,
		CoverImage:   coverImage,
		CategoryId:   item.CategoryID,
		CategoryName: item.Category.Name,
		Location:     item.Location,
		LostTime:     item.LostTime,
		Status:       item.Status,
		Publisher: dto.UserBrief{
			UserId:   item.Author.ID,
			Nickname: item.Author.NickName,
			Avatar:   item.Author.Avatar,
		},
		ViewCount:  item.ViewCount,
		ClaimCount: item.ClaimCount,
		CreateTime: item.CreatedAt,
	}
}
func toItemList(items []model.Item, total int64, page, pageSize int) *dto.ItemsList {
	list := make([]dto.ItemBrief, 0, len(items))
	for i := range items {
		list = append(list, toItemBrief(&items[i]))
	}
	return &dto.ItemsList{
		PageMeta: dto.PageMeta{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		Items: list,
	}
}
func GenerateItem(userID uint, req *dto.ItemCreateRequest) (*dto.ItemStatusResponse, error) {

	if req.LostTime.After(time.Now()) {
		return nil, errcode.ErrBadRequest
	}
	if req.LostTime.Before(minLostTime) {
		return nil, errcode.ErrBadRequest
	}
	if err := checkCategory(req.CategoryId); err != nil {
		return nil, err
	}
	if err := checkImages(userID, req.Images); err != nil {
		return nil, err
	}
	images, err := json.Marshal(req.Images)
	if err != nil {
		logger.Logger.Error("序列化图片列表失败", zap.Uint("user_Id", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	item := &model.Item{
		AuthorID:     userID,
		CategoryID:   req.CategoryId,
		Type:         req.Type,
		Title:        req.Title,
		Description:  req.Description,
		Location:     req.Location,
		Images:       images,
		LostTime:     req.LostTime,
		ContactType:  req.ContactType,
		ContactValue: req.ContactValue,
		Status:       model.ItemStatusPending,
	}
	if err := dao.GenerateItem(item); err != nil {
		logger.Logger.Error("发布物品失败", zap.Uint("user_Id", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return &dto.ItemStatusResponse{
		ItemId: item.ID,
		Status: item.Status,
	}, nil

}

func GetItemDetail(userID, itemID uint, role string) (*dto.ItemDetail, error) {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return nil, err
	}

	owner := item.AuthorID == userID
	admin := isAdmin(role)
	if (item.Status == model.ItemStatusPending || item.Status == model.ItemStatusRejected) && !owner && !admin {
		return nil, errcode.ErrLostInfoNotFound
	}
	viewCount := item.ViewCount
	if delta, ok := dao.IncrItemView(itemID); ok && delta > 0 {
		viewCount += uint(delta)
	}
	contactValue := item.ContactValue
	if !owner && !admin {
		full := false
		if userID != 0 {
			count, err := dao.CountApprovedClaim(itemID, userID)
			if err != nil {
				logger.Logger.Error("查询认领失败", zap.Uint("userId", userID), zap.Uint("itemId", itemID), zap.Error(err))
				return nil, errcode.ErrInternalServer
			}
			full = count > 0
		}
		if !full {
			contactValue = maskContactValue(item.ContactType, item.ContactValue)
		}
	}
	favorited := false
	if userID != 0 {
		favorited, err = dao.IsFavorited(itemID, userID)
		if err != nil {
			logger.Logger.Error("查询收藏失败", zap.Uint("userId", userID), zap.Uint("itemId", itemID), zap.Error(err))
			return nil, errcode.ErrInternalServer
		}
	}

	var rejectedReason *string
	var closeRemark *string
	if item.RejectReason != "" {
		rejectedReason = &item.RejectReason
	}
	if item.CloseRemark != "" {
		closeRemark = &item.CloseRemark
	}
	return &dto.ItemDetail{
		ItemId:       item.ID,
		Type:         item.Type,
		Title:        item.Title,
		Category:     item.Category,
		Description:  item.Description,
		Location:     item.Location,
		LostTime:     item.LostTime,
		Images:       itemImageURLs(item),
		ContactType:  item.ContactType,
		ContactValue: contactValue,
		Status:       item.Status,
		Publisher: dto.UserBrief{
			UserId:   item.Author.ID,
			Nickname: item.Author.NickName,
			Avatar:   item.Author.Avatar,
		},
		ViewCount:      viewCount,
		ClaimCount:     item.ClaimCount,
		IsFavorited:    favorited,
		RejectedReason: rejectedReason,
		CloseRemark:    closeRemark,
		CreateTime:     item.CreatedAt,
		UpdateTime:     item.UpdatedAt,
	}, nil
}

func UpdateItem(UserID, itemID uint, req *dto.ItemUpdateRequest) (*dto.ItemStatusResponse, error) {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return nil, err
	}
	if item.AuthorID != UserID {
		return nil, errcode.ErrLostInfoNotFound
	}
	if item.Status == model.ItemStatusClaimed {
		return nil, errcode.ErrInfoStatusNotAllow
	}
	fields := map[string]any{
		"status":        model.ItemStatusPending,
		"reject_reason": "",
		"review_time":   nil,
		"close_remark":  "",
	}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Location != nil {
		fields["location"] = *req.Location
	}
	if req.LostTime != nil {
		if req.LostTime.After(time.Now()) {
			return nil, errcode.ErrBadRequest
		}
		if req.LostTime.Before(minLostTime) {
			return nil, errcode.ErrBadRequest
		}
		fields["lost_time"] = *req.LostTime
	}
	if req.Images != nil {
		if err := checkImages(UserID, req.Images); err != nil {
			return nil, err
		}
		images, err := json.Marshal(req.Images)
		if err != nil {
			logger.Logger.Error("序列化图片列表失败", zap.Uint("user_Id", UserID), zap.Error(err))
			return nil, errcode.ErrInternalServer
		}
		fields["images"] = images
	}
	if req.CategoryId != nil {
		if err := checkCategory(*req.CategoryId); err != nil {
			return nil, err
		}
		fields["category_id"] = *req.CategoryId
	}
	if req.ContactType != nil {
		fields["contact_type"] = *req.ContactType
	}
	if req.ContactValue != nil {
		fields["contact_value"] = *req.ContactValue
	}
	if err := dao.UpdateItem(itemID, fields); err != nil {
		logger.Logger.Error("更新物品失败", zap.Uint("user_Id", UserID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return &dto.ItemStatusResponse{
		ItemId: item.ID,
		Status: model.ItemStatusPending,
	}, nil
}
func DeleteItem(userID, itemID uint, role string) error {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return err
	}
	if item.AuthorID != userID && role != model.RoleSysAdmin {
		return errcode.ErrLostInfoNotFound
	}
	if err := dao.DeleteItem(itemID); err != nil {
		logger.Logger.Error("删除物品失败", zap.Uint("user_Id", userID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
func ListItems(role string, req *dto.ItemListRequest) (*dto.ItemsList, error) {
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	q := dao.ItemListQuery{
		ItemType:   req.Type,
		CategoryID: req.CategoryId,
		Keyword:    req.Keyword,
		Location:   req.Location,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Offset:     (page - 1) * pageSize,
		Limit:      pageSize,
	}
	if isAdmin(role) {
		if req.Status != "" {
			q.Statuses = []string{req.Status}
		}
	} else if slices.Contains(publicItemStatus, req.Status) {
		q.Statuses = []string{req.Status}
	} else {
		q.Statuses = publicItemStatus
	}
	items, total, err := dao.ListItems(q)
	if err != nil {
		logger.Logger.Error("查询物品列表失败", zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return toItemList(items, total, page, pageSize), nil

}
func ListMyItems(userID uint, req *dto.ItemListRequest) (*dto.ItemsList, error) {
	page, pageSize := normalizePageNum(req.Page, req.PageSize)
	q := dao.ItemListQuery{
		ItemType: req.Type,
		AuthorID: userID,
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	}
	if req.Status != "" {
		q.Statuses = []string{req.Status}
	}
	items, total, err := dao.ListItems(q)
	if err != nil {
		logger.Logger.Error("查询我的物品列表失败", zap.Uint("user_Id", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return toItemList(items, total, page, pageSize), nil

}

func UpdateItemStatus(userID, itemID uint, role, status, remark string) error {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return err
	}
	owner := item.AuthorID == userID
	admin := isAdmin(role)
	if !owner && !admin {
		return errcode.ErrLostInfoNotFound
	}
	if !admin && status != model.ItemStatusClosed && status != model.ItemStatusApproved {
		return errcode.ErrForbidden
	}
	if !admin && (item.Status == model.ItemStatusPending || item.Status == model.ItemStatusRejected) {
		return errcode.ErrForbidden
	}
	if !slices.Contains(itemStatusTransition[item.Status], status) {
		return errcode.ErrInfoStatusNotAllow
	}
	var fields = map[string]any{
		"status": status,
	}
	switch {
	case item.Status == model.ItemStatusPending: // pengding -> approved/rejected
		fields["admin_id"] = userID
		fields["review_time"] = time.Now()
		fields["reject_reason"] = ""
		if status == model.ItemStatusRejected {
			fields["reject_reason"] = remark
		}

	default:
		if status == model.ItemStatusClosed {
			fields["close_remark"] = remark
		} else {
			fields["close_remark"] = ""
		}
	}
	if err := dao.UpdateItem(itemID, fields); err != nil {
		logger.Logger.Error("更新物品状态失败", zap.Uint("user_Id", userID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}
