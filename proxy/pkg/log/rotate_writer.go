package log

import (
	"chatgpt-proxy/pkg/config"
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

func getRotateWriter() io.Writer {
	cnf := config.GetConf()
	writer := &lumberjack.Logger{
		// 文件名
		Filename: cnf.Log.LogPath,
		// 单个文件大小单位MB
		MaxSize: 1,
		// 最多保留文件数
		MaxBackups: 15,
		// 最长保留时间，天
		MaxAge:    7,
		LocalTime: true,
	}
	return writer
}
