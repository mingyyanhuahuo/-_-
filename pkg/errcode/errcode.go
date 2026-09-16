package errcode

type BizError struct {
	HttpStatus int
	Code       int
	Message    string
}

func (e *BizError) Error() string {
	return e.Message
}

func New(httpStatus int, code int, msg string) *BizError {
	return &BizError{
		HttpStatus: httpStatus,
		Code:       code,
		Message:    msg,
	}
}

var (
	ErrSuccess          = &BizError{200, 200, "success"}
	ErrBadRequest       = &BizError{400, 1001, "参数错误"}
	ErrUnauthorized     = &BizError{401, 1002, "未登录"}
	ErrForbidden        = &BizError{403, 1003, "权限不足"}
	ErrResourceNotFound = &BizError{404, 1004, "资源不存在"}
	ErrTooManyRequests  = &BizError{429, 1005, "请求过于频繁"}
	ErrTokenExpired     = &BizError{401, 1006, "登录已过期"}

	ErrUserExist           = &BizError{409, 2001, "用户名已存在"}
	ErrUserPwdWrong        = &BizError{400, 2002, "用户名或密码错误"}
	ErrOldPwdWrong         = &BizError{400, 2004, "原密码错误"}
	ErrStudentIDRegistered = &BizError{400, 2005, "学号已被注册"}
	ErrPhoneRegistered     = &BizError{400, 2006, "手机号已被注册"}
	ErrEmailRegistered     = &BizError{400, 2007, "邮箱已被注册"}
	ErrNoBindContact       = &BizError{400, 2008, "该用户未绑定手机号与邮箱"}

	ErrLostInfoNotFound   = &BizError{404, 3001, "发布信息不存在"}
	ErrInfoStatusNotAllow = &BizError{409, 3002, "信息状态不允许该操作"}
	ErrClaimReqNotFound   = &BizError{404, 3003, "认领申请不存在"}
	ErrClaimReqHandled    = &BizError{409, 3004, "认领申请已处理"}
	ErrCannotClaimSelf    = &BizError{400, 3005, "不能认领自己发布的物品"}
	ErrHasPendingClaim    = &BizError{409, 3006, "该物品已有待处理的认领申请"}
	ErrCategoryHasItem    = &BizError{409, 3007, "分类下存在物品，无法删除"}
	ErrCategoryNameExist  = &BizError{409, 3008, "分类名称已存在"}
	ErrCategoryInvalid    = &BizError{400, 3012, "分类不存在或已禁用"}

	ErrAnnouncementNotFound       = &BizError{404, 3009, "公告不存在"}
	ErrAnnouncementStatusNotAllow = &BizError{409, 3010, "公告状态不允许该操作"}
	ErrAnnouncementListFailed     = &BizError{500, 3011, "公告列表查询失败"}

	ErrFileTypeNotSupport = &BizError{400, 4001, "文件类型不支持"}
	ErrFileTooLarge       = &BizError{400, 4002, "文件大小超限"}
	ErrFileNotFound       = &BizError{404, 4003, "文件不存在或已删除"}

	ErrInternalServer = &BizError{500, 5000, "服务器内部错误"}
)
