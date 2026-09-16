package service

import (
	"encoding/json"
	"errors"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var PublicItemStatus = []string{
	model.ItemStatusApproved,
	model.ItemStatusClaimed,
	model.ItemStatusClosed,
}
var itemStatusTransition = map[string][]string{
	model.ItemStatusApproved: {model.ItemStatusClaimed, model.ItemStatusClosed},
	model.ItemStatusClaimed:  {model.ItemStatusApproved},
}

func isAdmin(role string) bool {
	return role == model.RoleLfAdmin || role == model.RoleSysAdmin
}
func itemURLs(item *model.Item) []string {
	urls := []string{}
	if len(item.Images) > 0 {
		if err := item.Images.Unmarshal(item.Images, &urls); err != nil {
			return []string{}
		}
		if urls == nil {
			urls = []string{}
		}

	}
	return urls
}
func checkCategory(categoryID uint) error {
	category, err := dao.GetCategoryByID(categoryID)
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
		return string(runes[:3]) + "****" + string(runes[7:])
	case model.ContactTypeEmail:
		atIndex := strings.Index(contactValue, "@")
		if atIndex <= 1 {
			return strings.Repeat("*", len(runes))
		}
		return string(runes[:1]) + "****" + contactValue[atIndex:]
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
			return nil, errcode.ErrItemNotFound
		}
		logger.Logger.Error("查询物品失败", zap.Uint("itemId", itemID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	return item, nil
}
func toItemBrief(item *model.Item) *dto.ItemBrief {
	coverImage := ""
	if urls := itemURLs(item); len(urls) > 0 {
		coverImage = urls[0]
	}
	return &dto.ItemBrief{
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
			Nickname: item.Author.Nickname,
			Avatar:   item.Author.Avatar,
		},
		ViewCount:  item.ViewCount,
		ClaimCount: item.ClaimCount,
		CreateTime: item.CreatedAt,
	}
}
func toItemList(items []*model.Item, total int64, page, pageSize int) []*dto.ItemsList {
	list := make([]*dto.ItemBrief, len(items))
	for i, item := range items {
		list = append(list, toItemBrief(item))
	}
	return &dto.ItemList{
		PageMeta: dto.PageMeta{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
		Items: list,
	}
}
func GenerateItem(UserID string, req *dto.ItemCreateRequest) (*dto.ItemStatusResponse, error) {
	if req.LostTime.After(time.Now()) {
		return nil, errcode.BadRequest
	}
	if err := checkCategory(req.CategoryId); err != nil {
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

// func GetItemDetail(role) (*dto.ItemDetail, error) {

// }

// func UpateItem(UserID string, itemID uint, req *dto.ItemUpdateRequest) (*dto.ItemStatusResponse, error) {
// 	if req.LostTime != nil && req.LostTime.After(time.Now()) {
// 		return nil, errcode.BadRequest
// 	}
// }
