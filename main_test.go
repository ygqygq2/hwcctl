package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	// 这是一个基础测试，确保 main 包可以正常导入
	if testing.Short() {
		t.Skip("跳过 main 包测试")
	}
}
