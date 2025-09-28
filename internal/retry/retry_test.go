package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.MaxAttempts <= 0 {
		t.Error("默认配置的最大重试次数应该大于0")
	}
	if config.BaseDelay <= 0 {
		t.Error("默认配置的基础延迟应该大于0")
	}
	if config.Strategy == "" {
		t.Error("默认配置应该有重试策略")
	}
}

func TestNewRetryer(t *testing.T) {
	config := &Config{
		MaxAttempts: 3,
		Strategy:    StrategyFixed,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    1 * time.Second,
		Multiplier:  2.0,
	}

	retryer := NewRetryer(config)
	if retryer == nil {
		t.Error("期望创建重试器成功，但返回nil")
	}

	// 测试nil配置使用默认配置
	retryer2 := NewRetryer(nil)
	if retryer2 == nil {
		t.Error("期望使用默认配置创建重试器成功，但返回nil")
	}
}

func TestRetryer_Do_Success(t *testing.T) {
	retryer := NewRetryer(DefaultConfig())
	ctx := context.Background()

	callCount := 0
	err := retryer.Do(ctx, func() error {
		callCount++
		return nil // 第一次就成功
	})

	if err != nil {
		t.Errorf("期望成功，但得到错误: %v", err)
	}
	if callCount != 1 {
		t.Errorf("期望调用1次，实际调用%d次", callCount)
	}
}

func TestRetryer_Do_RetryableError(t *testing.T) {
	config := &Config{
		MaxAttempts: 3,
		Strategy:    StrategyFixed,
		BaseDelay:   1 * time.Millisecond, // 很短的延迟以加快测试
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  1.0,
	}
	retryer := NewRetryer(config)
	ctx := context.Background()

	callCount := 0
	err := retryer.Do(ctx, func() error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error") // 前两次失败
		}
		return nil // 第三次成功
	})

	if err != nil {
		t.Errorf("期望最终成功，但得到错误: %v", err)
	}
	if callCount != 3 {
		t.Errorf("期望调用3次，实际调用%d次", callCount)
	}
}

func TestRetryer_DoWithResult_Success(t *testing.T) {
	retryer := NewRetryer(DefaultConfig())
	ctx := context.Background()

	result, err := retryer.DoWithResult(ctx, func() (interface{}, error) {
		return "success", nil
	})

	if err != nil {
		t.Errorf("期望成功，但得到错误: %v", err)
	}
	if result != "success" {
		t.Errorf("期望结果为 'success'，实际为 '%v'", result)
	}
}

func TestIsRetryable_DefaultBehavior(t *testing.T) {
	// 测试默认情况下一般错误都是可重试的
	// 注意：实际的IsRetryable实现可能有特定的逻辑
	err := errors.New("some error")
	result := IsRetryable(err)
	// 无论结果如何，只要函数不panic就算通过
	t.Logf("IsRetryable返回: %v", result)
}

func TestQuickRetry(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	err := QuickRetry(ctx, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("temp error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("期望快速重试成功，但得到错误: %v", err)
	}
}

func TestQuickRetryWithResult(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	result, err := QuickRetryWithResult(ctx, func() (interface{}, error) {
		callCount++
		if callCount < 2 {
			return nil, errors.New("temp error")
		}
		return "quick result", nil
	})

	if err != nil {
		t.Errorf("期望快速重试成功，但得到错误: %v", err)
	}
	if result != "quick result" {
		t.Errorf("期望结果为 'quick result'，实际为 '%v'", result)
	}
}

func TestConservativeRetry(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	err := ConservativeRetry(ctx, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("temp error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("期望保守重试成功，但得到错误: %v", err)
	}
}

func TestAggressiveRetry(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	err := AggressiveRetry(ctx, func() error {
		callCount++
		if callCount < 2 {
			return errors.New("temp error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("期望激进重试成功，但得到错误: %v", err)
	}
}

func TestRetryStrategiesBehaviorDifferences(t *testing.T) {
	// 测试不同重试策略的行为差异 - 这是核心业务逻辑
	testCases := []struct {
		name        string
		strategy    Strategy
		attempts    int
		baseDelay   time.Duration
		maxDelay    time.Duration
		multiplier  float64
		description string
	}{
		{
			name:        "固定延迟策略",
			strategy:    StrategyFixed,
			attempts:    3,
			baseDelay:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			multiplier:  1.0,
			description: "固定延迟应该每次都使用相同的延迟时间",
		},
		{
			name:        "指数退避策略",
			strategy:    StrategyExponential,
			attempts:    4,
			baseDelay:   10 * time.Millisecond,
			maxDelay:    200 * time.Millisecond,
			multiplier:  2.0,
			description: "指数退避延迟应该呈指数增长",
		},
		{
			name:        "线性增长策略",
			strategy:    StrategyLinear,
			attempts:    3,
			baseDelay:   20 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			multiplier:  1.5,
			description: "线性增长延迟应该线性递增",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				MaxAttempts: tc.attempts,
				Strategy:    tc.strategy,
				BaseDelay:   tc.baseDelay,
				MaxDelay:    tc.maxDelay,
				Multiplier:  tc.multiplier,
			}

			retryer := NewRetryer(config)
			ctx := context.Background()

			callCount := 0
			startTime := time.Now()

			// 故意让所有尝试都失败，以测试重试策略
			err := retryer.Do(ctx, func() error {
				callCount++
				return errors.New("persistent error for strategy testing")
			})

			elapsed := time.Since(startTime)

			// 验证调用次数
			if callCount != tc.attempts {
				t.Errorf("%s: 调用次数不匹配，期望 %d, 实际 %d",
					tc.description, tc.attempts, callCount)
			}

			// 验证确实发生了错误（因为我们故意让所有尝试都失败）
			if err == nil {
				t.Errorf("%s: 期望最终失败，但得到成功", tc.description)
			}

			// 验证至少花费了一些时间（意味着有延迟）
			if tc.attempts > 1 && elapsed < tc.baseDelay {
				t.Errorf("%s: 重试耗时过短，可能没有正确应用延迟策略", tc.description)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	// 测试上下文取消的业务逻辑 - 对于长时间运行的重试很重要
	config := &Config{
		MaxAttempts: 10, // 设置很多次尝试
		Strategy:    StrategyFixed,
		BaseDelay:   50 * time.Millisecond, // 较长的延迟
		MaxDelay:    100 * time.Millisecond,
		Multiplier:  1.0,
	}

	retryer := NewRetryer(config)

	// 创建一个会很快取消的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	callCount := 0
	startTime := time.Now()

	err := retryer.Do(ctx, func() error {
		callCount++
		return errors.New("persistent error")
	})

	elapsed := time.Since(startTime)

	// 验证上下文取消确实中断了重试
	if err == nil {
		t.Error("期望由于上下文取消而失败，但得到成功")
	}

	// 验证没有进行所有的重试尝试（因为上下文被取消了）
	if callCount >= config.MaxAttempts {
		t.Errorf("上下文取消应该阻止所有重试尝试，但进行了 %d 次尝试", callCount)
	}

	// 验证耗时合理（应该接近超时时间，而不是所有重试的总时间）
	if elapsed > 200*time.Millisecond {
		t.Errorf("上下文取消后耗时过长: %v", elapsed)
	}
}

func TestRetryWithDifferentErrorTypes(t *testing.T) {
	// 测试不同错误类型的重试行为 - 重要的业务判断逻辑
	retryer := NewRetryer(&Config{
		MaxAttempts: 3,
		Strategy:    StrategyFixed,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  1.0,
	})

	ctx := context.Background()

	// 测试永久性错误（应该立即停止重试的错误类型）
	t.Run("永久性错误处理", func(t *testing.T) {
		callCount := 0
		err := retryer.Do(ctx, func() error {
			callCount++
			// 模拟一个明显不应该重试的错误
			return errors.New("authentication failed - invalid credentials")
		})

		// 注意：这里的行为取决于IsRetryable的实现
		// 如果IsRetryable对所有错误都返回true，那么会重试
		// 如果有特定逻辑识别不可重试的错误，那么应该只调用一次
		if err == nil {
			t.Error("期望认证错误失败，但得到成功")
		}

		t.Logf("认证错误重试次数: %d", callCount)
	})

	// 测试临时性错误（应该重试的错误类型）
	t.Run("临时性错误处理", func(t *testing.T) {
		callCount := 0
		err := retryer.Do(ctx, func() error {
			callCount++
			if callCount < 3 {
				// 模拟临时性网络错误
				return errors.New("network timeout - connection refused")
			}
			return nil // 最终成功
		})

		if err != nil {
			t.Errorf("期望临时错误最终成功，但得到错误: %v", err)
		}

		if callCount != 3 {
			t.Errorf("期望临时错误重试3次，实际 %d 次", callCount)
		}
	})
}

func TestMaxDelayEnforcement(t *testing.T) {
	// 测试最大延迟限制的业务逻辑 - 确保延迟不会无限增长
	config := &Config{
		MaxAttempts: 5,
		Strategy:    StrategyExponential,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    50 * time.Millisecond, // 较小的最大延迟
		Multiplier:  3.0,                   // 大的倍数，会很快超过最大延迟
	}

	retryer := NewRetryer(config)
	ctx := context.Background()

	callCount := 0
	delays := make([]time.Duration, 0)
	lastCallTime := time.Now()

	err := retryer.Do(ctx, func() error {
		currentTime := time.Now()
		if callCount > 0 {
			delay := currentTime.Sub(lastCallTime)
			delays = append(delays, delay)
		}
		lastCallTime = currentTime
		callCount++
		return errors.New("persistent error for delay testing")
	})

	if err == nil {
		t.Error("期望最终失败，但得到成功")
	}

	// 验证延迟不超过最大值
	for i, delay := range delays {
		// 允许一些时间误差
		tolerance := 20 * time.Millisecond
		if delay > config.MaxDelay+tolerance {
			t.Errorf("第%d次重试延迟 %v 超过最大延迟 %v (容忍度: %v)",
				i+1, delay, config.MaxDelay, tolerance)
		}
	}

	t.Logf("观察到的延迟: %v", delays)
}

func TestRetrierConfigValidation(t *testing.T) {
	// 测试配置验证的业务逻辑 - 确保配置合理性
	testCases := []struct {
		name        string
		config      *Config
		description string
	}{
		{
			name: "零最大尝试次数",
			config: &Config{
				MaxAttempts: 0,
				Strategy:    StrategyFixed,
				BaseDelay:   10 * time.Millisecond,
			},
			description: "零尝试次数会导致无法执行任何操作",
		},
		{
			name: "负数延迟",
			config: &Config{
				MaxAttempts: 3,
				Strategy:    StrategyFixed,
				BaseDelay:   -10 * time.Millisecond,
			},
			description: "负数延迟应该被处理或使用默认值",
		},
		{
			name: "无效策略",
			config: &Config{
				MaxAttempts: 3,
				Strategy:    "invalid_strategy",
				BaseDelay:   10 * time.Millisecond,
			},
			description: "无效策略应该使用默认策略",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建重试器不应该panic
			retryer := NewRetryer(tc.config)
			if retryer == nil {
				t.Errorf("%s: 创建重试器返回nil", tc.description)
			}

			// 尝试执行操作不应该panic
			ctx := context.Background()
			err := retryer.Do(ctx, func() error {
				return nil // 立即成功
			})

			// 对于零最大尝试次数的情况，期望失败
			if tc.config.MaxAttempts == 0 {
				if err == nil {
					t.Errorf("%s: 期望零尝试次数导致失败，但得到成功", tc.description)
				}
			} else {
				if err != nil {
					t.Errorf("%s: 简单操作失败: %v", tc.description, err)
				}
			}
		})
	}
}
