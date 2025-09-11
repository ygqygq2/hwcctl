package auth

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 华为云认证配置
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
}

// Profile 配置文件中的 profile
type Profile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	Output          string `yaml:"output"`
}

// ConfigFile 配置文件结构
type ConfigFile struct {
	Default Profile `yaml:"default"`
}

// NewConfig 创建新的认证配置
func NewConfig(accessKey, secretKey, region string) *Config {
	return &Config{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}
}

// LoadConfig 从多个来源加载认证配置（优先级：命令行参数 > 环境变量 > 配置文件）
func LoadConfig(accessKeyFlag, secretKeyFlag, regionFlag string) (*Config, error) {
	config := &Config{}

	// 1. 尝试从配置文件读取
	configFile := loadConfigFile()
	if configFile != nil {
		config.AccessKey = configFile.Default.AccessKeyID
		config.SecretKey = configFile.Default.SecretAccessKey
		config.Region = configFile.Default.Region
	}

	// 2. 从环境变量覆盖
	if envAccessKey := os.Getenv("HUAWEICLOUD_ACCESS_KEY"); envAccessKey != "" {
		config.AccessKey = envAccessKey
	}
	if envSecretKey := os.Getenv("HUAWEICLOUD_SECRET_KEY"); envSecretKey != "" {
		config.SecretKey = envSecretKey
	}
	if envRegion := os.Getenv("HUAWEICLOUD_REGION"); envRegion != "" {
		config.Region = envRegion
	}

	// 3. 从命令行参数覆盖
	if accessKeyFlag != "" {
		config.AccessKey = accessKeyFlag
	}
	if secretKeyFlag != "" {
		config.SecretKey = secretKeyFlag
	}
	if regionFlag != "" {
		config.Region = regionFlag
	}

	// 设置默认值
	if config.Region == "" {
		config.Region = "cn-north-1"
	}

	return config, nil
}

// LoadFromEnv 从环境变量加载认证配置（保持向后兼容）
func LoadFromEnv() (*Config, error) {
	return LoadConfig("", "", "")
}

// loadConfigFile 加载配置文件
func loadConfigFile() *ConfigFile {
	configPath := getConfigPath()
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config ConfigFile
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil
	}

	return &config
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./.hwcctl/config"
	}
	return filepath.Join(homeDir, ".hwcctl", "config")
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.AccessKey == "" {
		return errors.New("Access Key 不能为空")
	}
	if c.SecretKey == "" {
		return errors.New("Secret Key 不能为空")
	}
	if c.Region == "" {
		return errors.New("区域不能为空")
	}
	return nil
}
