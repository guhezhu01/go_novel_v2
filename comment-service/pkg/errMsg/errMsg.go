package errMsg

const (
	SUCCESS     = 200
	ERROR       = 500
	TokenFailed = 10001
	// 40xxx评论
	CommentNotExist = 4001
	CommentAddWrong = 4002
)

var codeMsg = map[uint32]string{
	SUCCESS:         "OK",
	ERROR:           "FAIL",
	TokenFailed:     "token验证失败",
	CommentNotExist: "评论不存在",
	CommentAddWrong: "点赞发生错误",
}

func GetErrMsg(code uint32) string {
	return codeMsg[code]
}
