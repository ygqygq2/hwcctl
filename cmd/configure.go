package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ygqygq2/hwcctl/internal/auth"
	"github.com/ygqygq2/hwcctl/internal/logx"
	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	Default Profile `yaml:"default"`
}

// Profile 配置文件中的 profile
type Profile struct {
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	Output          string `yaml:"output"`
}

// configureCmd 代表配置命令
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "配置华为云凭证和设置",
	Long: `交互式配置华为云访问凭证、默认区域和输出格式。
类似于 AWS CLI 的 configure 命令，将配置保存到 ~/.hwcctl/config 文件中。`,
	RunE: runConfigure,
}

func runConfigure(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// 获取当前配置
	config := loadConfig()

	fmt.Println("华为云 CLI 配置")
	fmt.Println("请输入你的华为云访问凭证信息:")

	// 配置 Access Key ID
	fmt.Printf("Huawei Cloud Access Key ID [%s]: ", maskString(config.Default.AccessKeyID))
	accessKey, _ := reader.ReadString('\n')
	accessKey = strings.TrimSpace(accessKey)
	if accessKey != "" {
		config.Default.AccessKeyID = accessKey
	}

	// 配置 Secret Access Key
	fmt.Printf("Huawei Cloud Secret Access Key [%s]: ", maskString(config.Default.SecretAccessKey))
	secretKey, _ := reader.ReadString('\n')
	secretKey = strings.TrimSpace(secretKey)
	if secretKey != "" {
		config.Default.SecretAccessKey = secretKey
	}

	// 配置默认区域
	fmt.Printf("Default region name [%s]: ", config.Default.Region)
	region, _ := reader.ReadString('\n')
	region = strings.TrimSpace(region)
	if region != "" {
		config.Default.Region = region
	}

	// 配置输出格式
	fmt.Printf("Default output format [%s]: ", config.Default.Output)
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)
	if output != "" {
		config.Default.Output = output
	}

	// 保存配置
	if err := saveConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	fmt.Println("✅ 配置已保存")
	logx.Infof("配置文件已保存到: %s", auth.ResolveConfigPath())

	return nil
}

// loadConfig 加载配置文件
func loadConfig() Config {
	configPath := auth.ResolveConfigPath()
	config := Config{
		Default: Profile{
			Region: "cn-north-1",
			Output: "table",
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		logx.Warnf("读取配置文件失败: %v", err)
		return config
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		logx.Warnf("解析配置文件失败: %v", err)
		return config
	}

	return config
}

// saveConfig 保存配置文件
func saveConfig(config Config) error {
	configPath := auth.ResolveConfigPath()
	configDir := filepath.Dir(configPath)

	// 创建配置目录
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 序列化配置
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// maskString 掩码字符串（用于显示部分信息）
func maskString(s string) string {
	if len(s) <= 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-4)
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
