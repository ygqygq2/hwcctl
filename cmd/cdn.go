package cmd

import (
	"github.com/spf13/cobra"
)

// cdnCmd 代表 CDN 相关命令
var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "内容分发网络 (CDN) 相关操作",
	Long:  `管理华为云内容分发网络，包括缓存刷新、预热等操作。`,
}

func init() {
	rootCmd.AddCommand(cdnCmd)
}
