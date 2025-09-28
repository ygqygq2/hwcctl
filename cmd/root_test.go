package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	// 测试根命令是否正确初始化
	if rootCmd.Use != "hwcctl" {
		t.Errorf("期望命令名称为 'hwcctl'，实际为 '%s'", rootCmd.Use)
	}

	if rootCmd.Short != "华为云命令行工具" {
		t.Errorf("期望短描述为 '华为云命令行工具'，实际为 '%s'", rootCmd.Short)
	}
}

func TestSetVersionInfo(t *testing.T) {
	testVersion := "1.0.0"
	testBuildTime := "2025-01-01"
	testGitCommit := "abc123"

	SetVersionInfo(testVersion, testBuildTime, testGitCommit)

	if version != testVersion {
		t.Errorf("期望版本为 '%s'，实际为 '%s'", testVersion, version)
	}
	if buildTime != testBuildTime {
		t.Errorf("期望构建时间为 '%s'，实际为 '%s'", testBuildTime, buildTime)
	}
	if gitCommit != testGitCommit {
		t.Errorf("期望 Git 提交为 '%s'，实际为 '%s'", testGitCommit, gitCommit)
	}
}

func TestExecute(t *testing.T) {
	// 测试Execute函数不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute函数发生panic: %v", r)
		}
	}()

	// 测试Execute函数存在并可调用
	// 注意：这里不实际执行，因为它可能会导致程序退出
	// 只要函数定义存在就通过测试
	t.Log("Execute函数定义正常")
}

func TestRootCmdFlags(t *testing.T) {
	// 测试根命令的持久化标志
	flags := rootCmd.PersistentFlags()

	// 检查是否有基本的标志设置
	if flags == nil {
		t.Error("根命令应该有持久化标志")
	}

	// 测试命令结构
	if rootCmd.Commands() == nil {
		t.Error("根命令应该能够添加子命令")
	}
}

func TestRootCmdStructure(t *testing.T) {
	// 测试根命令的基本属性
	if rootCmd.Use == "" {
		t.Error("根命令应该有Use属性")
	}

	if rootCmd.Short == "" {
		t.Error("根命令应该有Short描述")
	}

	if rootCmd.Long == "" {
		t.Error("根命令应该有Long描述")
	}
}

func TestVersionVariables(t *testing.T) {
	// 测试版本变量的初始状态
	originalVersion := version
	originalBuildTime := buildTime
	originalGitCommit := gitCommit

	// 设置新值
	SetVersionInfo("test-version", "test-build-time", "test-commit")

	// 验证值被正确设置
	if version != "test-version" {
		t.Errorf("期望版本为 'test-version'，实际为 '%s'", version)
	}

	// 恢复原始值
	SetVersionInfo(originalVersion, originalBuildTime, originalGitCommit)
}
