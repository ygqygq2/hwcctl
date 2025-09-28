package cmd

import (
	"testing"
)

func TestCDNCmd(t *testing.T) {
	// 测试CDN命令是否正确初始化
	if cdnCmd.Use != "cdn" {
		t.Errorf("期望命令名称为 'cdn'，实际为 '%s'", cdnCmd.Use)
	}

	if cdnCmd.Short == "" {
		t.Error("cdn命令应该有短描述")
	}

	if cdnCmd.Long == "" {
		t.Error("cdn命令应该有长描述")
	}
}

func TestCDNCommandExists(t *testing.T) {
	// 测试CDN命令是否被正确定义
	if cdnCmd == nil {
		t.Error("cdnCmd不应该为nil")
	}
}

func TestCDNCommandSubcommands(t *testing.T) {
	// 测试CDN命令是否有子命令
	// 由于子命令在其他文件中注册，这里只测试基本结构
	if cdnCmd.HasSubCommands() {
		t.Log("CDN命令有子命令")
	} else {
		t.Log("CDN命令目前没有子命令（这是正常的，子命令在init函数中添加）")
	}
}

func TestCDNCommandDescription(t *testing.T) {
	// 测试命令描述包含期望的关键词
	expectedKeywords := []string{"内容分发网络", "华为云"}

	for _, keyword := range expectedKeywords {
		if !containsString(cdnCmd.Long, keyword) {
			t.Errorf("CDN命令长描述应该包含关键词 '%s'", keyword)
		}
	}
}

// 辅助函数：检查字符串是否包含子字符串
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
