package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	Success State = "Success"
	Fail    State = "Fail"
)

type State string

const (
	logSavePath = "./log"
	logSaveName = "app.log"
)

var logClient *logrus.Logger
var logFileFullName string

//日志配置
var logConfig = Config{}

type Config struct {
	//日志路径
	LogSavePath string
	//日志文件名字
	LogSaveName string
}

func Init(config *Config) *logrus.Logger {
	if config == nil {
		config = &Config{}
	}
	if len(config.LogSavePath) == 0 {
		config.LogSavePath = logSavePath
	}
	if len(config.LogSaveName) == 0 {
		config.LogSaveName = logSaveName
	}

	logConfig = *config
	logClient = newLogger()
	Infof("App logClient init success;Log path is %s ", logFileFullName)
	return logClient
}

func Info(args ...interface{}) {
	args = append(args, sourced())
	logClient.Infoln(args...)
}
func Infof(format string, args ...interface{}) {
	format += "%s"
	args = append(args, sourced())
	logClient.Infof(format, args...)
}
func Debug(args ...interface{}) {
	args = append(args, sourced())
	logClient.Debugln(args...)
}
func Debugf(format string, args ...interface{}) {
	format += "%s"
	args = append(args, sourced())
	logClient.Debugf(format, args...)
}
func Warn(args ...interface{}) {
	args = append(args, sourced())
	logClient.Warnln(args...)
}
func Warnf(format string, args ...interface{}) {
	format += "%s"
	args = append(args, sourced())
	logClient.Warnf(format, args...)
}
func Error(args ...interface{}) {
	args = append(args, sourced())
	logClient.Errorln(args...)
}
func Errorf(format string, args ...interface{}) {
	format += "%s"
	args = append(args, sourced())
	logClient.Errorf(format, args...)
}

//操作日志模板
//data 补充数据
func ActionLog(state State, module, userName, action, params interface{}, data ...interface{}) {

	if len(data) == 0 {
		Infof("| Action log | | %s | | Module [%s] | user [%s] | action [%s] | result [%s] | params [%+v] ", state, module, userName, action, state, params)
	} else {
		Infof("| Action log | | %s | | Module [%s] | user [%s] | action [%s] | result [%s] | params [%+v] |  %+v", state, module, userName, action, state, params, data)

	}
}
func sourced() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
	}
	return fmt.Sprintf("| source:%s:%d", file, line)
}

func newLogger() *logrus.Logger {
	logClient := logrus.New()
	//禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend) err", err)
	}
	logClient.Out = src
	logClient.SetLevel(logrus.DebugLevel)

	//日志文件路径创建
	exist, _ := PathExists(logConfig.LogSavePath)
	if !exist {
		// 创建文件夹
		err = os.MkdirAll(logConfig.LogSavePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	//日志文件全路径
	logFileFullName = logConfig.LogSavePath + "/" + logConfig.LogSaveName

	logWriter, err := rotatelogs.New(
		logFileFullName+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(logFileFullName),  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.WarnLevel:  logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
	logClient.AddHook(lfHook)

	logger := logrus.New()
	logger.Hooks.Add(lfHook)
	return logger
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
