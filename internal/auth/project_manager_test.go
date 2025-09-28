package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestProjectManager_GetProjectID(t *testing.T) {
	// 创建测试配置
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "test-project-id",
		EnterpriseProjectID: "0",
	}

	// 初始化项目管理器
	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	// 测试获取项目ID
	projectID, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	expectedID := "test-project-id"
	if projectID != expectedID {
		t.Errorf("期望项目ID: %s, 实际: %s", expectedID, projectID)
	}

	// 验证缓存状态
	if !pm.IsLoaded() {
		t.Error("项目管理器应该已加载")
	}

	cachedID := pm.GetCachedProjectID()
	if cachedID != expectedID {
		t.Errorf("缓存的项目ID不正确: %s", cachedID)
	}
}

func TestProjectManager_EnterpriseProjectID(t *testing.T) {
	// 测试使用企业项目ID的情况
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "", // 空项目ID
		EnterpriseProjectID: "enterprise-123",
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	projectID, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	expectedID := "enterprise-123"
	if projectID != expectedID {
		t.Errorf("期望项目ID: %s, 实际: %s", expectedID, projectID)
	}
}

func TestProjectManager_DefaultProjectID(t *testing.T) {
	// 测试默认项目ID的情况
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "",
		EnterpriseProjectID: "0", // 默认企业项目
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	projectID, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	expectedID := "0"
	if projectID != expectedID {
		t.Errorf("期望项目ID: %s, 实际: %s", expectedID, projectID)
	}
}

func TestProjectManager_RefreshProjectID(t *testing.T) {
	// 测试刷新项目ID功能
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "initial-id",
		EnterpriseProjectID: "0",
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	// 第一次获取
	projectID1, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	// 修改配置
	testConfig.ProjectID = "new-project-id"

	// 刷新缓存
	err = pm.RefreshProjectID()
	if err != nil {
		t.Fatalf("刷新项目ID失败: %v", err)
	}

	// 再次获取
	projectID2, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	if projectID1 == projectID2 {
		t.Error("刷新后项目ID应该不同")
	}

	if projectID2 != "new-project-id" {
		t.Errorf("期望刷新后的项目ID: new-project-id, 实际: %s", projectID2)
	}
}

func TestProjectManager_ConcurrentAccess(t *testing.T) {
	// 测试并发访问
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "concurrent-test-id",
		EnterpriseProjectID: "0",
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	// 启动多个goroutine同时获取项目ID
	results := make(chan string, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			projectID, err := pm.GetProjectID()
			if err != nil {
				errors <- err
				return
			}
			results <- projectID
		}()
	}

	// 收集结果
	var projectIDs []string
	for i := 0; i < 10; i++ {
		select {
		case id := <-results:
			projectIDs = append(projectIDs, id)
		case err := <-errors:
			t.Fatalf("并发获取项目ID失败: %v", err)
		}
	}

	// 验证所有结果都相同
	expectedID := "concurrent-test-id"
	for i, id := range projectIDs {
		if id != expectedID {
			t.Errorf("第%d个结果不正确: 期望 %s, 实际 %s", i, expectedID, id)
		}
	}
}

func TestGetUnifiedProjectID(t *testing.T) {
	// 测试全局函数
	// 先加载配置初始化项目管理器
	_, err := LoadConfig("", "", "", "")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	projectID, err := GetUnifiedProjectID()
	if err != nil {
		t.Fatalf("获取统一项目ID失败: %v", err)
	}

	if projectID == "" {
		t.Error("项目ID不应该为空")
	}

	t.Logf("获取到的统一项目ID: %s", projectID)
}

func TestProjectManager_UninitializedConfig(t *testing.T) {
	// 测试未初始化配置的情况
	pm := GetProjectManager()

	// 重置状态
	pm.mutex.Lock()
	pm.config = nil
	pm.loaded = false
	pm.projectID = ""
	pm.mutex.Unlock()

	_, err := pm.GetProjectID()
	if err == nil {
		t.Error("期望返回错误，但没有错误")
	}

	expectedError := "项目管理器未初始化"
	if err.Error() != expectedError {
		t.Errorf("期望错误: %s, 实际: %s", expectedError, err.Error())
	}
}

func TestProjectManager_EmptyConfig(t *testing.T) {
	// 测试空配置的情况
	testConfig := &Config{
		AccessKey:           "",
		SecretKey:           "",
		Region:              "",
		DomainID:            "",
		ProjectID:           "",
		EnterpriseProjectID: "",
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	projectID, err := pm.GetProjectID()
	if err != nil {
		t.Fatalf("获取项目ID失败: %v", err)
	}

	// 应该返回默认值
	expectedID := "0"
	if projectID != expectedID {
		t.Errorf("期望项目ID: %s, 实际: %s", expectedID, projectID)
	}
}

func TestProjectManager_PriorityOrder(t *testing.T) {
	// 测试优先级顺序：project_id > enterprise_project_id (非"0") > 默认值"0"
	testCases := []struct {
		name                string
		projectID           string
		enterpriseProjectID string
		expectedID          string
	}{
		{
			name:                "ProjectID优先",
			projectID:           "project-123",
			enterpriseProjectID: "enterprise-456",
			expectedID:          "project-123",
		},
		{
			name:                "EnterpriseProjectID次优",
			projectID:           "",
			enterpriseProjectID: "enterprise-456",
			expectedID:          "enterprise-456",
		},
		{
			name:                "EnterpriseProjectID为0时使用默认",
			projectID:           "",
			enterpriseProjectID: "0",
			expectedID:          "0",
		},
		{
			name:                "都为空时使用默认",
			projectID:           "",
			enterpriseProjectID: "",
			expectedID:          "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testConfig := &Config{
				AccessKey:           "test-ak",
				SecretKey:           "test-sk",
				Region:              "cn-north-4",
				DomainID:            "test-domain-id",
				ProjectID:           tc.projectID,
				EnterpriseProjectID: tc.enterpriseProjectID,
			}

			pm := GetProjectManager()
			pm.InitWithConfig(testConfig)

			projectID, err := pm.GetProjectID()
			if err != nil {
				t.Fatalf("获取项目ID失败: %v", err)
			}

			if projectID != tc.expectedID {
				t.Errorf("期望项目ID: %s, 实际: %s", tc.expectedID, projectID)
			}
		})
	}
}

func TestProjectManager_ThreadSafety(t *testing.T) {
	// 测试线程安全性 - 混合读写操作
	testConfig := &Config{
		AccessKey:           "test-ak",
		SecretKey:           "test-sk",
		Region:              "cn-north-4",
		DomainID:            "test-domain-id",
		ProjectID:           "thread-safety-test",
		EnterpriseProjectID: "0",
	}

	pm := GetProjectManager()
	pm.InitWithConfig(testConfig)

	done := make(chan bool, 20)
	errors := make(chan error, 20)

	// 启动多个读取操作
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				projectID, err := pm.GetProjectID()
				if err != nil {
					errors <- fmt.Errorf("读取操作%d失败: %v", id, err)
					return
				}
				if projectID != "thread-safety-test" {
					errors <- fmt.Errorf("读取操作%d得到错误的ID: %s", id, projectID)
					return
				}
			}
		}(i)
	}

	// 启动多个状态检查操作
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()
			for j := 0; j < 10; j++ {
				_ = pm.IsLoaded()
				_ = pm.GetCachedProjectID()
			}
		}(i)
	}

	// 启动少量刷新操作
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()
			time.Sleep(time.Millisecond * 10) // 稍微延迟以让其他操作先开始
			err := pm.RefreshProjectID()
			if err != nil {
				errors <- fmt.Errorf("刷新操作%d失败: %v", id, err)
				return
			}
		}(i)
	}

	// 等待所有操作完成
	for i := 0; i < 20; i++ {
		select {
		case <-done:
			// 操作完成
		case err := <-errors:
			t.Error(err)
		case <-time.After(5 * time.Second):
			t.Fatal("测试超时")
		}
	}
}
