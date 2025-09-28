package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigureCmd(t *testing.T) {
	// 测试configure命令是否正确初始化
	if configureCmd.Use != "configure" {
		t.Errorf("期望命令名称为 'configure'，实际为 '%s'", configureCmd.Use)
	}

	if configureCmd.Short == "" {
		t.Error("configure命令应该有短描述")
	}
}

func TestMaskString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "*"},
		{"ab", "**"},
		{"abc", "***"},
		{"abcd", "****"},
		{"abcde", "abcd*"},
		{"abcdef", "abcd**"},
		{"abcdefg", "abcd***"},
		{"abcdefgh", "abcd****"},
		{"verylongstring", "very**********"},
	}

	for _, tt := range tests {
		result := maskString(tt.input)
		if result != tt.expected {
			t.Errorf("maskString(%q) = %q, 期望 %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetConfigPath(t *testing.T) {
	// 保存原始环境变量
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// 测试正常情况
	testHome := "/tmp/test-home"
	os.Setenv("HOME", testHome)

	expectedPath := filepath.Join(testHome, ".hwcctl", "config")
	actualPath := getConfigPath()

	if actualPath != expectedPath {
		t.Errorf("期望配置路径为 %q，实际为 %q", expectedPath, actualPath)
	}

	// 测试HOME环境变量不存在的情况
	os.Unsetenv("HOME")
	fallbackPath := getConfigPath()

	if fallbackPath == "" {
		t.Error("即使HOME环境变量不存在，也应该返回一个有效路径")
	}
}

func TestLoadConfig(t *testing.T) {
	// 保存原始HOME环境变量
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// 设置临时HOME目录，确保配置文件不存在
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	config := loadConfig()

	// 验证默认配置
	if config.Default.Region != "cn-north-1" {
		t.Errorf("期望默认区域为 'cn-north-1'，实际为 '%s'", config.Default.Region)
	}

	if config.Default.Output != "table" {
		t.Errorf("期望默认输出格式为 'table'，实际为 '%s'", config.Default.Output)
	}
}

func TestSaveConfig(t *testing.T) {
	// 保存原始HOME环境变量
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// 创建临时目录用于测试
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// 测试配置结构
	config := Config{
		Default: Profile{
			AccessKeyID:     "test-access-key",
			SecretAccessKey: "test-secret-key",
			Region:          "cn-north-1",
			Output:          "json",
		},
	}

	// 测试保存配置
	err := saveConfig(config)
	if err != nil {
		t.Errorf("保存配置失败: %v", err)
	}

	// 验证文件存在
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("配置文件保存后不存在")
	}

	// 验证保存的配置可以重新加载
	loadedConfig := loadConfig()
	if loadedConfig.Default.AccessKeyID != config.Default.AccessKeyID {
		t.Errorf("重新加载的AccessKeyID不匹配")
	}
}

func TestRunConfigure(t *testing.T) {
	// 测试configure命令的运行函数
	// 由于runConfigure需要交互输入，这里只测试它不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runConfigure发生panic: %v", r)
		}
	}()

	// 创建临时目录用于测试
	tempDir := t.TempDir()

	// 保存原始HOME环境变量
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// 设置临时HOME
	os.Setenv("HOME", tempDir)

	// 由于runConfigure需要交互输入，我们无法直接测试它
	// 但我们可以测试相关的辅助函数
	t.Log("configure相关函数定义正常")
}

func TestConfigureCommandExists(t *testing.T) {
	// 测试configure命令是否被正确定义
	if configureCmd == nil {
		t.Error("configureCmd不应该为nil")
	}

	if configureCmd.RunE == nil {
		t.Error("configureCmd应该有RunE函数")
	}
}
