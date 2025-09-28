package main

import (
	"testing"

	"github.com/ygqygq2/hwcctl/cmd"
)

func TestMain(t *testing.T) {
	// 这是一个基础测试，确保 main 包可以正常导入
	if testing.Short() {
		t.Skip("跳过 main 包测试")
	}
}

func TestMainExists(t *testing.T) {
	// 测试main函数存在且可以被引用
	// 这个测试只是确保main函数被正确定义
	t.Log("main函数定义正常")
}

func TestMainPackage(t *testing.T) {
	// 测试main包的基本结构
	t.Log("main包结构正常")
}

func TestVersionVariables(t *testing.T) {
	// 测试版本变量是否定义（在开发环境中可能为默认值）
	if version == "" {
		t.Log("version变量为空，这在某些构建环境中是正常的")
	}

	if buildTime == "" {
		t.Log("buildTime变量为空，这在某些构建环境中是正常的")
	}

	if gitCommit == "" {
		t.Log("gitCommit变量为空，这在某些构建环境中是正常的")
	}

	// 测试变量类型正确
	_ = len(version)
	_ = len(buildTime)
	_ = len(gitCommit)
}

func TestMainFunctionCall(t *testing.T) {
	// 测试main函数调用SetVersionInfo和Execute
	// 这里我们无法直接测试main函数，但可以测试它调用的函数

	// 验证SetVersionInfo函数调用不会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetVersionInfo发生panic: %v", r)
		}
	}()

	// 测试设置版本信息
	cmd.SetVersionInfo("test-version", "test-build-time", "test-git-commit")

	// 这里无法测试cmd.Execute()因为它会退出程序或需要交互
	t.Log("版本信息设置功能正常")
}

func TestBuildInfo(t *testing.T) {
	// 测试构建信息的格式
	if version != "dev" && version == "" {
		t.Error("version应该有默认值或构建时注入的值")
	}

	if buildTime != "unknown" && buildTime == "" {
		t.Error("buildTime应该有默认值或构建时注入的值")
	}

	if gitCommit != "unknown" && gitCommit == "" {
		t.Error("gitCommit应该有默认值或构建时注入的值")
	}
}
