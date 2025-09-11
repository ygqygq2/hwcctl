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
