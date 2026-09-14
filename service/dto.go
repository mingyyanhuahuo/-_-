package service

import (
	"lostfound/model"
	"time"
)

const (
	defaultPageSize = 10
	maxPageSize     = 50
)

// ________________________________________
type UserBrief struct {
	UserId   uint   `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ________________________________________
// announcement
type PageMeta struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type Announcement struct {
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	IsTop     bool       `json:"isTop"`
	PublishAt *time.Time `json:"publishAt"`
}

type AnnouncementsList struct {
	PageMeta
	Announcements []model.Announcement `json:"list"`
}

// ________________________________________
// item and claim
type ItemBrief struct {
	ItemId       uint      `json:"itemId"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	CoverImage   string    `json:"coverImage"`
	CategoryId   uint      `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Location     string    `json:"location"`
	LostTime     time.Time `json:"lostTime"`
	Status       string    `json:"status"`
	Publisher    UserBrief `json:"publisher"`
	ViewCount    uint      `json:"viewCount"`
	ClaimCount   uint      `json:"claimCount"`
	CreateTime   time.Time `json:"createTime"`
}

type ClaimBrief struct {
	ClaimId        uint      `json:"claimId"`
	ItemId         uint      `json:"itemId"`
	ItemTitle      string    `json:"itemTitle"`
	ItemCoverImage string    `json:"itemCoverImage"`
	ClaimReason    string    `json:"claimReason"`
	Status         string    `json:"status"`
	Applicant      UserBrief `json:"applicant"`
	CreateTime     time.Time `json:"createTime"`
}
type ItemDetail struct {
	ItemId         uint           `json:"itemId"`
	Type           string         `json:"type"`
	Title          string         `json:"title"`
	Category       model.ItemType `json:"category"` // 复用 model，ItemCount 已带 gorm:"-"
	Description    string         `json:"description"`
	Location       string         `json:"location"`
	LostTime       time.Time      `json:"lostTime"`
	Images         []string       `json:"images"`
	ContactType    string         `json:"contactType"`
	ContactValue   string         `json:"contactValue"`
	Status         string         `json:"status"`
	Publisher      UserBrief      `json:"publisher"`
	ViewCount      uint           `json:"viewCount"`
	ClaimCount     uint           `json:"claimCount"`
	IsFavorited    bool           `json:"isFavorited"`
	RejectedReason *string        `json:"rejectedReason"`
	CreateTime     time.Time      `json:"createTime"`
	UpdateTime     time.Time      `json:"updateTime"`
}

type ClaimDetail struct {
	ClaimBrief
	ProofImages  []string   `json:"proofImages"`
	ContactValue string     `json:"contactValue"`
	Remark       *string    `json:"remark"`
	ReviewerId   *uint      `json:"reviewerId"`
	ReviewTime   *time.Time `json:"reviewTime"`
}
