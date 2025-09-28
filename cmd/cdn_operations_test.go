package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCDNRefreshCmd(t *testing.T) {
	// 测试CDN刷新命令是否正确初始化
	if cdnRefreshCmd.Use != "refresh" {
		t.Errorf("期望命令名称为 'refresh'，实际为 '%s'", cdnRefreshCmd.Use)
	}

	if cdnRefreshCmd.Short == "" {
		t.Error("refresh命令应该有短描述")
	}

	if cdnRefreshCmd.RunE == nil {
		t.Error("refresh命令应该有RunE函数")
	}
}

func TestCDNPreloadCmd(t *testing.T) {
	// 测试CDN预热命令是否正确初始化
	if cdnPreloadCmd.Use != "preload" {
		t.Errorf("期望命令名称为 'preload'，实际为 '%s'", cdnPreloadCmd.Use)
	}

	if cdnPreloadCmd.Short == "" {
		t.Error("preload命令应该有短描述")
	}

	if cdnPreloadCmd.RunE == nil {
		t.Error("preload命令应该有RunE函数")
	}
}

func TestCDNTaskCmd(t *testing.T) {
	// 测试CDN任务状态查询命令是否正确初始化
	if cdnTaskCmd.Use != "task [task-id]" {
		t.Errorf("期望命令名称为 'task [task-id]'，实际为 '%s'", cdnTaskCmd.Use)
	}

	if cdnTaskCmd.Short == "" {
		t.Error("task命令应该有短描述")
	}

	if cdnTaskCmd.RunE == nil {
		t.Error("task命令应该有RunE函数")
	}

	if cdnTaskCmd.Args == nil {
		t.Error("task命令应该有Args验证函数")
	}
}

func TestRunCDNRefreshValidation(t *testing.T) {
	// 创建模拟的命令
	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().String("output", "table", "Output format")

	cmd := &cobra.Command{}
	cmd.Flags().StringSlice("urls", []string{}, "URLs to refresh")
	cmd.Flags().String("type", "url", "Refresh type")

	// 添加子命令到根命令
	rootCmd.AddCommand(cmd)

	// 测试没有提供URLs的情况
	err := runCDNRefresh(cmd, []string{})
	if err == nil {
		t.Error("期望在没有URLs时返回错误")
	}

	if err != nil && !strings.Contains(err.Error(), "请指定要刷新的 URL 或目录") {
		t.Errorf("期望包含特定验证错误消息，实际为: %v", err)
	}
}

func TestRunCDNPreloadValidation(t *testing.T) {
	// 创建模拟的命令
	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().String("output", "table", "Output format")

	cmd := &cobra.Command{}
	cmd.Flags().StringSlice("urls", []string{}, "URLs to preload")

	// 添加子命令到根命令
	rootCmd.AddCommand(cmd)

	// 测试没有提供URLs的情况
	err := runCDNPreload(cmd, []string{})
	if err == nil {
		t.Error("期望在没有URLs时返回错误")
	}

	if err != nil && !strings.Contains(err.Error(), "请指定要预热的 URL") {
		t.Errorf("期望包含特定验证错误消息，实际为: %v", err)
	}
}

func TestRunCDNTaskValidation(t *testing.T) {
	// 创建模拟的命令
	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().String("output", "table", "Output format")

	cmd := &cobra.Command{}

	// 添加子命令到根命令
	rootCmd.AddCommand(cmd)

	// 测试提供任务ID的情况
	args := []string{"test-task-id"}

	// 由于需要真实的配置和客户端，这里只测试函数不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runCDNTask发生panic: %v", r)
		}
	}()

	// 这会失败，但不应该panic
	_ = runCDNTask(cmd, args)
}

func TestCDNCommandsInitialization(t *testing.T) {
	// 测试命令是否被正确初始化
	if cdnRefreshCmd == nil {
		t.Error("cdnRefreshCmd不应该为nil")
	}

	if cdnPreloadCmd == nil {
		t.Error("cdnPreloadCmd不应该为nil")
	}

	if cdnTaskCmd == nil {
		t.Error("cdnTaskCmd不应该为nil")
	}
}

func TestCDNRefreshCommandFlags(t *testing.T) {
	// 测试refresh命令的标志
	urlsFlag := cdnRefreshCmd.Flags().Lookup("urls")
	if urlsFlag == nil {
		t.Error("refresh命令应该有urls标志")
	}

	typeFlag := cdnRefreshCmd.Flags().Lookup("type")
	if typeFlag == nil {
		t.Error("refresh命令应该有type标志")
	}

	if typeFlag.DefValue != "url" {
		t.Errorf("type标志的默认值应该为 'url'，实际为 '%s'", typeFlag.DefValue)
	}
}

func TestCDNPreloadCommandFlags(t *testing.T) {
	// 测试preload命令的标志
	urlsFlag := cdnPreloadCmd.Flags().Lookup("urls")
	if urlsFlag == nil {
		t.Error("preload命令应该有urls标志")
	}
}
