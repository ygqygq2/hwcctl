package utils

import (
	"testing"
)

func TestStringSliceContains(t *testing.T) {
	slice := []string{"apple", "banana", "orange"}

	tests := []struct {
		item     string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := StringSliceContains(slice, tt.item)
			if result != tt.expected {
				t.Errorf("StringSliceContains(%v, %s) = %v，期望 %v", slice, tt.item, result, tt.expected)
			}
		})
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"a", "", "b", "  ", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{"", "  ", "\t"},
			expected: []string{},
		},
		{
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
	}

	for i, tt := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := RemoveEmptyStrings(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("长度不匹配：期望 %d，实际 %d", len(tt.expected), len(result))
				return
			}
			for j, v := range result {
				if v != tt.expected[j] {
					t.Errorf("索引 %d：期望 %s，实际 %s", j, tt.expected[j], v)
				}
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello world", 5, "hello..."},
		{"hello", 10, "hello"},
		{"hello world", 11, "hello world"},
		{"", 5, ""},
		{"test", 0, "..."},
	}

	for i, tt := range tests {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			result := TruncateString(tt.input, tt.length)
			if result != tt.expected {
				t.Errorf("TruncateString(%s, %d) = %s，期望 %s", tt.input, tt.length, result, tt.expected)
			}
		})
	}
}

func TestFormatOutput(t *testing.T) {
	testData := map[string]interface{}{
		"name":   "test",
		"value":  123,
		"active": true,
	}

	tests := []struct {
		format    string
		shouldErr bool
	}{
		{"json", false},
		{"yaml", false},
		{"table", false},
		{"invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			err := FormatOutput(testData, tt.format)
			if tt.shouldErr && err == nil {
				t.Error("期望出错，但没有错误")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("期望成功，实际出错: %v", err)
			}
		})
	}
}
