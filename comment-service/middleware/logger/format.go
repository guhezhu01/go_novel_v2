package logger

import (
	"comment-service/middleware"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logFormatter{
		JSONFormatter: &logrus.JSONFormatter{},
	})
}

type logFormatter struct {
	*logrus.JSONFormatter
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	requestId, ok := middleware.GetRequestID(entry.Context, "comment-service")
	if !ok {
		return nil, nil
	}
	if len(requestId) > 0 {
		entry.Data["comment-service"] = requestId
	}

	return f.JSONFormatter.Format(entry)
}
