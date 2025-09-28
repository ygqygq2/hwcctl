package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter("json")
	if formatter == nil {
		t.Error("期望创建格式化器成功，但返回nil")
	}
}

func TestFormatter_PrintJSON(t *testing.T) {
	formatter := NewFormatter("json")
	testData := map[string]string{"key": "value"}

	output := captureOutput(func() {
		formatter.Print(testData)
	})

	if !strings.Contains(output, `"key": "value"`) {
		t.Errorf("JSON输出格式不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintYAML(t *testing.T) {
	formatter := NewFormatter("yaml")
	testData := map[string]string{"key": "value"}

	output := captureOutput(func() {
		formatter.Print(testData)
	})

	if !strings.Contains(output, "key: value") {
		t.Errorf("YAML输出格式不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintTable(t *testing.T) {
	formatter := NewFormatter("table")
	testData := map[string]interface{}{"name": "test", "value": 123}

	// 只测试不panic，不检查具体输出内容
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("表格输出发生panic: %v", r)
		}
	}()

	formatter.Print(testData)
	// 如果没有panic，就算通过
}

func TestFormatter_PrintText(t *testing.T) {
	formatter := NewFormatter("text")
	testData := "test data"

	output := captureOutput(func() {
		formatter.Print(testData)
	})

	if !strings.Contains(output, "test data") {
		t.Errorf("文本输出格式不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintSuccess(t *testing.T) {
	formatter := NewFormatter("text")

	output := captureOutput(func() {
		formatter.PrintSuccess("操作成功")
	})

	if !strings.Contains(output, "✅") || !strings.Contains(output, "操作成功") {
		t.Errorf("成功消息输出不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintError(t *testing.T) {
	formatter := NewFormatter("text")

	output := captureOutput(func() {
		formatter.PrintError("操作失败")
	})

	if !strings.Contains(output, "❌") || !strings.Contains(output, "操作失败") {
		t.Errorf("错误消息输出不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintWarning(t *testing.T) {
	formatter := NewFormatter("text")

	output := captureOutput(func() {
		formatter.PrintWarning("警告信息")
	})

	if !strings.Contains(output, "⚠️") || !strings.Contains(output, "警告信息") {
		t.Errorf("警告消息输出不正确，实际输出: %s", output)
	}
}

func TestFormatter_PrintInfo(t *testing.T) {
	formatter := NewFormatter("text")

	output := captureOutput(func() {
		formatter.PrintInfo("提示信息")
	})

	if !strings.Contains(output, "ℹ️") || !strings.Contains(output, "提示信息") {
		t.Errorf("信息消息输出不正确，实际输出: %s", output)
	}
}

func TestFormatter_InvalidFormat(t *testing.T) {
	// 测试无效格式的处理
	formatter := NewFormatter("invalid-format")

	output := captureOutput(func() {
		formatter.Print("test data")
	})

	// 对于无效格式，应该有某种输出
	if len(output) == 0 {
		t.Error("即使格式无效，也应该有输出")
	}
}

func TestFormatter_NilData(t *testing.T) {
	// 测试nil数据的处理
	formatter := NewFormatter("json")

	output := captureOutput(func() {
		formatter.Print(nil)
	})

	if len(output) == 0 {
		t.Error("即使数据为nil，也应该有输出")
	}
}

func TestFormatter_ComplexData(t *testing.T) {
	// 测试复杂数据结构
	formatter := NewFormatter("json")

	complexData := map[string]interface{}{
		"string": "test",
		"number": 123,
		"bool":   true,
		"array":  []string{"a", "b", "c"},
		"nested": map[string]string{"key": "value"},
	}

	output := captureOutput(func() {
		formatter.Print(complexData)
	})

	if !strings.Contains(output, "test") {
		t.Error("复杂数据结构输出应该包含测试数据")
	}
}
