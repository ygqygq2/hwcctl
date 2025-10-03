package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ygqygq2/hwcctl/cmd"
)

var (
	// 构建时注入的版本信息
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 尝试从构建信息中填充版本信息（适配 `go install module@version` 的场景）
	// 当未通过 ldflags 注入时，尽量从 debug.BuildInfo 获得版本、提交与构建时间。
	if info, ok := debug.ReadBuildInfo(); ok {
		// 主模块版本，例如 v0.0.3；开发态通常为 (devel)
		if version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			// 与 goreleaser 输出保持一致，去掉前缀 v
			version = strings.TrimPrefix(info.Main.Version, "v")
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if gitCommit == "unknown" && s.Value != "" {
					gitCommit = s.Value
				}
			case "vcs.time":
				if buildTime == "unknown" && s.Value != "" {
					buildTime = s.Value
				}
			}
		}
	}
	// 设置版本信息
	cmd.SetVersionInfo(version, buildTime, gitCommit)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}
