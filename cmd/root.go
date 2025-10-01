package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ygqygq2/hwcctl/internal/auth"
	"github.com/ygqygq2/hwcctl/internal/logx"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// SetVersionInfo 设置版本信息
func SetVersionInfo(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
	// 更新 rootCmd 的版本信息
	rootCmd.Version = v
	// 重新设置版本模板
	rootCmd.SetVersionTemplate(`hwcctl 版本 {{.Version}}
构建时间: ` + bt + `
Git 提交: ` + gc + `
`)
}

// rootCmd 代表没有调用子命令时的基础命令
var rootCmd = &cobra.Command{
	Use:   "hwcctl",
	Short: "华为云命令行工具",
	Long: `hwcctl 是一个华为云命令行工具，用于调用华为云 API 进行各种运维操作。

类似于 AWS CLI，但专门针对华为云服务设计。`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 处理调试标志
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logx.SetLevel("debug")
		}

		configPath, _ := cmd.Flags().GetString("config")
		if configPath == "" {
			configPath = os.Getenv("HWCCTL_CONFIG")
		}
		auth.SetConfigPath(configPath)
	},
}

// Execute 添加所有子命令到根命令并适当设置标志
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 设置版本模板
	rootCmd.SetVersionTemplate(`hwcctl 版本 {{.Version}}
构建时间: ` + buildTime + `
Git 提交: ` + gitCommit + `
`)

	// 全局标志
	rootCmd.PersistentFlags().StringP("region", "r", "", "华为云区域 (也可使用环境变量 HUAWEICLOUD_REGION)")
	rootCmd.PersistentFlags().String("access-key-id", "", "Access Key ID (也可使用环境变量 HUAWEICLOUD_ACCESS_KEY)")
	rootCmd.PersistentFlags().String("secret-access-key", "", "Secret Access Key (也可使用环境变量 HUAWEICLOUD_SECRET_KEY)")
	rootCmd.PersistentFlags().String("domain-id", "", "Domain ID (也可使用环境变量 HUAWEICLOUD_DOMAIN_ID)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "调试模式")
	rootCmd.PersistentFlags().String("output", "table", "输出格式 (table|json|yaml)")
	rootCmd.PersistentFlags().String("profile", "", "使用指定的配置文件 profile")
	rootCmd.PersistentFlags().String("config", "", "指定配置文件路径 (默认 ~/.hwcctl/config)")
	rootCmd.PersistentFlags().String("endpoint-url", "", "覆盖默认的服务端点 URL")

	// 设置使用帮助
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "help [command]",
		Short:  "显示命令帮助信息",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rootCmd.Help()
			} else {
				subCmd, _, err := rootCmd.Find(args)
				if err != nil {
					fmt.Fprintf(os.Stderr, "未找到命令: %s\n", args[0])
					os.Exit(1)
				}
				subCmd.Help()
			}
		},
	})
}
