package log

import (
	"fmt"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type logger struct {
	mu sync.Mutex
}

var l *logger

func init() {
	l = new()
}
func new() *logger {
	l := &logger{}
	return l
}

func Println(v ...any) {
	l.Output(fmt.Sprintln(v...))
}

func (l logger) Output(str string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	scr, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
		"data":        str,
		"serviceName": serviceName,
	})
	entry.Info()

}
