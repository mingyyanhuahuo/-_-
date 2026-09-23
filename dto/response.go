package dto

import (
	"lostfound/model"
	"time"
)

type UserBrief struct {
	UserId   uint   `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type PageMeta struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type UserAdminList struct {
	PageMeta
	List []model.User `json:"list"`
}

type AnnouncementsList struct {
	PageMeta
	Announcements []model.Announcement `json:"list"`
}

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
	Category       model.ItemType `json:"category"`
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
	CloseRemark    *string        `json:"closeRemark"`
	CreateTime     time.Time      `json:"createTime"`
	UpdateTime     time.Time      `json:"updateTime"`
}

type ClaimDetail struct {
	ClaimBrief
	ProofImages  []string   `json:"proofImages"`
	ContactValue string     `json:"contactValue"`
	Remark       *string    `json:"remark"`
	ReviewerID   *uint      `json:"reviewerId"`
	ReviewTime   *time.Time `json:"reviewTime"`
}

type UpLoadResponse struct {
	FileId   string `json:"fileId"`
	Url      string `json:"url"`
	Size     uint   `json:"size"`
	MimeType string `json:"mimeType"`
}

type ItemStatusResponse struct {
	ItemId uint   `json:"itemId"`
	Status string `json:"status"`
}
type ItemsList struct {
	PageMeta
	Items []ItemBrief `json:"list"`
}

type ClaimStatusResponse struct {
	ClaimId uint   `json:"claimId"`
	Status  string `json:"status"`
}

type ClaimsListRep struct {
	PageMeta
	List []ClaimBrief `json:"list"`
}
type AuditItemBrief struct {
	ItemId       uint      `json:"itemId"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Images       []string  `json:"images"`
	CategoryId   uint      `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	Location     string    `json:"location"`
	LostTime     time.Time `json:"lostTime"`
	Status       string    `json:"status"`
	Publisher    UserBrief `json:"publisher"`
	CreateTime   time.Time `json:"createTime"`
}
type AuditItemsList struct {
	PageMeta
	List []AuditItemBrief `json:"list"`
}

type AuditClaimBrief struct {
	ClaimId        uint       `json:"claimId"`
	ItemId         uint       `json:"itemId"`
	ItemTitle      string     `json:"itemTitle"`
	ItemCoverImage string     `json:"itemCoverImage"`
	ClaimReason    string     `json:"claimReason"`
	ProofImages    []string   `json:"proofImages"`
	ContactValue   string     `json:"contactValue"`
	Status         string     `json:"status"`
	Applicant      UserBrief  `json:"applicant"`
	Remark         *string    `json:"remark"`
	ReviewerID     *uint      `json:"reviewerId"`
	ReviewTime     *time.Time `json:"reviewTime"`
	CreateTime     time.Time  `json:"createTime"`
}
type AuditClaimsList struct {
	PageMeta
	List []AuditClaimBrief `json:"list"`
}
