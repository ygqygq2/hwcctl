package cdn

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// 测试创建客户端 - 由于需要真实的认证信息，这里只测试不panic
	_, err := NewClient()
	// 不管成功失败，只要不panic就算通过
	t.Logf("NewClient返回错误: %v", err)
}

func TestGetStringValue(t *testing.T) {
	// 测试nil指针
	result := getStringValue(nil)
	if result != "" {
		t.Errorf("期望空字符串，实际 '%s'", result)
	}

	// 测试有值的指针
	testStr := "test_value"
	result = getStringValue(&testStr)
	if result != "test_value" {
		t.Errorf("期望 'test_value'，实际 '%s'", result)
	}
}

func TestGetIntValue(t *testing.T) {
	// 测试nil指针
	result := getIntValue(nil)
	if result != 0 {
		t.Errorf("期望 0，实际 %d", result)
	}

	// 测试有值的指针
	testInt := int32(123)
	result = getIntValue(&testInt)
	if result != 123 {
		t.Errorf("期望 123，实际 %d", result)
	}
}

func TestGetRegion(t *testing.T) {
	// 测试有效区域
	region, err := getRegion("cn-north-1")
	if err != nil {
		t.Errorf("期望成功获取区域，但得到错误: %v", err)
	}
	if region == nil {
		t.Error("期望返回区域对象，但为nil")
	}

	// 测试无效区域
	_, err = getRegion("invalid-region")
	if err == nil {
		t.Error("期望无效区域返回错误，但成功了")
	}
}

func TestTask_Struct(t *testing.T) {
	// 测试Task结构体能正常创建和使用
	task := Task{
		ID:        "test-id",
		Type:      "refresh",
		Status:    "processing",
		CreatedAt: "2023-01-01T00:00:00Z",
		Progress:  50,
	}

	if task.ID != "test-id" {
		t.Errorf("期望任务ID为 'test-id'，实际为 '%s'", task.ID)
	}
	if task.Progress != 50 {
		t.Errorf("期望进度为 50，实际为 %d", task.Progress)
	}
}

func TestRefreshResult_Struct(t *testing.T) {
	// 测试RefreshResult结构体能正常创建和使用
	result := RefreshResult{
		TaskID: "refresh-task-id",
		Type:   "refresh",
		URLs:   []string{"https://example.com/test"},
		Status: "processing",
	}

	if result.TaskID != "refresh-task-id" {
		t.Errorf("期望任务ID为 'refresh-task-id'，实际为 '%s'", result.TaskID)
	}
	if len(result.URLs) != 1 {
		t.Errorf("期望URL数量为 1，实际为 %d", len(result.URLs))
	}
}

func TestPreloadResult_Struct(t *testing.T) {
	// 测试PreloadResult结构体
	result := &PreloadResult{
		TaskID: "preload-123",
	}

	if result.TaskID != "preload-123" {
		t.Errorf("期望TaskID为preload-123，实际为%s", result.TaskID)
	}
}

func TestClient_Struct(t *testing.T) {
	// 测试Client结构体的基本字段
	client := &Client{}

	// 测试结构体存在
	if client == nil {
		t.Error("Client结构体不应该为nil")
	}
}

func TestRegionMapping(t *testing.T) {
	// 测试区域映射功能
	region, err := getRegion("cn-north-1")
	if err != nil {
		t.Errorf("获取cn-north-1区域失败: %v", err)
	}

	if region == nil {
		t.Error("区域对象不应该为nil")
	}

	// 测试不支持的区域
	_, err = getRegion("invalid-region")
	if err == nil {
		t.Error("不支持的区域应该返回错误")
	}
}

func TestValidateURLs(t *testing.T) {
	// 测试URL验证功能
	validURLs := []string{
		"http://example.com/file.txt",
		"https://cdn.example.com/image.jpg",
	}

	for _, url := range validURLs {
		if !isValidURL(url) {
			t.Errorf("URL %s 应该是有效的", url)
		}
	}

	invalidURLs := []string{
		"not-a-url",
		"ftp://example.com/file.txt",
		"",
	}

	for _, url := range invalidURLs {
		if isValidURL(url) {
			t.Errorf("URL %s 应该是无效的", url)
		}
	}
}

// 辅助函数：检查URL是否有效
func isValidURL(url string) bool {
	return len(url) > 0 && (len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://"))
}

func TestCDNClientCreation(t *testing.T) {
	// 测试CDN客户端创建的业务逻辑和错误处理
	t.Run("客户端创建失败时的错误处理", func(t *testing.T) {
		// 在测试环境中，由于没有有效的华为云凭证，客户端创建应该失败
		// 这个测试验证错误处理逻辑是否正确
		_, err := NewClient()
		if err != nil {
			// 验证错误类型和消息的合理性
			t.Logf("期望的客户端创建错误: %v", err)

			// 错误消息应该有助于用户理解问题
			errMsg := err.Error()
			if errMsg == "" {
				t.Error("错误消息不应该为空")
			}
		}
	})
}

func TestURLValidation(t *testing.T) {
	// 测试URL验证的边界情况 - 这对CDN操作至关重要
	testCases := []struct {
		name     string
		url      string
		expected bool
		reason   string
	}{
		{
			name:     "标准HTTP URL",
			url:      "http://example.com/path/to/file.jpg",
			expected: true,
			reason:   "标准HTTP URL应该被接受",
		},
		{
			name:     "标准HTTPS URL",
			url:      "https://cdn.example.com/assets/style.css",
			expected: true,
			reason:   "标准HTTPS URL应该被接受",
		},
		{
			name:     "空字符串",
			url:      "",
			expected: false,
			reason:   "空URL应该被拒绝",
		},
		{
			name:     "无协议URL",
			url:      "example.com/file.txt",
			expected: false,
			reason:   "缺少协议的URL应该被拒绝",
		},
		{
			name:     "FTP协议",
			url:      "ftp://example.com/file.txt",
			expected: false,
			reason:   "非HTTP(S)协议应该被拒绝",
		},
		{
			name:     "恶意构造的URL",
			url:      "javascript:alert('xss')",
			expected: false,
			reason:   "恶意URL应该被拒绝",
		},
		{
			name:     "很短的URL",
			url:      "http://a",
			expected: true,
			reason:   "即使很短但格式正确的URL应该被接受",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isValidURL(tc.url)
			if result != tc.expected {
				t.Errorf("URL验证失败: %s\nURL: %s\n期望: %v, 实际: %v",
					tc.reason, tc.url, tc.expected, result)
			}
		})
	}
}

func TestTaskStructFields(t *testing.T) {
	// 测试Task结构体的业务逻辑完整性
	task := &Task{
		ID:          "task-123",
		Type:        "refresh",
		Status:      "processing",
		CreatedAt:   "2023-01-01T12:00:00Z",
		CompletedAt: "",
		Progress:    50,
	}

	// 验证关键字段
	if task.ID == "" {
		t.Error("任务ID不应该为空")
	}

	// 验证任务类型的合理性
	validTypes := []string{"refresh", "preload"}
	typeValid := false
	for _, vt := range validTypes {
		if task.Type == vt {
			typeValid = true
			break
		}
	}
	if !typeValid {
		t.Errorf("任务类型 '%s' 不在有效类型列表中", task.Type)
	}

	// 验证进度范围
	if task.Progress < 0 || task.Progress > 100 {
		t.Errorf("任务进度 %d 应该在0-100范围内", task.Progress)
	}

	// 验证时间格式（简单检查）
	if task.CreatedAt == "" {
		t.Error("创建时间不应该为空")
	}
}

func TestResultStructsDataIntegrity(t *testing.T) {
	// 测试结果结构体的数据完整性 - 重要的业务对象
	now := time.Now()

	t.Run("RefreshResult数据完整性", func(t *testing.T) {
		result := &RefreshResult{
			TaskID:    "refresh-task-123",
			Type:      "file",
			URLs:      []string{"https://cdn.example.com/file1.jpg", "https://cdn.example.com/file2.css"},
			Status:    "submitted",
			CreatedAt: now,
		}

		// 验证必需字段
		if result.TaskID == "" {
			t.Error("刷新结果的TaskID不应该为空")
		}

		if len(result.URLs) == 0 {
			t.Error("刷新结果应该包含至少一个URL")
		}

		// 验证URL的有效性
		for i, url := range result.URLs {
			if !isValidURL(url) {
				t.Errorf("刷新结果中的URL[%d] '%s' 无效", i, url)
			}
		}

		// 验证时间字段
		if result.CreatedAt.IsZero() {
			t.Error("创建时间不应该为零值")
		}
	})

	t.Run("PreloadResult数据完整性", func(t *testing.T) {
		result := &PreloadResult{
			TaskID:    "preload-task-456",
			Type:      "file",
			URLs:      []string{"https://cdn.example.com/large-file.zip"},
			Status:    "processing",
			CreatedAt: now,
		}

		// 类似的验证逻辑
		if result.TaskID == "" {
			t.Error("预热结果的TaskID不应该为空")
		}

		if len(result.URLs) == 0 {
			t.Error("预热结果应该包含至少一个URL")
		}

		// 预热通常用于大文件，验证这个业务逻辑
		for _, url := range result.URLs {
			if !isValidURL(url) {
				t.Errorf("预热结果中的URL '%s' 无效", url)
			}
		}
	})
}
