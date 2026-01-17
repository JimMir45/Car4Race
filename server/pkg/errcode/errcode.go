package errcode

// 业务错误码定义
// 格式：4xxxx，其中 4 表示客户端错误，后四位为具体错误码
const (
	// 通用错误 400xx
	CodeInvalidParam = 40001 // 参数错误/格式不正确

	// 认证错误 400xx
	CodeUnauthorized     = 40005 // 未登录或登录已过期
	CodeInvalidCode      = 40007 // 验证码错误或已过期
	CodePhoneRegistered  = 40008 // 该手机号已注册（保留，当前合并登录不使用）
	CodeInvalidInvite    = 40009 // 邀请码无效或已被使用
	CodeRateLimitExceed  = 40010 // 请求过于频繁
	CodeQueueRequired    = 40011 // 当前访问人数较多，请稍后再试（保留）

	// 下载错误 400xx
	CodeDownloadExpired  = 40003 // 下载链接已过期
	CodeDownloadExceeded = 40004 // 下载次数已用完

	// 权限错误 403xx
	CodeForbidden       = 40301 // 无权限
	CodeAdminRequired   = 40302 // 需要管理员权限
	CodeNotPurchased    = 40303 // 未购买该课程
	CodeAlreadyPurchased = 40304 // 已购买该课程

	// 资源错误 404xx
	CodeNotFound       = 40401 // 资源不存在
	CodeUserNotFound   = 40402 // 用户不存在
	CodeCourseNotFound = 40403 // 课程不存在
)

// 错误码对应的消息
var codeMessages = map[int]string{
	CodeInvalidParam:     "参数错误",
	CodeUnauthorized:     "未登录或登录已过期",
	CodeInvalidCode:      "验证码错误或已过期",
	CodePhoneRegistered:  "该手机号已注册，请直接登录",
	CodeInvalidInvite:    "邀请码无效或已被使用",
	CodeRateLimitExceed:  "请求过于频繁，请稍后再试",
	CodeQueueRequired:    "当前访问人数较多，请稍后再试",
	CodeDownloadExpired:  "下载链接已过期",
	CodeDownloadExceeded: "下载次数已用完，请联系客服",
	CodeForbidden:        "无权限",
	CodeAdminRequired:    "需要管理员权限",
	CodeNotPurchased:     "未购买该课程",
	CodeAlreadyPurchased: "您已购买该课程",
	CodeNotFound:         "资源不存在",
	CodeUserNotFound:     "用户不存在",
	CodeCourseNotFound:   "课程不存在",
}

// Message 获取错误码对应的消息
func Message(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Error 自定义错误类型，包含错误码
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// New 创建带错误码的错误
func New(code int) *Error {
	return &Error{
		Code:    code,
		Message: Message(code),
	}
}

// NewWithMessage 创建带自定义消息的错误
func NewWithMessage(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Is 检查错误是否为指定错误码
func Is(err error, code int) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == code
	}
	return false
}

// GetCode 从错误中获取错误码
func GetCode(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return 0
}
