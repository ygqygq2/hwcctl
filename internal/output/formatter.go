package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format 输出格式枚举
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatText  Format = "text"
)

// Formatter 输出格式化器
type Formatter struct {
	format Format
	writer *tabwriter.Writer
}

// NewFormatter 创建新的格式化器
func NewFormatter(format string) *Formatter {
	f := &Formatter{
		format: Format(format),
		writer: tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0),
	}
	return f
}

// Print 根据格式输出数据
func (f *Formatter) Print(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.printJSON(data)
	case FormatYAML:
		return f.printYAML(data)
	case FormatTable:
		return f.printTable(data)
	case FormatText:
		return f.printText(data)
	default:
		return f.printTable(data)
	}
}

// printJSON 输出 JSON 格式
func (f *Formatter) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML 输出 YAML 格式
func (f *Formatter) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(data)
}

// printTable 输出表格格式
func (f *Formatter) printTable(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return f.printSliceAsTable(v)
	case reflect.Struct:
		return f.printStructAsTable(v)
	case reflect.Map:
		return f.printMapAsTable(v)
	default:
		fmt.Println(data)
		return nil
	}
}

// printSliceAsTable 打印切片为表格
func (f *Formatter) printSliceAsTable(v reflect.Value) error {
	if v.Len() == 0 {
		fmt.Println("No data found.")
		return nil
	}

	// 获取第一个元素的类型来确定列名
	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	if first.Kind() == reflect.Struct {
		// 打印表头
		t := first.Type()
		headers := make([]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("table")
			if tag != "" {
				headers[i] = tag
			} else {
				headers[i] = strings.ToUpper(field.Name)
			}
		}
		fmt.Fprintln(f.writer, strings.Join(headers, "\t"))

		// 打印数据行
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			row := make([]string, item.NumField())
			for j := 0; j < item.NumField(); j++ {
				field := item.Field(j)
				row[j] = fmt.Sprintf("%v", field.Interface())
			}
			fmt.Fprintln(f.writer, strings.Join(row, "\t"))
		}
	} else {
		// 简单类型的切片
		for i := 0; i < v.Len(); i++ {
			fmt.Println(v.Index(i).Interface())
		}
	}

	return f.writer.Flush()
}

// printStructAsTable 打印结构体为表格
func (f *Formatter) printStructAsTable(v reflect.Value) error {
	t := v.Type()
	
	// 打印两列表格：字段名 | 值
	fmt.Fprintln(f.writer, "FIELD\tVALUE")
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		
		name := field.Name
		tag := field.Tag.Get("table")
		if tag != "" {
			name = tag
		}
		
		value := v.Field(i).Interface()
		fmt.Fprintf(f.writer, "%s\t%v\n", name, value)
	}
	
	return f.writer.Flush()
}

// printMapAsTable 打印 map 为表格
func (f *Formatter) printMapAsTable(v reflect.Value) error {
	fmt.Fprintln(f.writer, "KEY\tVALUE")
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		fmt.Fprintf(f.writer, "%v\t%v\n", key.Interface(), value.Interface())
	}
	return f.writer.Flush()
}

// printText 输出纯文本格式
func (f *Formatter) printText(data interface{}) error {
	fmt.Println(data)
	return nil
}

// PrintSuccess 打印成功消息
func (f *Formatter) PrintSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// PrintError 打印错误消息
func (f *Formatter) PrintError(message string) {
	fmt.Printf("❌ %s\n", message)
}

// PrintWarning 打印警告消息
func (f *Formatter) PrintWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// PrintInfo 打印信息消息
func (f *Formatter) PrintInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}
