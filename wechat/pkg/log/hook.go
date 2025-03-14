package log

import (
	nativeLog "log"

	"github.com/sirupsen/logrus"
)

// 为了实现以下三个错误级别的场景下，除了日志除了输出到文件，同步输出到终端
type errorHook struct {
}

func (*errorHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
func (*errorHook) Fire(entry *logrus.Entry) error {
	nativeLog.Println(entry.Message, entry.Data)
	return nil
}
