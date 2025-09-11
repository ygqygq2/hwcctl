package logx

import (
	"testing"
)

func TestSetLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected LogLevel
	}{
		{"debug", DEBUG},
		{"info", INFO},
		{"warn", WARN},
		{"error", ERROR},
		{"invalid", INFO}, // 默认级别
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			SetLevel(tt.level)
			if currentLevel != tt.expected {
				t.Errorf("期望日志级别为 %v，实际为 %v", tt.expected, currentLevel)
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	// 测试不同级别的日志输出
	// 这些测试主要确保函数不会 panic
	SetLevel("debug")
	Debug("调试信息")
	Debugf("调试信息: %s", "test")

	Info("信息日志")
	Infof("信息日志: %s", "test")

	Warn("警告日志")
	Warnf("警告日志: %s", "test")

	Error("错误日志")
	Errorf("错误日志: %s", "test")

	Printf("普通输出: %s\n", "test")
	Println("普通输出")

	// 测试不同级别下的过滤
	SetLevel("error")
	Debug("这条调试信息不应该显示")
	Info("这条信息日志不应该显示")
	Error("这条错误日志应该显示")
}
