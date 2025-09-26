package auth

import (
	"testing"
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
