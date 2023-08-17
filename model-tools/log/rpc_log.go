package log

import (
	"fmt"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Println(v ...any) {
	data := fmt.Sprintln(v...)
	scr, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	logger := logrus.New()
	logger.Out = scr
	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",
		// 清除历史 (WithMaxAge和WithRotationCount只能选其一)
		retalog.WithMaxAge(7*24*time.Hour), //默认每7天清除下日志文件
		// 日志周期(默认每86400秒/一天旋转一次)
		retalog.WithRotationTime(24*time.Hour),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(Hook)

	entry := logger.WithFields(logrus.Fields{
		"data":        data,
		"serviceName": "comment",
	})
	entry.Info()
}
