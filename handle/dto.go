package handle

import "time"

type ItemCreateRequest struct {
	Type         string    `json:"type" binding:"required,oneof=lost found"`
	Title        string    `json:"title" binding:"required,min=2,max=50"`
	CategoryId   uint      `json:"categoryId" binding:"required"`
	Description  string    `json:"description" binding:"max=1000"`
	Location     string    `json:"location" binding:"required,min=2,max=100"`
	LostTime     time.Time `json:"lostTime" binding:"required"`
	Images       []string  `json:"images" binding:"max=5,dive,url"`
	ContactType  string    `json:"contactType" binding:"required,oneof=phone wechat qq email"`
	ContactValue string    `json:"contactValue" binding:"required,max=100"`
}

type ItemUpdateRequest struct {
	Title        *string    `json:"title" binding:"omitempty,min=2,max=50"`
	CategoryId   *uint      `json:"categoryId"`
	Description  *string    `json:"description" binding:"omitempty,max=1000"`
	Location     *string    `json:"location" binding:"omitempty,min=2,max=100"`
	LostTime     *time.Time `json:"lostTime"`
	Images       []string   `json:"images" binding:"omitempty,max=5,dive,url"`
	ContactType  *string    `json:"contactType" binding:"omitempty,oneof=phone wechat qq email"`
	ContactValue *string    `json:"contactValue" binding:"omitempty,max=100"`
}

type ClaimCreateRequest struct {
	ItemId       uint     `json:"itemId" binding:"required"`
	ClaimReason  string   `json:"claimReason" binding:"required,min=10,max=500"`
	ProofImages  []string `json:"proofImages" binding:"max=3,dive,url"`
	ContactValue string   `json:"contactValue" binding:"max=100"`
}

type UpdateItemStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Remark string `json:"remark" binding:"max=200"`
}

type AuditRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
	Remark string `json:"remark" binding:"max=200"`
}

type AnnouncementRequest struct {
	Title     string     `json:"title" binding:"required,max=100"`
	Content   string     `json:"content" binding:"required,max=5000"`
	IsTop     bool       `json:"isTop"`
	PublishAt *time.Time `json:"publishAt"`
}
