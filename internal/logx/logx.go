package logx

import (
	"fmt"
	"log"
	"os"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel = INFO
	debugLogger  = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
	infoLogger   = log.New(os.Stdout, "[INFO]  ", log.LstdFlags)
	warnLogger   = log.New(os.Stdout, "[WARN]  ", log.LstdFlags)
	errorLogger  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)

// SetLevel 设置日志级别
func SetLevel(level string) {
	switch level {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		currentLevel = INFO
	}
}

// Debug 输出调试日志
func Debug(v ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Print(v...)
	}
}

// Debugf 格式化输出调试日志
func Debugf(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Printf(format, v...)
	}
}

// Info 输出信息日志
func Info(v ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Print(v...)
	}
}

// Infof 格式化输出信息日志
func Infof(format string, v ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Printf(format, v...)
	}
}

// Warn 输出警告日志
func Warn(v ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Print(v...)
	}
}

// Warnf 格式化输出警告日志
func Warnf(format string, v ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Printf(format, v...)
	}
}

// Error 输出错误日志
func Error(v ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Print(v...)
	}
}

// Errorf 格式化输出错误日志
func Errorf(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Printf(format, v...)
	}
}

// Fatal 输出致命错误日志并退出
func Fatal(v ...interface{}) {
	errorLogger.Print(v...)
	os.Exit(1)
}

// Fatalf 格式化输出致命错误日志并退出
func Fatalf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
	os.Exit(1)
}

// Printf 普通输出（不带日志级别标识）
func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Println 普通输出（不带日志级别标识）
func Println(v ...interface{}) {
	fmt.Println(v...)
}
