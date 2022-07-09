package config

import (
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type LogConfig struct {
	Path        string        `yaml:"path"`
	MaxAge      time.Duration `yaml:"maxAge"`
	RotationAge time.Duration `yaml:"rotationAge"`
	Level       string        `yaml:"level"`
}

type Formatter struct {
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelText := strings.ToUpper(entry.Level.String())
	buf := bytes.NewBuffer(make([]byte, 0, 32))
	//time
	buf.WriteString("[")
	buf.WriteString(entry.Time.Format(timeFormat))
	buf.WriteString("] ")

	//level
	buf.WriteString("[")
	buf.WriteString(levelText)
	buf.WriteString("] ")

	//message
	buf.WriteString(entry.Message)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

func InitLog(config *LogConfig) {
	if config == nil {
		panic("log.Config is nil")
	}
	//设置返回文件、行号、调用函数
	logrus.SetReportCaller(true)
	//设置日志格式
	f := &Formatter{}
	logrus.SetFormatter(f)
	//验证日志路径
	// 如果路径文件夹不存在，自动创建
	logPath := config.Path
	if logPath != "" {
		if _, err := os.Stat(logPath); err != nil {
			if os.IsNotExist(err) {
				logPath = path.Clean(logPath)
				logDir := path.Dir(logPath)
				err := os.MkdirAll(logDir, os.ModePerm)
				if err != nil {
					fmt.Printf("get info of log directory with error: %s\n", err)
					panic(err)
				}
			}
		} else {
			//如果文件存在，确认是否是软连接，如果不是软连接需要重命名
			//兼容源日志
			_, err := os.Readlink(logPath)
			if err != nil {
				//如果目标文件存在，rename会直接替换，先检查目标文件是否存在
				distPath := logPath + "." + time.Now().Format("20060102")
				if _, err := os.Stat(distPath); err == nil {
					panic(fmt.Sprintf("distPath(%s) is exist", distPath))
				}
				if err := os.Rename(logPath, distPath); err != nil {
					panic(fmt.Sprintf("rename error, err:%s, logPath:%s, distPath:%s", err, logPath, distPath))
				}
			}
		}
		absPath, _ := filepath.Abs(logPath)
		writer, err := rotatelogs.New(
			absPath+".%Y%m%d",
			rotatelogs.WithLinkName(absPath),
			rotatelogs.WithMaxAge(config.MaxAge),
			rotatelogs.WithRotationTime(config.RotationAge),
		)
		if err != nil {
			logrus.Println("failed to create log writer with error:", err)
			panic(err)
		}
		logrus.SetOutput(writer)
	}
	//设置日志等级
	lvl, err := logrus.ParseLevel(config.Level)
	if err != nil {
		panic(fmt.Sprintf("logLevel error, logLevel:%s", config.Level))
	}
	logrus.SetLevel(lvl)
}
