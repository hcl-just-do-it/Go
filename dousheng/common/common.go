package common

type ErrNo int

const (
	OK              ErrNo = 0 // 成功
	ParamInvalid    ErrNo = 1 // 参数不合法
	UserHasExisted  ErrNo = 2 // 该 Username 已存在
	UserHasDeleted  ErrNo = 3 // 用户已删除
	UserNotExisted  ErrNo = 4 // 用户不存在
	WrongPassword   ErrNo = 5 // 密码错误
	LoginRequired   ErrNo = 6 // 用户未登录
	OperationFailed ErrNo = 7 // 操作失败

)

type Response struct {
	StatusCode ErrNo  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
