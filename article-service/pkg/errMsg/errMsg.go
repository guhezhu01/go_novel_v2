package errMsg

const (
	SUCCESS = 200
	ERROR   = 500
	// ArtNotExist 小说不存在错误
	ArtNotExist = 2001
	// CateNotExist 分类错误
	CateNotExist = 3002
)

var codeMsg = map[uint32]string{
	SUCCESS:      "OK",
	ERROR:        "FAIL",
	ArtNotExist:  "小说不存在",
	CateNotExist: "该分类不存在",
}

func GetErrMsg(code uint32) string {
	return codeMsg[code]
}
