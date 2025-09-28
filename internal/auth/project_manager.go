package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ProjectManager 项目管理器，负责统一管理项目ID
type ProjectManager struct {
	config    *Config
	projectID string
	mutex     sync.RWMutex
	loaded    bool
}

// 全局项目管理器实例
var (
	globalProjectManager *ProjectManager
	projectManagerOnce   sync.Once
)

// GetProjectManager 获取全局项目管理器实例
func GetProjectManager() *ProjectManager {
	projectManagerOnce.Do(func() {
		globalProjectManager = &ProjectManager{}
	})
	return globalProjectManager
}

// InitWithConfig 使用配置初始化项目管理器
func (pm *ProjectManager) InitWithConfig(config *Config) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.config = config
	pm.loaded = false // 重置加载状态
}

// GetProjectID 获取项目ID，支持懒加载
func (pm *ProjectManager) GetProjectID() (string, error) {
	// 先尝试读取缓存
	pm.mutex.RLock()
	if pm.loaded && pm.projectID != "" {
		projectID := pm.projectID
		pm.mutex.RUnlock()
		return projectID, nil
	}
	pm.mutex.RUnlock()

	// 需要加载项目ID
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// 双重检查，避免并发时重复加载
	if pm.loaded && pm.projectID != "" {
		return pm.projectID, nil
	}

	// 检查配置是否已初始化
	if pm.config == nil {
		return "", fmt.Errorf("项目管理器未初始化")
	}

	// 1. 优先使用配置文件中的 project_id
	if pm.config.ProjectID != "" {
		pm.projectID = pm.config.ProjectID
		pm.loaded = true
		return pm.projectID, nil
	}

	// 2. 尝试使用 enterprise_project_id
	if pm.config.EnterpriseProjectID != "" && pm.config.EnterpriseProjectID != "0" {
		pm.projectID = pm.config.EnterpriseProjectID
		pm.loaded = true
		return pm.projectID, nil
	}

	// 3. 如果都没有，尝试从API获取
	if pm.config.AccessKey != "" && pm.config.SecretKey != "" && pm.config.DomainID != "" {
		projectID, err := pm.fetchProjectIDFromAPI()
		if err != nil {
			// API获取失败，使用默认值
			pm.projectID = "0" // 默认企业项目
			pm.loaded = true
			return pm.projectID, nil
		}
		pm.projectID = projectID
		pm.loaded = true
		return pm.projectID, nil
	}

	// 4. 最后使用默认值
	pm.projectID = "0"
	pm.loaded = true
	return pm.projectID, nil
}

// fetchProjectIDFromAPI 从华为云API获取项目ID
func (pm *ProjectManager) fetchProjectIDFromAPI() (string, error) {
	url := "https://iam.myhuaweicloud.com/v3/projects"

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 这里需要实现华为云签名，暂时简化处理
	// TODO: 实现完整的华为云签名算法

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var projectsResponse ProjectsResponse
	err = json.Unmarshal(body, &projectsResponse)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 查找当前区域对应的项目ID
	for _, project := range projectsResponse.Projects {
		if project.Name == pm.config.Region && project.Enabled {
			return project.ID, nil
		}
	}

	return "", fmt.Errorf("在区域 %s 中未找到启用的项目", pm.config.Region)
}

// RefreshProjectID 强制刷新项目ID缓存
func (pm *ProjectManager) RefreshProjectID() error {
	pm.mutex.Lock()
	pm.loaded = false
	pm.projectID = ""
	pm.mutex.Unlock()

	// 重新获取（不能在锁内调用GetProjectID，会死锁）
	_, err := pm.GetProjectID()
	return err
} // IsLoaded 检查项目ID是否已加载
func (pm *ProjectManager) IsLoaded() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.loaded
}

// GetCachedProjectID 获取缓存的项目ID，不触发懒加载
func (pm *ProjectManager) GetCachedProjectID() string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.projectID
}
