package model

import (
	"time"

	"gorm.io/gorm"
)

type Claim struct {
	ID            uint           `json:"claimId"`
	ItemID        uint           `gorm:"index" json:"itemId"`
	UserID        uint           `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"applicant"`
	ReviewerID    *uint          `json:"reviewerId"`
	ClaimsReason  string         `gorm:"size:500" json:"claimReason"`
	ProofImages   string         `gorm:"type:text" json:"proofImages"`
	ContactType   string         `gorm:"size:16" json:"-"` // phone | wechat | qq | email
	ContactValue  string         `gorm:"size:100" json:"contactValue"`
	PendingStatus string         `gorm:"size:16;default:pending;index" json:"status"` // pending | approved | rejected | cancelled
	Remark        string         `gorm:"size:200" json:"remark"`
	ReviewTime    *time.Time     `json:"reviewTime"`
	CreatedAt     time.Time      `json:"createTime"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
