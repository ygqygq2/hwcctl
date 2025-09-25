package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/ygqygq2/hwcctl/internal/errors"
)

// Strategy 重试策略
type Strategy string

const (
	// StrategyFixed 固定间隔重试
	StrategyFixed Strategy = "fixed"
	// StrategyExponential 指数退避重试
	StrategyExponential Strategy = "exponential"
	// StrategyLinear 线性递增重试
	StrategyLinear Strategy = "linear"
)

// Config 重试配置
type Config struct {
	// 最大重试次数
	MaxAttempts int
	// 重试策略
	Strategy Strategy
	// 基础延迟时间
	BaseDelay time.Duration
	// 最大延迟时间
	MaxDelay time.Duration
	// 延迟倍数（用于指数退避）
	Multiplier float64
	// 随机化因子（0-1，用于添加抖动）
	Jitter float64
	// 是否启用调试日志
	Debug bool
}

// DefaultConfig 默认重试配置
func DefaultConfig() *Config {
	return &Config{
		MaxAttempts: 3,
		Strategy:    StrategyExponential,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
		Jitter:      0.1,
		Debug:       false,
	}
}

// Retryer 重试器
type Retryer struct {
	config *Config
}

// NewRetryer 创建新的重试器
func NewRetryer(config *Config) *Retryer {
	if config == nil {
		config = DefaultConfig()
	}

	return &Retryer{
		config: config,
	}
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func() error

// Do 执行重试逻辑
func (r *Retryer) Do(ctx context.Context, fn RetryableFunc) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		if r.config.Debug {
			fmt.Printf("[DEBUG] 重试尝试 %d/%d\n", attempt, r.config.MaxAttempts)
		}

		// 执行函数
		err := fn()
		if err == nil {
			if attempt > 1 && r.config.Debug {
				fmt.Printf("[INFO] 重试成功，尝试次数: %d\n", attempt)
			}
			return nil
		}

		lastErr = err

		// 检查是否是华为云错误且不可重试
		if hwErr, ok := err.(*errors.HuaweiCloudError); ok {
			if !hwErr.IsRetryable() {
				if r.config.Debug {
					fmt.Printf("[DEBUG] 错误不可重试: %v\n", err)
				}
				return err
			}
		}

		// 如果是最后一次尝试，直接返回错误
		if attempt == r.config.MaxAttempts {
			break
		}

		// 计算延迟时间
		delay := r.calculateDelay(attempt)

		if r.config.Debug {
			fmt.Printf("[DEBUG] 重试失败: %v, 等待 %v 后重试\n", err, delay)
		}

		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 继续下一次重试
		}
	}

	// 如果最大重试次数为 1，直接返回原始错误，不添加重试信息
	if r.config.MaxAttempts == 1 {
		return lastErr
	}

	return fmt.Errorf("重试失败，已达到最大重试次数 %d: %w", r.config.MaxAttempts, lastErr)
}

// DoWithResult 执行重试逻辑并返回结果
func (r *Retryer) DoWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	var result interface{}

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		if r.config.Debug {
			fmt.Printf("[DEBUG] 重试尝试 %d/%d\n", attempt, r.config.MaxAttempts)
		}

		// 执行函数
		res, err := fn()
		if err == nil {
			if attempt > 1 && r.config.Debug {
				fmt.Printf("[INFO] 重试成功，尝试次数: %d\n", attempt)
			}
			return res, nil
		}

		lastErr = err
		result = res

		// 检查是否是华为云错误且不可重试
		if hwErr, ok := err.(*errors.HuaweiCloudError); ok {
			if !hwErr.IsRetryable() {
				if r.config.Debug {
					fmt.Printf("[DEBUG] 错误不可重试: %v\n", err)
				}
				return result, err
			}
		}

		// 如果是最后一次尝试，直接返回错误
		if attempt == r.config.MaxAttempts {
			break
		}

		// 计算延迟时间
		delay := r.calculateDelay(attempt)

		if r.config.Debug {
			fmt.Printf("[DEBUG] 重试失败: %v, 等待 %v 后重试\n", err, delay)
		}

		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(delay):
			// 继续下一次重试
		}
	}

	// 如果最大重试次数为 1，直接返回原始错误，不添加重试信息
	if r.config.MaxAttempts == 1 {
		return result, lastErr
	}

	return result, fmt.Errorf("重试失败，已达到最大重试次数 %d: %w", r.config.MaxAttempts, lastErr)
}

// calculateDelay 计算延迟时间
func (r *Retryer) calculateDelay(attempt int) time.Duration {
	var delay time.Duration

	switch r.config.Strategy {
	case StrategyFixed:
		delay = r.config.BaseDelay

	case StrategyLinear:
		delay = time.Duration(attempt) * r.config.BaseDelay

	case StrategyExponential:
		delay = time.Duration(float64(r.config.BaseDelay) * math.Pow(r.config.Multiplier, float64(attempt-1)))

	default:
		delay = r.config.BaseDelay
	}

	// 应用最大延迟限制
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	// 添加抖动
	if r.config.Jitter > 0 {
		jitterAmount := float64(delay) * r.config.Jitter
		jitter := time.Duration(rand.Float64() * jitterAmount)
		delay += jitter
	}

	return delay
}

// IsRetryable 检查错误是否可重试
func IsRetryable(err error) bool {
	if hwErr, ok := err.(*errors.HuaweiCloudError); ok {
		return hwErr.IsRetryable()
	}

	// 检查其他类型的可重试错误
	errMsg := err.Error()
	retryablePatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"network",
		"temporary",
		"服务不可用",
		"超时",
	}

	for _, pattern := range retryablePatterns {
		if contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// contains 检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Quick retry functions for common scenarios

// QuickRetry 快速重试（默认配置）
func QuickRetry(ctx context.Context, fn RetryableFunc) error {
	retryer := NewRetryer(DefaultConfig())
	return retryer.Do(ctx, fn)
}

// QuickRetryWithResult 快速重试并返回结果
func QuickRetryWithResult(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	retryer := NewRetryer(DefaultConfig())
	return retryer.DoWithResult(ctx, fn)
}

// AggressiveRetry 激进重试（更多次数，更短间隔）
func AggressiveRetry(ctx context.Context, fn RetryableFunc) error {
	config := &Config{
		MaxAttempts: 5,
		Strategy:    StrategyExponential,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  1.5,
		Jitter:      0.1,
	}
	retryer := NewRetryer(config)
	return retryer.Do(ctx, fn)
}

// ConservativeRetry 保守重试（更少次数，更长间隔）
func ConservativeRetry(ctx context.Context, fn RetryableFunc) error {
	config := &Config{
		MaxAttempts: 2,
		Strategy:    StrategyFixed,
		BaseDelay:   3 * time.Second,
		MaxDelay:    30 * time.Second,
		Jitter:      0.2,
	}
	retryer := NewRetryer(config)
	return retryer.Do(ctx, fn)
}
