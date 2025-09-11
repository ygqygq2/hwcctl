package cdn

// Client CDN 客户端封装
type Client struct{}

// Task 任务信息
type Task struct {
	ID          string
	Type        string
	Status      string
	CreatedAt   string
	CompletedAt string
	Progress    int
}

// NewClient 创建新的 CDN 客户端
func NewClient() (*Client, error) {
	return &Client{}, nil
}

// RefreshCache 刷新 CDN 缓存
func (c *Client) RefreshCache(urls []string, refreshType string) (string, error) {
	return "task-123456", nil
}

// PreloadCache 预热 CDN 缓存
func (c *Client) PreloadCache(urls []string) (string, error) {
	return "task-789012", nil
}

// GetTaskStatus 查询任务状态
func (c *Client) GetTaskStatus(taskId string) (*Task, error) {
	return &Task{
		ID:          taskId,
		Type:        "refresh",
		Status:      "task_inprogress",
		CreatedAt:   "2025-09-11 10:00:00",
		CompletedAt: "",
		Progress:    50,
	}, nil
}
