# 04. CDN 管理

华为云内容分发网络 (CDN) 管理功能，支持缓存刷新、内容预热和任务查询。

## 功能概述

| 功能     | 命令                 | 说明                      |
| -------- | -------------------- | ------------------------- |
| 缓存刷新 | `hwcctl cdn refresh` | 刷新指定 URL 或目录的缓存 |
| 内容预热 | `hwcctl cdn preload` | 预热内容到边缘节点        |
| 任务查询 | `hwcctl cdn task`    | 查询刷新/预热任务状态     |

## 缓存刷新

### 基本用法

```bash
# 刷新单个文件
hwcctl cdn refresh --urls "https://example.com/image.jpg"

# 刷新多个文件
hwcctl cdn refresh --urls "https://example.com/file1.jpg,https://example.com/file2.css,https://example.com/file3.js"

# 使用数组格式
hwcctl cdn refresh --urls "https://example.com/file1.jpg" --urls "https://example.com/file2.css"
```

### 刷新类型

#### 文件刷新（默认）

```bash
hwcctl cdn refresh --type url --urls "https://example.com/specific-file.jpg"
```

#### 目录刷新

```bash
hwcctl cdn refresh --type directory --urls "https://example.com/images/"
```

### 批量刷新示例

```bash
# 刷新网站的静态资源
hwcctl cdn refresh --urls "https://cdn.example.com/css/,https://cdn.example.com/js/,https://cdn.example.com/images/" --type directory

# 刷新特定文件列表
hwcctl cdn refresh --urls "https://cdn.example.com/index.html,https://cdn.example.com/main.css,https://cdn.example.com/app.js"
```

## 内容预热

### 基本用法

```bash
# 预热单个文件
hwcctl cdn preload --urls "https://example.com/popular-video.mp4"

# 预热多个文件
hwcctl cdn preload --urls "https://example.com/hot-image1.jpg,https://example.com/hot-image2.jpg"
```

### 使用场景

1. **新内容发布**：预热新发布的热门内容
2. **活动准备**：活动前预热相关资源
3. **性能优化**：预热用户经常访问的资源

```bash
# 活动前预热
hwcctl cdn preload --urls "https://cdn.example.com/promo/banner.jpg,https://cdn.example.com/promo/video.mp4"
```

## 任务查询

### 查询任务状态

```bash
# 查询特定任务
hwcctl cdn task task-123456789

# 使用调试模式查看详细信息
hwcctl cdn task task-123456789 --debug
```

### 任务状态说明

| 状态             | 说明       |
| ---------------- | ---------- |
| `task_inprocess` | 任务处理中 |
| `task_done`      | 任务完成   |
| `task_fail`      | 任务失败   |

## 输出格式

### 表格格式（默认）

```bash
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output table
```

### JSON 格式

```bash
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output json
```

### YAML 格式

```bash
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output yaml
```

## 调试和监控

### 启用调试模式

```bash
hwcctl cdn refresh --urls "https://example.com/test.jpg" --debug
```

调试模式会显示：

- 详细的 API 请求信息
- 认证过程
- 错误详情
- 网络请求时间

### 详细输出

```bash
hwcctl cdn refresh --urls "https://example.com/test.jpg" --verbose
```

## 最佳实践

### 1. 批量操作

```bash
# 推荐：一次刷新多个相关文件
hwcctl cdn refresh --urls "https://cdn.example.com/v1.2.0/main.css,https://cdn.example.com/v1.2.0/main.js"

# 避免：多次单独刷新
hwcctl cdn refresh --urls "https://cdn.example.com/v1.2.0/main.css"
hwcctl cdn refresh --urls "https://cdn.example.com/v1.2.0/main.js"
```

### 2. 目录 vs 文件刷新

```bash
# 目录刷新：适合批量更新
hwcctl cdn refresh --type directory --urls "https://cdn.example.com/assets/"

# 文件刷新：适合精确控制
hwcctl cdn refresh --type url --urls "https://cdn.example.com/assets/main.css"
```

### 3. 预热策略

```bash
# 发布新版本时的推荐流程：
# 1. 先刷新旧版本缓存
hwcctl cdn refresh --type directory --urls "https://cdn.example.com/v1.1.0/"

# 2. 再预热新版本内容
hwcctl cdn preload --urls "https://cdn.example.com/v1.2.0/main.css,https://cdn.example.com/v1.2.0/main.js"
```

### 4. 监控任务执行

```bash
# 获取任务 ID 并查询状态
TASK_ID=$(hwcctl cdn refresh --urls "https://example.com/test.jpg" --output json | jq -r '.task_id')
hwcctl cdn task "$TASK_ID"
```

## 配置要求

CDN 功能需要以下配置：

1. **Domain ID**：CDN 服务必需的账号标识
2. **适当的权限**：确保 Access Key 具有 CDN 操作权限
3. **已配置域名**：URL 必须是已在华为云 CDN 配置的域名

## 限制和配额

华为云 CDN API 限制：

- **URL 数量**：单次最多 100 个 URL
- **请求频率**：建议间隔 1 秒以上
- **URL 长度**：最大 1000 字符

## 错误处理

常见错误及解决方案：

### 认证错误

```bash
Error: 创建认证信息失败
```

**解决方案**：检查 Domain ID 配置，参考 [配置指南](./03-configuration.md)

### 权限错误

```bash
Error: [HTTP403] Forbidden
```

**解决方案**：确保 Access Key 具有 CDN 服务权限

### 网络超时

```bash
Error: context deadline exceeded
```

**解决方案**：检查网络连接，或启用重试机制

更多故障排查请参考 [故障排查指南](./08-troubleshooting.md)。
