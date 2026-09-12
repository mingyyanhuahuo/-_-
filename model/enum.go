package model

// 物品类型
const (
	ItemTypeLost  = "lost"
	ItemTypeFound = "found"
)

// 物品状态
const (
	ItemStatusPending  = "pending"
	ItemStatusApproved = "approved"
	ItemStatusRejected = "rejected"
	ItemStatusClaimed  = "claimed"
	ItemStatusClosed   = "closed"
)

// 认领申请状态
const (
	ClaimStatusPending   = "pending"
	ClaimStatusApproved  = "approved"
	ClaimStatusRejected  = "rejected"
	ClaimStatusCancelled = "cancelled"
)

// 联系方式类型
const (
	ContactTypePhone  = "phone"
	ContactTypeWechat = "wechat"
	ContactTypeQQ     = "qq"
	ContactTypeEmail  = "email"
)

// 用户角色
const (
	RoleStudent  = "student"
	RoleLfAdmin  = "lf_admin"
	RoleSysAdmin = "sys_admin"
)

// 审核动作
const (
	ReviewActionApprove = "approve"
	ReviewActionReject  = "reject"
)

// 公告状态
const (
	AnnouncementStatusDraft     = "draft"
	AnnouncementStatusPublished = "published"
	AnnouncementStatusOffline   = "offline"
)

// 审计日志：操作类型
const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
	ActionLogin  = "LOGIN"
	ActionLogout = "LOGOUT"
	ActionAudit  = "AUDIT"
	ActionExport = "EXPORT"
)

// 审计日志：目标类型
const (
	TargetItem         = "item"
	TargetClaim        = "claim"
	TargetUser         = "user"
	TargetCategory     = "category"
	TargetAnnouncement = "announcement"
	TargetFile         = "file"
)

// 通知类型
const (
	NotifyClaimCreated   = "claim_created"
	NotifyClaimApproved  = "claim_approved"
	NotifyClaimRejected  = "claim_rejected"
	NotifyItemApproved   = "item_approved"
	NotifyItemRejected   = "item_rejected"
	NotifyCommentCreated = "comment_created"
	NotifySystem         = "system"
)
