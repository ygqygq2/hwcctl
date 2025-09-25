package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ygqygq2/hwcctl/internal/auth"
	"github.com/ygqygq2/hwcctl/internal/cdn"
	hwErrors "github.com/ygqygq2/hwcctl/internal/errors"
	"github.com/ygqygq2/hwcctl/internal/logx"
	"github.com/ygqygq2/hwcctl/internal/output"
	"github.com/ygqygq2/hwcctl/internal/retry"
)

// refreshCmd 代表 CDN 刷新命令
var cdnRefreshCmd = &cobra.Command{
	Use:          "refresh",
	Short:        "刷新 CDN 缓存",
	Long:         `刷新指定的 URL 或目录的 CDN 缓存，支持批量刷新。`,
	RunE:         runCDNRefresh,
	SilenceUsage: true, // 发生错误时不显示用法信息
}

// preloadCmd 代表 CDN 预热命令
var cdnPreloadCmd = &cobra.Command{
	Use:          "preload",
	Short:        "预热 CDN 缓存",
	Long:         `预热指定的 URL 到 CDN 边缘节点，提高访问速度。`,
	RunE:         runCDNPreload,
	SilenceUsage: true, // 发生错误时不显示用法信息
}

func runCDNRefresh(cmd *cobra.Command, args []string) error {
	// 获取参数
	urls, _ := cmd.Flags().GetStringSlice("urls")
	refreshType, _ := cmd.Flags().GetString("type")

	if len(urls) == 0 {
		return hwErrors.NewValidationError("请指定要刷新的 URL 或目录")
	}

	// 获取输出格式
	outputFormat, _ := cmd.Root().PersistentFlags().GetString("output")
	formatter := output.NewFormatter(outputFormat)

	// 获取配置信息以确定重试策略
	config, err := auth.LoadConfig("", "", "", "")
	if err != nil {
		return hwErrors.NewServerError(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 创建重试器 - 根据配置决定是否重试
	var retryer *retry.Retryer
	if config.EnableRetry && config.MaxRetries > 0 {
		retryConfig := retry.DefaultConfig()
		retryConfig.MaxAttempts = config.MaxRetries + 1 // MaxRetries 是重试次数，需要加1为总尝试次数
		retryer = retry.NewRetryer(retryConfig)
	} else {
		// 不重试，只执行一次
		retryConfig := &retry.Config{
			MaxAttempts: 1,
			Strategy:    retry.StrategyFixed,
			BaseDelay:   0,
		}
		retryer = retry.NewRetryer(retryConfig)
	}

	logx.Infof("开始刷新 CDN 缓存，类型: %s", refreshType)
	logx.Infof("待刷新的 URL/目录: %s", strings.Join(urls, ", "))

	// 使用重试机制执行刷新
	ctx := context.Background()
	result, err := retryer.DoWithResult(ctx, func() (interface{}, error) {
		// 创建 CDN 客户端
		client, err := cdn.NewClient()
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("创建 CDN 客户端失败: %v", err))
		}

		// 执行刷新
		taskId, err := client.RefreshCache(urls, refreshType)
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("刷新 CDN 缓存失败: %v", err))
		}

		return map[string]interface{}{
			"task_id":    taskId,
			"type":       "refresh",
			"urls":       urls,
			"status":     "submitted",
			"created_at": time.Now().Format(time.RFC3339),
		}, nil
	})

	if err != nil {
		formatter.PrintError(fmt.Sprintf("CDN 缓存刷新失败: %v", err))
		return err
	}

	// 输出结果
	if outputFormat == "table" || outputFormat == "text" {
		taskData := result.(map[string]interface{})
		taskId := taskData["task_id"].(string)

		logx.Infof("CDN 缓存刷新任务已提交，任务 ID: %s", taskId)
		formatter.PrintSuccess("CDN 缓存刷新任务已提交成功")
		fmt.Printf("任务 ID: %s\n", taskId)
		fmt.Printf("可以使用以下命令查询任务状态:\n")
		fmt.Printf("hwcctl cdn task %s\n", taskId)
	} else {
		formatter.Print(result)
	}

	return nil
}

func runCDNPreload(cmd *cobra.Command, args []string) error {
	// 获取参数
	urls, _ := cmd.Flags().GetStringSlice("urls")

	if len(urls) == 0 {
		return hwErrors.NewValidationError("请指定要预热的 URL")
	}

	// 获取输出格式
	outputFormat, _ := cmd.Root().PersistentFlags().GetString("output")
	formatter := output.NewFormatter(outputFormat)

	// 获取配置信息以确定重试策略
	config, err := auth.LoadConfig("", "", "", "")
	if err != nil {
		return hwErrors.NewServerError(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 创建重试器 - 根据配置决定是否重试
	var retryer *retry.Retryer
	if config.EnableRetry && config.MaxRetries > 0 {
		retryConfig := retry.DefaultConfig()
		retryConfig.MaxAttempts = config.MaxRetries + 1 // MaxRetries 是重试次数，需要加1为总尝试次数
		retryer = retry.NewRetryer(retryConfig)
	} else {
		// 不重试，只执行一次
		retryConfig := &retry.Config{
			MaxAttempts: 1,
			Strategy:    retry.StrategyFixed,
			BaseDelay:   0,
		}
		retryer = retry.NewRetryer(retryConfig)
	}

	logx.Infof("开始预热 CDN 缓存")
	logx.Infof("待预热的 URL: %s", strings.Join(urls, ", "))

	// 使用重试机制执行预热
	ctx := context.Background()
	result, err := retryer.DoWithResult(ctx, func() (interface{}, error) {
		// 创建 CDN 客户端
		client, err := cdn.NewClient()
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("创建 CDN 客户端失败: %v", err))
		}

		// 执行预热
		taskId, err := client.PreloadCache(urls)
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("预热 CDN 缓存失败: %v", err))
		}

		return map[string]interface{}{
			"task_id":    taskId,
			"type":       "preload",
			"urls":       urls,
			"status":     "submitted",
			"created_at": time.Now().Format(time.RFC3339),
		}, nil
	})

	if err != nil {
		formatter.PrintError(fmt.Sprintf("CDN 缓存预热失败: %v", err))
		return err
	}

	// 输出结果
	if outputFormat == "table" || outputFormat == "text" {
		taskData := result.(map[string]interface{})
		taskId := taskData["task_id"].(string)

		logx.Infof("CDN 缓存预热任务已提交，任务 ID: %s", taskId)
		formatter.PrintSuccess("CDN 缓存预热任务已提交成功")
		fmt.Printf("任务 ID: %s\n", taskId)
		fmt.Printf("可以使用以下命令查询任务状态:\n")
		fmt.Printf("hwcctl cdn task %s\n", taskId)
	} else {
		formatter.Print(result)
	}

	return nil
}

// taskCmd 代表查询 CDN 任务状态命令
var cdnTaskCmd = &cobra.Command{
	Use:          "task [task-id]",
	Short:        "查询 CDN 任务状态",
	Long:         `查询指定任务 ID 的 CDN 刷新或预热任务状态。`,
	Args:         cobra.ExactArgs(1),
	RunE:         runCDNTask,
	SilenceUsage: true, // 发生错误时不显示用法信息
}

func runCDNTask(cmd *cobra.Command, args []string) error {
	taskId := args[0]

	// 获取输出格式
	outputFormat, _ := cmd.Root().PersistentFlags().GetString("output")
	formatter := output.NewFormatter(outputFormat)

	// 获取配置信息以确定重试策略
	config, err := auth.LoadConfig("", "", "", "")
	if err != nil {
		return hwErrors.NewServerError(fmt.Sprintf("加载配置失败: %v", err))
	}

	// 创建重试器 - 根据配置决定是否重试
	var retryer *retry.Retryer
	if config.EnableRetry && config.MaxRetries > 0 {
		retryConfig := retry.DefaultConfig()
		retryConfig.MaxAttempts = config.MaxRetries + 1 // MaxRetries 是重试次数，需要加1为总尝试次数
		retryer = retry.NewRetryer(retryConfig)
	} else {
		// 不重试，只执行一次
		retryConfig := &retry.Config{
			MaxAttempts: 1,
			Strategy:    retry.StrategyFixed,
			BaseDelay:   0,
		}
		retryer = retry.NewRetryer(retryConfig)
	}

	logx.Infof("查询 CDN 任务状态，任务 ID: %s", taskId)

	// 使用重试机制查询任务状态
	ctx := context.Background()
	result, err := retryer.DoWithResult(ctx, func() (interface{}, error) {
		// 创建 CDN 客户端
		client, err := cdn.NewClient()
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("创建 CDN 客户端失败: %v", err))
		}

		// 查询任务状态
		task, err := client.GetTaskStatus(taskId)
		if err != nil {
			return nil, hwErrors.NewServerError(fmt.Sprintf("查询任务状态失败: %v", err))
		}

		return task, nil
	})

	if err != nil {
		formatter.PrintError(fmt.Sprintf("查询 CDN 任务状态失败: %v", err))
		return err
	}

	// 输出结果
	if outputFormat == "table" || outputFormat == "text" {
		task := result.(*cdn.Task)
		fmt.Printf("📋 CDN 任务状态\n")
		fmt.Printf("任务 ID: %s\n", task.ID)
		fmt.Printf("任务类型: %s\n", task.Type)
		fmt.Printf("任务状态: %s\n", task.Status)
		fmt.Printf("创建时间: %s\n", task.CreatedAt)
		if task.CompletedAt != "" {
			fmt.Printf("完成时间: %s\n", task.CompletedAt)
		}
		fmt.Printf("处理进度: %d%%\n", task.Progress)
	} else {
		formatter.Print(result)
	}

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
