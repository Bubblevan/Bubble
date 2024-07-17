package core

import (
	"backend/global"
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = red
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	log := global.Config.Logger

	// 自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		// 自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// 自定义输出格式
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s %s %s \n", log.Prefix, timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m %s\n", log.Prefix, timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

func InitLogger() *logrus.Logger {
	mLog := logrus.New()
	mLog.SetOutput(os.Stdout) // 设置输出类型
	mLog.SetReportCaller(true)
	mLog.SetFormatter(&LogFormatter{})
	mLog.SetLevel(logrus.DebugLevel)
	return mLog
}

func InitDefaultLogger() {
	// 全局log
	logrus.SetOutput(os.Stdout)          // 设置输出类型
	logrus.SetReportCaller(true)         // 开启返回函数名和行号
	logrus.SetFormatter(&LogFormatter{}) // 设置自己定义的Formatter
	logrus.SetLevel(logrus.DebugLevel)   // 设置最低的level
}
