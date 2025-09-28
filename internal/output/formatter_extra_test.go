package output

import (
	"testing"
)

func TestPrintSliceAsTable(t *testing.T) {
	formatter := NewFormatter("table")

	// 测试字符串切片
	stringSlice := []string{"item1", "item2", "item3"}

	err := formatter.Print(stringSlice)
	if err != nil {
		t.Errorf("打印字符串切片失败: %v", err)
	}
}

func TestPrintStructAsTable(t *testing.T) {
	formatter := NewFormatter("table")

	// 测试结构体
	type testStruct struct {
		Name string `table:"姓名"`
		Age  int    `table:"年龄"`
	}

	testData := testStruct{Name: "John", Age: 30}

	err := formatter.Print(testData)
	if err != nil {
		t.Errorf("打印结构体失败: %v", err)
	}
}

func TestPrintComplexSliceAsTable(t *testing.T) {
	formatter := NewFormatter("table")

	// 测试复杂类型切片
	type item struct {
		ID   int    `table:"编号"`
		Name string `table:"名称"`
	}

	items := []item{
		{ID: 1, Name: "Item1"},
		{ID: 2, Name: "Item2"},
	}

	err := formatter.Print(items)
	if err != nil {
		t.Errorf("打印复杂切片失败: %v", err)
	}
}

func TestFormatterEdgeCases(t *testing.T) {
	// 测试空数据
	formatter := NewFormatter("table")

	// 测试空切片
	emptySlice := []string{}
	err := formatter.Print(emptySlice)
	if err != nil {
		t.Errorf("打印空切片失败: %v", err)
	}

	// 测试nil值
	err = formatter.Print(nil)
	if err != nil {
		t.Errorf("打印nil失败: %v", err)
	}

	// 测试无效格式 - 应该回退到table格式
	invalidFormatter := NewFormatter("invalid")
	err = invalidFormatter.Print("test")
	if err != nil {
		t.Errorf("无效格式处理失败: %v", err)
	}
}

func TestFormatterMessages(t *testing.T) {
	formatter := NewFormatter("table")

	// 测试各种消息类型
	formatter.PrintSuccess("测试成功")
	formatter.PrintError("测试错误")
	formatter.PrintWarning("测试警告")
	formatter.PrintInfo("测试信息")
}

func TestFormatterMapAsTable(t *testing.T) {
	formatter := NewFormatter("table")

	// 测试map输出
	testMap := map[string]interface{}{
		"name": "test",
		"age":  25,
		"city": "Beijing",
	}

	err := formatter.Print(testMap)
	if err != nil {
		t.Errorf("打印map失败: %v", err)
	}
}

func TestFormatterDifferentFormats(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	// 测试JSON格式
	jsonFormatter := NewFormatter("json")
	err := jsonFormatter.Print(testData)
	if err != nil {
		t.Errorf("JSON格式输出失败: %v", err)
	}

	// 测试YAML格式
	yamlFormatter := NewFormatter("yaml")
	err = yamlFormatter.Print(testData)
	if err != nil {
		t.Errorf("YAML格式输出失败: %v", err)
	}

	// 测试Text格式
	textFormatter := NewFormatter("text")
	err = textFormatter.Print("simple text")
	if err != nil {
		t.Errorf("Text格式输出失败: %v", err)
	}
}
