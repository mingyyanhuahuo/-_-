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
	ErrUnauthorized     = &BizError{401, 1002, "未登录或登录已过期"}
	ErrForbidden        = &BizError{403, 1003, "权限不足"}
	ErrResourceNotFound = &BizError{404, 1004, "资源不存在"}
	ErrTooManyRequests  = &BizError{429, 1005, "请求过于频繁"}

	ErrUserExist           = &BizError{400, 2001, "用户名已存在"}
	ErrUserPwdWrong        = &BizError{400, 2002, "用户名或密码错误"}
	ErrOldPwdWrong         = &BizError{400, 2004, "原密码错误"}
	ErrStudentIDRegistered = &BizError{400, 2005, "学号已被注册"}
	ErrPhoneRegistered     = &BizError{400, 2006, "手机号已被注册"}
	ErrEmailRegistered     = &BizError{400, 2007, "邮箱已被注册"}
	ErrNoBindContact       = &BizError{400, 2008, "该用户未绑定手机号与邮箱"}

	ErrLostInfoNotFound   = &BizError{400, 3001, "发布信息不存在"}
	ErrInfoStatusNotAllow = &BizError{400, 3002, "信息状态不允许该操作"}
	ErrClaimReqNotFound   = &BizError{400, 3003, "认领申请不存在"}
	ErrClaimReqHandled    = &BizError{400, 3004, "认领申请已处理"}
	ErrCannotClaimSelf    = &BizError{400, 3005, "不能认领自己发布的物品"}
	ErrHasPendingClaim    = &BizError{400, 3006, "该物品已有待处理的认领申请"}
	ErrCategoryHasItem    = &BizError{400, 3007, "分类下存在物品，无法删除"}
	ErrCategoryNameExist  = &BizError{400, 3008, "分类名称已存在"}

	ErrFileTypeNotSupport = &BizError{400, 4001, "文件类型不支持"}
	ErrFileTooLarge       = &BizError{400, 4002, "文件大小超限"}
	ErrFileNotFound       = &BizError{400, 4003, "文件不存在或已删除"}

	ErrServerInternal = &BizError{500, 5000, "服务器内部错误"}
)
