package service

import (
	"errors"
	"fmt"
	"lostfound/dao"
	"lostfound/dto"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func toCommentBrief(comment *model.Comment) dto.CommentBrief {
	return dto.CommentBrief{
		CommentId: comment.ID,
		ItemId:    comment.ItemID,
		Content:   comment.Content,
		Author: dto.UserBrief{
			UserId:   comment.AuthorID,
			Nickname: comment.Author.NickName,
			Avatar:   comment.Author.Avatar,
		},
		CreateTime: comment.CreatedAt,
	}
}

func ListComments(itemID uint, page, pageSize int) (*dto.CommentList, error) {
	if _, err := getItemOrFail(itemID); err != nil {
		return nil, err
	}
	page, pageSize = normalizePageNum(page, pageSize)
	comments, total, err := dao.ListComments(itemID, (page-1)*pageSize, pageSize)
	if err != nil {
		logger.Logger.Error("查询留言列表失败", zap.Uint("itemId", itemID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	list := make([]dto.CommentBrief, 0, len(comments))
	for _, comment := range comments {
		list = append(list, toCommentBrief(comment))
	}
	return &dto.CommentList{
		PageMeta: dto.PageMeta{Total: total, Page: page, PageSize: pageSize},
		List:     list,
	}, nil
}

func CreateComment(userID, itemID uint, req *dto.CommentCreateRequest) (*dto.CommentBrief, error) {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return nil, err
	}
	comment := &model.Comment{
		ItemID:   itemID,
		AuthorID: userID,
		Content:  req.Content,
	}
	if err := dao.CreateComment(comment); err != nil {
		logger.Logger.Error("发表留言失败", zap.Uint("itemId", itemID), zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	if err := dao.GenerateNotification([]model.Notification{{
		UserID:      item.AuthorID,
		Type:        model.NotifyCommentCreated,
		Content:     fmt.Sprintf("您的物品《%s》收到新留言", item.Title),
		RelatedID:   &comment.ID,
		RelatedType: model.TargetComment,
	}}); err != nil {
		logger.Logger.Error("发送留言通知失败", zap.Uint("itemId", itemID), zap.Uint("userId", userID), zap.Error(err))
	}
	author, err := dao.GetUserByID(userID)
	if err != nil {
		logger.Logger.Error("获取留言作者失败", zap.Uint("userId", userID), zap.Error(err))
		return nil, errcode.ErrInternalServer
	}
	comment.Author = *author
	brief := toCommentBrief(comment)
	return &brief, nil
}

func DeleteComment(userID, itemID, commentID uint, role string) error {
	item, err := getItemOrFail(itemID)
	if err != nil {
		return err
	}
	comment, err := dao.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrResourceNotFound
		}
		logger.Logger.Error("查询留言失败", zap.Uint("commentId", commentID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	if comment.ItemID != itemID {
		return errcode.ErrResourceNotFound
	}
	if comment.AuthorID != userID && item.AuthorID != userID && !isAdmin(role) {
		return errcode.ErrForbidden
	}
	if err := dao.DeleteComment(commentID); err != nil {
		logger.Logger.Error("删除留言失败", zap.Uint("commentId", commentID), zap.Error(err))
		return errcode.ErrInternalServer
	}
	return nil
}