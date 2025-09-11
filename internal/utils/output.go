package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// OutputFormat 输出格式类型
type OutputFormat string

const (
	TableFormat OutputFormat = "table"
	JSONFormat  OutputFormat = "json"
	YAMLFormat  OutputFormat = "yaml"
)

// FormatOutput 根据指定格式输出数据
func FormatOutput(data interface{}, format string) error {
	switch OutputFormat(strings.ToLower(format)) {
	case JSONFormat:
		return outputJSON(data)
	case YAMLFormat:
		return outputYAML(data)
	case TableFormat:
		return outputTable(data)
	default:
		return fmt.Errorf("不支持的输出格式: %s", format)
	}
}

// outputJSON 以 JSON 格式输出
func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// outputYAML 以 YAML 格式输出
func outputYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(data)
}

// outputTable 以表格格式输出
func outputTable(data interface{}) error {
	// TODO: 实现表格输出格式
	// 这里先简单输出，后续可以使用第三方表格库如 tablewriter
	fmt.Printf("%+v\n", data)
	return nil
}

// StringSliceContains 检查字符串切片是否包含指定元素
func StringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveEmptyStrings 移除字符串切片中的空字符串
func RemoveEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

// TruncateString 截断字符串到指定长度
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
