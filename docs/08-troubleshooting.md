# 08. 故障排查

常见问题和解决方案指南。

## 认证相关问题

### Domain ID 错误

**问题**：CDN 操作失败，提示需要 Domain ID

```
Error: 创建认证信息失败: CDN 服务需要提供 Domain ID
```

**解决方案**：

1. 获取 Domain ID：

   - 登录华为云控制台
   - 右上角用户名 → "我的凭证"
   - 查看"账号 ID"（即 Domain ID）

2. 配置 Domain ID：

   ```bash
   # 方式一：环境变量
   export HUAWEICLOUD_DOMAIN_ID="your-domain-id"

   # 方式二：配置文件
   echo 'domain_id: "your-domain-id"' >> ~/.hwcctl/config

   # 方式三：命令行参数
   hwcctl --domain-id "your-domain-id" cdn refresh --urls "..."
   ```

### Access Key 权限不足

**问题**：操作被拒绝

```
Error: [HTTP403] Forbidden
```

**解决方案**：

1. 确认 Access Key 具有以下权限：

   - CDN FullAccess 或
   - CDN ReadOnlyAccess + CDN RefreshCache + CDN PreloadCache

2. 检查 IAM 用户权限配置
3. 确认 Access Key 和 Secret Key 正确

### 网络连接问题

**问题**：请求超时

```
Error: context deadline exceeded (Client.Timeout exceeded while awaiting headers)
```

**解决方案**：

1. 检查网络连接：

   ```bash
   curl -I https://cdn.myhuaweicloud.com
   ```

2. 检查防火墙设置
3. 启用重试机制：
   ```yaml
   # ~/.hwcctl/config
   default:
     enable_retry: true
     max_retries: 3
   ```

## 配置相关问题

### 配置文件不生效

**问题**：手动创建的配置文件不生效

**解决方案**：

1. 检查文件路径：`~/.hwcctl/config`
2. 检查文件格式（YAML）
3. 检查文件权限：

   ```bash
   chmod 600 ~/.hwcctl/config
   ```

4. 验证配置格式：
   ```bash
   # 使用调试模式查看配置加载
   hwcctl --debug version
   ```

### 环境变量未加载

**问题**：设置了环境变量但不生效

**解决方案**：

1. 检查环境变量名称：

   ```bash
   echo $HUAWEICLOUD_ACCESS_KEY
   echo $HUAWEICLOUD_SECRET_KEY
   echo $HUAWEICLOUD_REGION
   echo $HUAWEICLOUD_DOMAIN_ID
   ```

2. 确保在当前 shell 中设置：

   ```bash
   # 临时设置
   export HUAWEICLOUD_ACCESS_KEY="your-key"

   # 永久设置
   echo 'export HUAWEICLOUD_ACCESS_KEY="your-key"' >> ~/.bashrc
   source ~/.bashrc
   ```

## CDN 相关问题

### URL 不在 CDN 配置中

**问题**：刷新或预热失败

```
Error: [HTTP404] Not Found
```

**解决方案**：

1. 确认域名已在华为云 CDN 控制台配置
2. 确认域名状态为"已上线"
3. 检查 URL 格式是否正确（必须包含 https:// 或 http://）

### 批量操作失败

**问题**：批量刷新部分失败

**解决方案**：

1. 检查 URL 数量限制（单次最多 100 个）
2. 检查单个 URL 长度（最大 1000 字符）
3. 分批处理：
   ```bash
   # 分批刷新
   hwcctl cdn refresh --urls "url1,url2,url3"
   sleep 2
   hwcctl cdn refresh --urls "url4,url5,url6"
   ```

## 性能相关问题

### 请求响应慢

**解决方案**：

1. 启用重试但减少重试次数：

   ```yaml
   default:
     enable_retry: true
     max_retries: 1
   ```

2. 使用批量操作减少 API 调用次数
3. 检查网络延迟：
   ```bash
   ping cdn.myhuaweicloud.com
   ```

### 内存使用过高

**解决方案**：

1. 减少批量操作的 URL 数量
2. 分批处理大量 URL
3. 监控系统资源使用

## 调试技巧

### 启用详细日志

```bash
# 调试模式
hwcctl --debug cdn refresh --urls "..."

# 详细输出
hwcctl --verbose cdn refresh --urls "..."
```

### 查看配置加载过程

```bash
hwcctl --debug version
```

### 检查 API 请求

使用 `--debug` 模式查看实际的 API 请求和响应。

### 验证认证信息

```bash
# 简单测试认证是否正确
hwcctl --debug version
```

## 日志分析

### 错误级别

- `[ERROR]` - 需要立即处理的错误
- `[WARN]` - 警告信息
- `[INFO]` - 常规信息
- `[DEBUG]` - 调试信息（需要 --debug 参数）

### 常见日志模式

```bash
# 认证成功
[DEBUG] 认证信息加载成功

# 网络请求
[DEBUG] 发起 CDN API 请求: POST https://cdn.myhuaweicloud.com/...

# 操作完成
[INFO] CDN 缓存刷新成功，任务 ID: task-123456
```

## 获取帮助

### 社区支持

- GitHub Issues: [提交问题](https://github.com/ygqygq2/hwcctl/issues)
- 项目文档: [查看文档](../README.md)

### 华为云支持

- [华为云 CDN 文档](https://support.huaweicloud.com/cdn/)
- [华为云工单系统](https://console.huaweicloud.com/ticket/)

### 提交 Bug 报告

提交问题时请包含：

1. hwcctl 版本：`hwcctl version`
2. 操作系统和架构
3. 完整的错误信息
4. 复现步骤
5. 配置信息（隐藏敏感数据）

**示例 Bug 报告**：

```
**环境信息**
- hwcctl 版本: v0.1.0
- 操作系统: Ubuntu 20.04
- 架构: amd64

**问题描述**
CDN 刷新操作超时

**错误信息**
```

Error: context deadline exceeded

```

**复现步骤**
1. 配置认证信息
2. 执行 `hwcctl cdn refresh --urls "https://example.com/test.jpg"`
3. 等待约 30 秒后出现超时错误

**配置信息**
- Region: cn-north-4
- 网络环境: 企业内网
```
