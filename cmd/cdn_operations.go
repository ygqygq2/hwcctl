package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ygqygq2/hwcctl/internal/cdn"
	"github.com/ygqygq2/hwcctl/internal/logx"
)

// refreshCmd 代表 CDN 刷新命令
var cdnRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "刷新 CDN 缓存",
	Long:  `刷新指定的 URL 或目录的 CDN 缓存，支持批量刷新。`,
	RunE:  runCDNRefresh,
}

// preloadCmd 代表 CDN 预热命令
var cdnPreloadCmd = &cobra.Command{
	Use:   "preload",
	Short: "预热 CDN 缓存",
	Long:  `预热指定的 URL 到 CDN 边缘节点，提高访问速度。`,
	RunE:  runCDNPreload,
}

func runCDNRefresh(cmd *cobra.Command, args []string) error {
	// 获取参数
	urls, _ := cmd.Flags().GetStringSlice("urls")
	refreshType, _ := cmd.Flags().GetString("type")
	
	if len(urls) == 0 {
		return fmt.Errorf("请指定要刷新的 URL 或目录")
	}

	// 创建 CDN 客户端
	client, err := cdn.NewClient()
	if err != nil {
		return fmt.Errorf("创建 CDN 客户端失败: %v", err)
	}

	logx.Infof("开始刷新 CDN 缓存，类型: %s", refreshType)
	logx.Infof("待刷新的 URL/目录: %s", strings.Join(urls, ", "))

	// 执行刷新
	taskId, err := client.RefreshCache(urls, refreshType)
	if err != nil {
		return fmt.Errorf("刷新 CDN 缓存失败: %v", err)
	}

	logx.Infof("CDN 缓存刷新任务已提交，任务 ID: %s", taskId)
	fmt.Printf("✅ CDN 缓存刷新任务已提交成功\n")
	fmt.Printf("任务 ID: %s\n", taskId)
	fmt.Printf("可以使用以下命令查询任务状态:\n")
	fmt.Printf("hwcctl cdn task %s\n", taskId)

	return nil
}

func runCDNPreload(cmd *cobra.Command, args []string) error {
	// 获取参数
	urls, _ := cmd.Flags().GetStringSlice("urls")
	
	if len(urls) == 0 {
		return fmt.Errorf("请指定要预热的 URL")
	}

	// 创建 CDN 客户端
	client, err := cdn.NewClient()
	if err != nil {
		return fmt.Errorf("创建 CDN 客户端失败: %v", err)
	}

	logx.Infof("开始预热 CDN 缓存")
	logx.Infof("待预热的 URL: %s", strings.Join(urls, ", "))

	// 执行预热
	taskId, err := client.PreloadCache(urls)
	if err != nil {
		return fmt.Errorf("预热 CDN 缓存失败: %v", err)
	}

	logx.Infof("CDN 缓存预热任务已提交，任务 ID: %s", taskId)
	fmt.Printf("✅ CDN 缓存预热任务已提交成功\n")
	fmt.Printf("任务 ID: %s\n", taskId)
	fmt.Printf("可以使用以下命令查询任务状态:\n")
	fmt.Printf("hwcctl cdn task %s\n", taskId)

	return nil
}

// taskCmd 代表查询 CDN 任务状态命令
var cdnTaskCmd = &cobra.Command{
	Use:   "task [task-id]",
	Short: "查询 CDN 任务状态",
	Long:  `查询指定任务 ID 的 CDN 刷新或预热任务状态。`,
	Args:  cobra.ExactArgs(1),
	RunE:  runCDNTask,
}

func runCDNTask(cmd *cobra.Command, args []string) error {
	taskId := args[0]

	// 创建 CDN 客户端
	client, err := cdn.NewClient()
	if err != nil {
		return fmt.Errorf("创建 CDN 客户端失败: %v", err)
	}

	logx.Infof("查询 CDN 任务状态，任务 ID: %s", taskId)

	// 查询任务状态
	task, err := client.GetTaskStatus(taskId)
	if err != nil {
		return fmt.Errorf("查询任务状态失败: %v", err)
	}

	fmt.Printf("📋 CDN 任务状态\n")
	fmt.Printf("任务 ID: %s\n", task.ID)
	fmt.Printf("任务类型: %s\n", task.Type)
	fmt.Printf("任务状态: %s\n", task.Status)
	fmt.Printf("创建时间: %s\n", task.CreatedAt)
	if task.CompletedAt != "" {
		fmt.Printf("完成时间: %s\n", task.CompletedAt)
	}
	fmt.Printf("处理进度: %d%%\n", task.Progress)

	return nil
}

func init() {
	// 添加 CDN 子命令
	cdnCmd.AddCommand(cdnRefreshCmd)
	cdnCmd.AddCommand(cdnPreloadCmd)
	cdnCmd.AddCommand(cdnTaskCmd)

	// CDN 刷新命令的标志
	cdnRefreshCmd.Flags().StringSlice("urls", []string{}, "要刷新的 URL 或目录列表（必需）")
	cdnRefreshCmd.Flags().String("type", "url", "刷新类型：url（文件）或 directory（目录）")
	cdnRefreshCmd.MarkFlagRequired("urls")

	// CDN 预热命令的标志
	cdnPreloadCmd.Flags().StringSlice("urls", []string{}, "要预热的 URL 列表（必需）")
	cdnPreloadCmd.MarkFlagRequired("urls")
}
