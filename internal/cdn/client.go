package cdn

import (
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/ygqygq2/hwcctl/internal/auth"
	hwErrors "github.com/ygqygq2/hwcctl/internal/errors"
	"github.com/ygqygq2/hwcctl/internal/logx"
)

// Client CDN 客户端封装
type Client struct {
	cdnClient *cdn.CdnClient
	region    string
}

// Task 任务信息
type Task struct {
	ID          string `json:"id" table:"任务ID"`
	Type        string `json:"type" table:"任务类型"`
	Status      string `json:"status" table:"任务状态"`
	CreatedAt   string `json:"created_at" table:"创建时间"`
	CompletedAt string `json:"completed_at,omitempty" table:"完成时间"`
	Progress    int    `json:"progress" table:"进度(%)"`
}

// RefreshResult 刷新结果
type RefreshResult struct {
	TaskID    string    `json:"task_id" table:"任务ID"`
	Type      string    `json:"type" table:"类型"`
	URLs      []string  `json:"urls" table:"URL列表"`
	Status    string    `json:"status" table:"状态"`
	CreatedAt time.Time `json:"created_at" table:"创建时间"`
}

// PreloadResult 预热结果
type PreloadResult struct {
	TaskID    string    `json:"task_id" table:"任务ID"`
	Type      string    `json:"type" table:"类型"`
	URLs      []string  `json:"urls" table:"URL列表"`
	Status    string    `json:"status" table:"状态"`
	CreatedAt time.Time `json:"created_at" table:"创建时间"`
}

// NewClient 创建新的 CDN 客户端
func NewClient() (*Client, error) {
	// 获取认证信息
	creds, err := auth.GetCredentials()
	if err != nil {
		return nil, hwErrors.NewAuthError(fmt.Sprintf("获取认证信息失败: %v", err))
	}

	// 验证必需的认证信息
	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
		return nil, hwErrors.NewAuthError("Access Key ID 和 Secret Access Key 不能为空")
	}

	if creds.Region == "" {
		return nil, hwErrors.NewValidationError("区域不能为空")
	}

	// 创建认证对象 - 使用全局认证方法 (CDN 服务需要)
	credentialsBuilder := global.NewCredentialsBuilder().
		WithAk(creds.AccessKeyID).
		WithSk(creds.SecretAccessKey)
	
	// 如果配置了 Domain ID，则使用它
	if creds.DomainID != "" {
		credentialsBuilder = credentialsBuilder.WithDomainId(creds.DomainID)
	}
	
	authCredentials, err := credentialsBuilder.SafeBuild()
	if err != nil {
		return nil, hwErrors.NewAuthError(fmt.Sprintf("创建认证信息失败: %v", err))
	}

	// 获取区域对象
	regionObj, err := getRegion(creds.Region)
	if err != nil {
		return nil, hwErrors.NewValidationError(fmt.Sprintf("不支持的区域: %s", creds.Region))
	}

	// 创建客户端配置
	hcConfig := config.DefaultHttpConfig().
		WithIgnoreSSLVerification(false).
		WithTimeout(30)

	// 创建 CDN 客户端
	hcClient, err := cdn.CdnClientBuilder().
		WithRegion(regionObj).
		WithCredential(authCredentials).
		WithHttpConfig(hcConfig).
		SafeBuild()
	if err != nil {
		return nil, hwErrors.NewServerError(fmt.Sprintf("创建CDN客户端失败: %v", err))
	}

	cdnClient := cdn.NewCdnClient(hcClient)

	logx.Debugf("CDN 客户端创建成功，区域: %s", creds.Region)

	return &Client{
		cdnClient: cdnClient,
		region:    creds.Region,
	}, nil
}

// RefreshCache 刷新 CDN 缓存
func (c *Client) RefreshCache(urls []string, refreshType string) (string, error) {
	logx.Debugf("开始刷新 CDN 缓存，类型: %s, URLs: %v", refreshType, urls)

	// 构建刷新请求体
	refreshTaskBody := &model.RefreshTaskRequestBody{
		Urls: urls,
	}

	// 设置刷新类型
	switch refreshType {
	case "url", "file":
		typeEnum := model.GetRefreshTaskRequestBodyTypeEnum().FILE
		refreshTaskBody.Type = &typeEnum
	case "directory", "dir":
		typeEnum := model.GetRefreshTaskRequestBodyTypeEnum().DIRECTORY
		refreshTaskBody.Type = &typeEnum
	default:
		return "", hwErrors.NewValidationError(fmt.Sprintf("不支持的刷新类型: %s，支持的类型: url, directory", refreshType))
	}

	// 构建请求
	refreshTaskRequest := &model.RefreshTaskRequest{
		RefreshTask: refreshTaskBody,
	}

	request := &model.CreateRefreshTasksRequest{
		Body: refreshTaskRequest,
	}

	// 发送请求
	response, err := c.cdnClient.CreateRefreshTasks(request)
	if err != nil {
		logx.Errorf("刷新 CDN 缓存失败: %v", err)
		return "", hwErrors.ParseHuaweiCloudError(500, err.Error())
	}

	if response.RefreshTask == nil {
		return "", hwErrors.NewServerError("刷新任务响应为空")
	}

	taskID := *response.RefreshTask
	logx.Infof("CDN 缓存刷新任务创建成功，任务ID: %s", taskID)

	return taskID, nil
}

// PreloadCache 预热 CDN 缓存
func (c *Client) PreloadCache(urls []string) (string, error) {
	logx.Debugf("开始预热 CDN 缓存，URLs: %v", urls)

	// 构建预热请求体
	preheatingTaskBody := &model.PreheatingTaskRequestBody{
		Urls: urls,
	}

	// 构建请求
	preheatingTaskRequest := &model.PreheatingTaskRequest{
		PreheatingTask: preheatingTaskBody,
	}

	request := &model.CreatePreheatingTasksRequest{
		Body: preheatingTaskRequest,
	}

	// 发送请求
	response, err := c.cdnClient.CreatePreheatingTasks(request)
	if err != nil {
		logx.Errorf("预热 CDN 缓存失败: %v", err)
		return "", hwErrors.ParseHuaweiCloudError(500, err.Error())
	}

	if response.PreheatingTask == nil {
		return "", hwErrors.NewServerError("预热任务响应为空")
	}

	taskID := *response.PreheatingTask
	logx.Infof("CDN 缓存预热任务创建成功，任务ID: %s", taskID)

	return taskID, nil
}

// GetTaskStatus 查询任务状态
func (c *Client) GetTaskStatus(taskID string) (*Task, error) {
	logx.Debugf("查询任务状态，任务ID: %s", taskID)

	// 构建查询请求
	request := &model.ShowHistoryTasksRequest{}
	
	// 设置查询范围（最近7天）
	endTime := time.Now().Unix() * 1000
	startTime := time.Now().AddDate(0, 0, -7).Unix() * 1000
	request.StartDate = &startTime
	request.EndDate = &endTime

	// 发送请求
	response, err := c.cdnClient.ShowHistoryTasks(request)
	if err != nil {
		logx.Errorf("查询任务状态失败: %v", err)
		return nil, hwErrors.ParseHuaweiCloudError(500, err.Error())
	}

	if response.Tasks == nil || len(*response.Tasks) == 0 {
		return nil, hwErrors.NewNotFoundError(fmt.Sprintf("任务 %s", taskID))
	}

	// 查找指定的任务
	tasks := *response.Tasks
	for _, task := range tasks {
		if task.Id != nil && *task.Id == taskID {
			return convertToTask(&task), nil
		}
	}

	return nil, hwErrors.NewNotFoundError(fmt.Sprintf("任务 %s", taskID))
}

// convertToTask 转换华为云任务对象为内部任务对象
func convertToTask(hwTask *model.TasksObject) *Task {
	task := &Task{
		ID: getStringValue(hwTask.Id),
	}

	// 转换任务类型
	if hwTask.TaskType != nil {
		switch *hwTask.TaskType {
		case model.GetTasksObjectTaskTypeEnum().REFRESH:
			task.Type = "refresh"
		case model.GetTasksObjectTaskTypeEnum().PREHEATING:
			task.Type = "preload"
		default:
			task.Type = hwTask.TaskType.Value()
		}
	}

	// 转换任务状态
	if hwTask.Status != nil {
		switch *hwTask.Status {
		case "task_inprocess":
			task.Status = "进行中"
		case "task_done":
			task.Status = "已完成"
		case "task_failed":
			task.Status = "失败"
		default:
			task.Status = *hwTask.Status
		}
	}

	// 转换时间
	if hwTask.CreateTime != nil {
		task.CreatedAt = time.Unix(*hwTask.CreateTime/1000, 0).Format("2006-01-02 15:04:05")
	}

	// 计算进度
	if hwTask.Processing != nil && hwTask.Total != nil {
		total := *hwTask.Total
		succeed := getIntValue(hwTask.Succeed)
		failed := getIntValue(hwTask.Failed)
		
		if total > 0 {
			completed := succeed + failed
			task.Progress = int(float64(completed) / float64(total) * 100)
		}
	}

	return task
}

// getRegion 获取华为云区域对象
func getRegion(regionName string) (*region.Region, error) {
	// 华为云支持的区域映射
	regionMap := map[string]*region.Region{
		"cn-north-1":     region.NewRegion("cn-north-1", "https://cdn.myhuaweicloud.com"),     // 华北-北京一
		"cn-north-4":     region.NewRegion("cn-north-4", "https://cdn.myhuaweicloud.com"),     // 华北-北京四
		"cn-east-2":      region.NewRegion("cn-east-2", "https://cdn.myhuaweicloud.com"),      // 华东-上海二
		"cn-east-3":      region.NewRegion("cn-east-3", "https://cdn.myhuaweicloud.com"),      // 华东-上海一
		"cn-south-1":     region.NewRegion("cn-south-1", "https://cdn.myhuaweicloud.com"),     // 华南-广州
		"cn-southwest-2": region.NewRegion("cn-southwest-2", "https://cdn.myhuaweicloud.com"), // 西南-贵阳一
		"ap-southeast-1": region.NewRegion("ap-southeast-1", "https://cdn.ap-southeast-1.myhuaweicloud.com"), // 亚太-新加坡
		"ap-southeast-2": region.NewRegion("ap-southeast-2", "https://cdn.ap-southeast-2.myhuaweicloud.com"), // 亚太-悉尼
		"ap-southeast-3": region.NewRegion("ap-southeast-3", "https://cdn.ap-southeast-3.myhuaweicloud.com"), // 亚太-吉隆坡
		"af-south-1":     region.NewRegion("af-south-1", "https://cdn.af-south-1.myhuaweicloud.com"),     // 非洲-约翰内斯堡
		"na-mexico-1":    region.NewRegion("na-mexico-1", "https://cdn.na-mexico-1.myhuaweicloud.com"),    // 拉美-墨西哥城一
		"la-north-2":     region.NewRegion("la-north-2", "https://cdn.la-north-2.myhuaweicloud.com"),     // 拉美-墨西哥城二
		"sa-brazil-1":    region.NewRegion("sa-brazil-1", "https://cdn.sa-brazil-1.myhuaweicloud.com"),    // 拉美-圣保罗一
	}

	if regionObj, exists := regionMap[regionName]; exists {
		return regionObj, nil
	}

	return nil, fmt.Errorf("不支持的区域: %s", regionName)
}

// getStringValue 安全获取字符串指针的值
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// getIntValue 安全获取整数指针的值
func getIntValue(ptr *int32) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}
