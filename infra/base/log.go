package base

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	// 定义日志的格式
	formatter := &prefixed.TextFormatter{}
	//开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	// 控制台高亮
	formatter.ForceFormatting = true
	formatter.ForceColors = true
	formatter.DisableColors = false
	//设置高亮显示的色彩样式
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	log.SetFormatter(formatter)
	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}
	//log.Info("测试Info")
	log.Debug("测试debug")
	// 日志文件
	logFileSettings()
}


func logFileSettings() {

	//配置日志输出目录
	logPath, _ := filepath.Abs("./logs")
	log.Infof("log dir: %s", logPath)
	logFileName := "resk"
	//日志文件最大保存时间，24小时
	maxAge := time.Hour * 24
	//日志切割时间间隔,1小时一个
	rotationTime := time.Hour * 1
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	//设置滚动日志输出
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	log.SetOutput(writer)

}
