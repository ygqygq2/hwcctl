package main

import (
	"fmt"
	"os"

	"github.com/ygqygq2/hwcctl/cmd"
)

var (
	// 构建时注入的版本信息
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 设置版本信息
	cmd.SetVersionInfo(version, buildTime, gitCommit)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}
