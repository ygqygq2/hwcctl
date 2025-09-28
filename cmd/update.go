package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
	"github.com/ygqygq2/hwcctl/internal/updater"
)

// updateCmd 代表自更新命令
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新 hwcctl 到最新版本",
	Long: `从 GitHub Releases 下载并安装最新版本的 hwcctl。

此命令会：
1. 检查 GitHub 上的最新版本
2. 下载适合当前操作系统和架构的二进制文件  
3. 验证下载文件的完整性
4. 安全地替换当前可执行文件`,
	RunE:         runUpdate,
	SilenceUsage: true,
}

var (
	// 更新相关标志
	updateForce   bool
	updateCheck   bool
	updateVersion string
)

func runUpdate(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Root().PersistentFlags().GetBool("verbose")
	debug, _ := cmd.Root().PersistentFlags().GetBool("debug")

	// 创建更新器配置
	config := &updater.Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: version,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Verbose:    verbose,
		Debug:      debug,
	}

	// 创建更新器实例
	u := updater.New(config)

	// 如果只是检查版本
	if updateCheck {
		return u.CheckUpdate()
	}

	// 执行更新
	return u.Update(updateForce, updateVersion)
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// 添加更新相关的标志
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "强制更新，即使已经是最新版本")
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "只检查是否有新版本，不执行更新")
	updateCmd.Flags().StringVar(&updateVersion, "version", "", "更新到指定版本（默认为最新版本）")
}
