package auth

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	accessKey := "test-access-key"
	secretKey := "test-secret-key"
	region := "cn-north-1"

	config := NewConfig(accessKey, secretKey, region)

	if config.AccessKey != accessKey {
		t.Errorf("期望 AccessKey 为 '%s'，实际为 '%s'", accessKey, config.AccessKey)
	}
	if config.SecretKey != secretKey {
		t.Errorf("期望 SecretKey 为 '%s'，实际为 '%s'", secretKey, config.SecretKey)
	}
	if config.Region != region {
		t.Errorf("期望 Region 为 '%s'，实际为 '%s'", region, config.Region)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// 保存原始环境变量
	originalAccessKey := os.Getenv("HUAWEICLOUD_ACCESS_KEY")
	originalSecretKey := os.Getenv("HUAWEICLOUD_SECRET_KEY")
	originalRegion := os.Getenv("HUAWEICLOUD_REGION")

	// 清理环境变量
	defer func() {
		os.Setenv("HUAWEICLOUD_ACCESS_KEY", originalAccessKey)
		os.Setenv("HUAWEICLOUD_SECRET_KEY", originalSecretKey)
		os.Setenv("HUAWEICLOUD_REGION", originalRegion)
	}()

	t.Run("成功加载环境变量", func(t *testing.T) {
		os.Setenv("HUAWEICLOUD_ACCESS_KEY", "test-key")
		os.Setenv("HUAWEICLOUD_SECRET_KEY", "test-secret")
		os.Setenv("HUAWEICLOUD_REGION", "cn-east-2")

		config, err := LoadFromEnv()
		if err != nil {
			t.Errorf("期望成功，实际出错: %v", err)
		}
		if config.AccessKey != "test-key" {
			t.Errorf("期望 AccessKey 为 'test-key'，实际为 '%s'", config.AccessKey)
		}
		if config.Region != "cn-east-2" {
			t.Errorf("期望 Region 为 'cn-east-2'，实际为 '%s'", config.Region)
		}
	})

	t.Run("缺少 AccessKey", func(t *testing.T) {
		os.Unsetenv("HUAWEICLOUD_ACCESS_KEY")
		os.Setenv("HUAWEICLOUD_SECRET_KEY", "test-secret")

		_, err := LoadFromEnv()
		if err == nil {
			t.Error("期望出错，但没有错误")
		}
	})

	t.Run("使用默认区域", func(t *testing.T) {
		os.Setenv("HUAWEICLOUD_ACCESS_KEY", "test-key")
		os.Setenv("HUAWEICLOUD_SECRET_KEY", "test-secret")
		os.Unsetenv("HUAWEICLOUD_REGION")

		config, err := LoadFromEnv()
		if err != nil {
			t.Errorf("期望成功，实际出错: %v", err)
		}
		if config.Region != "cn-north-1" {
			t.Errorf("期望默认区域为 'cn-north-1'，实际为 '%s'", config.Region)
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("有效配置", func(t *testing.T) {
		config := &Config{
			AccessKey: "test-key",
			SecretKey: "test-secret",
			Region:    "cn-north-1",
		}
		if err := config.Validate(); err != nil {
			t.Errorf("期望配置有效，实际出错: %v", err)
		}
	})

	t.Run("缺少 AccessKey", func(t *testing.T) {
		config := &Config{
			SecretKey: "test-secret",
			Region:    "cn-north-1",
		}
		if err := config.Validate(); err == nil {
			t.Error("期望验证失败，但没有错误")
		}
	})

	t.Run("缺少 SecretKey", func(t *testing.T) {
		config := &Config{
			AccessKey: "test-key",
			Region:    "cn-north-1",
		}
		if err := config.Validate(); err == nil {
			t.Error("期望验证失败，但没有错误")
		}
	})

	t.Run("缺少 Region", func(t *testing.T) {
		config := &Config{
			AccessKey: "test-key",
			SecretKey: "test-secret",
		}
		if err := config.Validate(); err == nil {
			t.Error("期望验证失败，但没有错误")
		}
	})
}
