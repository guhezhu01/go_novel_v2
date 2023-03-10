package errMsg

const (
	SUCCESS = 200
	ERROR   = 500

	// UsernameUsed code=1000...用户模块的错误
	UsernameUsed        = 1001
	PasswordWrong       = 1002
	UserNotExist        = 1003
	TokenExist          = 1004
	TokenRuntime        = 1005
	TokenWrong          = 1006
	TokenTypeWrong      = 1007
	UserNoRight         = 1008
	UpdatePasswordWrong = 1009
)

var codeMsg = map[uint32]string{
	SUCCESS:        "OK",
	ERROR:          "FAIL",
	UsernameUsed:   "用户名已存在！",
	PasswordWrong:  "密码错误",
	UserNotExist:   "用户不存在",
	TokenExist:     "TOKEN不存在,请重新登陆",
	TokenRuntime:   "TOKEN已过期,请重新登陆",
	TokenWrong:     "TOKEN不正确,请重新登陆",
	TokenTypeWrong: "TOKEN格式错误,请重新登陆",
	UserNoRight:    "该用户无权限",
}

func GetErrMsg(code uint32) string {
	return codeMsg[code]
}
