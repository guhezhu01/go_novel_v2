package log

import "log"

var filePath string
var serviceName string

func InitRpcLog(filePathSrc, name string) bool {
	if filePathSrc == "" {
		log.Println("filePath is nil")
		return false
	}
	if name == "" {
		log.Println("serviceName is nil")
		return false
	}
	filePath = filePathSrc
	serviceName = name
	return true
}
