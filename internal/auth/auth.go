package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 华为云认证配置
type Config struct {
	AccessKey           string
	SecretKey           string
	Region              string
	DomainID            string
	ProjectID           string `yaml:"project_id"`            // 项目ID
	EnterpriseProjectID string `yaml:"enterprise_project_id"` // 企业项目ID，默认为 "0"
	MaxRetries          int    `yaml:"max_retries"`           // 最大重试次数，默认 0（不重试）
	EnableRetry         bool   `yaml:"enable_retry"`          // 是否启用重试，默认 false
}

// Credentials 华为云认证凭证
type Credentials struct {
	AccessKeyID         string
	SecretAccessKey     string
	Region              string
	DomainID            string
	ProjectID           string
	EnterpriseProjectID string
}

// Profile 配置文件中的 profile
type Profile struct {
	AccessKeyID         string `yaml:"access_key_id"`
	SecretAccessKey     string `yaml:"secret_access_key"`
	Region              string `yaml:"region"`
	DomainID            string `yaml:"domain_id"`
	ProjectID           string `yaml:"project_id"`            // 项目ID
	EnterpriseProjectID string `yaml:"enterprise_project_id"` // 企业项目ID，默认为 "0"
	Output              string `yaml:"output"`
	MaxRetries          int    `yaml:"max_retries"`  // 最大重试次数，默认 0
	EnableRetry         bool   `yaml:"enable_retry"` // 是否启用重试，默认 false
}

// ConfigFile 配置文件结构
type ConfigFile struct {
	Default Profile `yaml:"default"`
}

// Project 项目信息
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DomainID    string `json:"domain_id"`
	Enabled     bool   `json:"enabled"`
}

// ProjectsResponse IAM项目列表响应
type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

// TokenResponse IAM Token响应
type TokenResponse struct {
	Token struct {
		ExpiresAt string `json:"expires_at"`
	} `json:"token"`
}

// NewConfig 创建新的认证配置
func NewConfig(accessKey, secretKey, region string) *Config {
	return &Config{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}
}

// LoadConfig 加载配置信息，优先级：命令行参数 > 环境变量 > 配置文件
func LoadConfig(accessKeyFlag, secretKeyFlag, regionFlag, domainIDFlag string) (*Config, error) {
	config := &Config{}

	// 1. 尝试从配置文件读取
	configFile := loadConfigFile()
	if configFile != nil {
		config.AccessKey = configFile.Default.AccessKeyID
		config.SecretKey = configFile.Default.SecretAccessKey
		config.Region = configFile.Default.Region
		config.DomainID = configFile.Default.DomainID
		config.ProjectID = configFile.Default.ProjectID
		config.EnterpriseProjectID = configFile.Default.EnterpriseProjectID
		config.MaxRetries = configFile.Default.MaxRetries
		config.EnableRetry = configFile.Default.EnableRetry
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
	if envDomainID := os.Getenv("HUAWEICLOUD_DOMAIN_ID"); envDomainID != "" {
		config.DomainID = envDomainID
	}
	if envProjectID := os.Getenv("HUAWEICLOUD_PROJECT_ID"); envProjectID != "" {
		config.ProjectID = envProjectID
	}
	if envEnterpriseProjectID := os.Getenv("HUAWEICLOUD_ENTERPRISE_PROJECT_ID"); envEnterpriseProjectID != "" {
		config.EnterpriseProjectID = envEnterpriseProjectID
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
	if domainIDFlag != "" {
		config.DomainID = domainIDFlag
	}

	// 设置默认值
	if config.Region == "" {
		config.Region = "cn-north-1"
	}

	// 设置企业项目ID默认值
	if config.EnterpriseProjectID == "" {
		config.EnterpriseProjectID = "0" // "0" 表示默认企业项目
	}

	// 初始化项目管理器
	projectManager := GetProjectManager()
	projectManager.InitWithConfig(config)

	return config, nil
}

// LoadFromEnv 从环境变量加载认证配置（保持向后兼容）
func LoadFromEnv() (*Config, error) {
	return LoadConfig("", "", "", "")
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
		return errors.New("access Key 不能为空")
	}
	if c.SecretKey == "" {
		return errors.New("secret Key 不能为空")
	}
	if c.Region == "" {
		return errors.New("区域不能为空")
	}
	return nil
}

// GetCredentials 获取华为云凭证信息，供 CDN 客户端使用
func GetCredentials() (*Credentials, error) {
	config, err := LoadConfig("", "", "", "")
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyID:         config.AccessKey,
		SecretAccessKey:     config.SecretKey,
		Region:              config.Region,
		DomainID:            config.DomainID,
		ProjectID:           config.ProjectID,
		EnterpriseProjectID: config.EnterpriseProjectID,
	}, nil
}

// GetUnifiedProjectID 获取统一的项目ID（懒加载模式）
func GetUnifiedProjectID() (string, error) {
	projectManager := GetProjectManager()
	return projectManager.GetProjectID()
}

// GetCredentialsWithFlags 根据命令行参数获取凭证信息
func GetCredentialsWithFlags(accessKey, secretKey, region, domainID string) (*Credentials, error) {
	config, err := LoadConfig(accessKey, secretKey, region, domainID)
	if err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyID:         config.AccessKey,
		SecretAccessKey:     config.SecretKey,
		Region:              config.Region,
		DomainID:            config.DomainID,
		ProjectID:           config.ProjectID,
		EnterpriseProjectID: config.EnterpriseProjectID,
	}, nil
}

// FetchProjects 获取项目列表
func (c *Config) FetchProjects() (*ProjectsResponse, error) {
	if c.AccessKey == "" || c.SecretKey == "" {
		return nil, errors.New("accessKey 和 secretKey 不能为空")
	}

	// 构建请求URL
	url := "https://iam.myhuaweicloud.com/v3/projects"

	// 创建HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 生成华为云签名
	err = c.signRequest(req)
	if err != nil {
		return nil, fmt.Errorf("签名请求失败: %v", err)
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var projectsResponse ProjectsResponse
	err = json.Unmarshal(body, &projectsResponse)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &projectsResponse, nil
} // signRequest 为请求添加华为云签名
func (c *Config) signRequest(req *http.Request) error {
	// 获取当前时间
	now := time.Now().UTC()
	timestamp := now.Format("20060102T150405Z")
	date := now.Format("20060102")

	// 设置必要的头部
	req.Header.Set("X-Amz-Date", timestamp)
	req.Header.Set("Host", req.Host)

	// 构建Canonical Request
	canonicalRequest := c.buildCanonicalRequest(req)

	// 构建String to Sign
	credentialScope := fmt.Sprintf("%s/%s/iam/aws4_request", date, c.Region)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		timestamp, credentialScope, c.sha256Hash(canonicalRequest))

	// 计算签名
	signature := c.calculateSignature(c.SecretKey, date, c.Region, "iam", stringToSign)

	// 构建Authorization头
	authorizationHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.AccessKey, credentialScope, c.getSignedHeaders(req), signature)

	req.Header.Set("Authorization", authorizationHeader)

	return nil
}

// buildCanonicalRequest 构建规范请求
func (c *Config) buildCanonicalRequest(req *http.Request) string {
	method := req.Method
	uri := req.URL.EscapedPath()
	if uri == "" {
		uri = "/"
	}

	query := req.URL.RawQuery

	// 构建规范头部
	var headerNames []string
	headerMap := make(map[string]string)

	for name, values := range req.Header {
		lowerName := strings.ToLower(name)
		headerNames = append(headerNames, lowerName)
		headerMap[lowerName] = strings.Join(values, ",")
	}

	sort.Strings(headerNames)

	var canonicalHeaders strings.Builder
	for _, name := range headerNames {
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", name, headerMap[name]))
	}

	signedHeaders := strings.Join(headerNames, ";")

	// 获取请求体的哈希
	var bodyHash string
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
		bodyHash = c.sha256Hash(string(body))
	} else {
		bodyHash = c.sha256Hash("")
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method, uri, query, canonicalHeaders.String(), signedHeaders, bodyHash)
}

// getSignedHeaders 获取签名头部列表
func (c *Config) getSignedHeaders(req *http.Request) string {
	var headerNames []string
	for name := range req.Header {
		headerNames = append(headerNames, strings.ToLower(name))
	}
	sort.Strings(headerNames)
	return strings.Join(headerNames, ";")
}

// calculateSignature 计算签名
func (c *Config) calculateSignature(key, date, region, service, stringToSign string) string {
	kDate := c.hmacSHA256([]byte("AWS4"+key), date)
	kRegion := c.hmacSHA256(kDate, region)
	kService := c.hmacSHA256(kRegion, service)
	kSigning := c.hmacSHA256(kService, "aws4_request")
	signature := c.hmacSHA256(kSigning, stringToSign)
	return hex.EncodeToString(signature)
}

// hmacSHA256 计算HMAC-SHA256
func (c *Config) hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// sha256Hash 计算SHA256哈希
func (c *Config) sha256Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// GetProjectIDByRegion 根据region获取对应的项目ID
func (c *Config) GetProjectIDByRegion(region string) (string, error) {
	projects, err := c.FetchProjects()
	if err != nil {
		return "", err
	}

	for _, project := range projects.Projects {
		if project.Name == region && project.Enabled {
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("在区域 %s 中未找到启用的项目", region)
}
