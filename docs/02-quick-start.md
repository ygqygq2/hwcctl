# 02. 快速开始

本指南将帮助您快速开始使用 hwcctl。

## 前置条件

1. ✅ 已安装 hwcctl ([安装指南](./01-installation.md))
2. ✅ 拥有华为云账号和访问密钥
3. ✅ 已获取 Domain ID（CDN 服务需要）

## 第一步：配置认证

### 方式一：交互式配置（推荐新手）

```bash
hwcctl configure
```

按提示输入：

- Access Key ID
- Secret Access Key
- Default region (如 `cn-north-4`)
- Output format (如 `table`)

### 方式二：环境变量配置

```bash
export HUAWEICLOUD_ACCESS_KEY="your-access-key"
export HUAWEICLOUD_SECRET_KEY="your-secret-key"
export HUAWEICLOUD_REGION="cn-north-4"
export HUAWEICLOUD_DOMAIN_ID="your-domain-id"
```

## 第二步：测试连接

```bash
# 查看版本信息
hwcctl version
```

## 第三步：使用 CDN 功能

### 刷新 CDN 缓存

```bash
# 刷新单个文件
hwcctl cdn refresh --urls "https://example.com/image.jpg"

# 刷新多个文件
hwcctl cdn refresh --urls "https://example.com/file1.jpg,https://example.com/file2.css"

# 刷新目录
hwcctl cdn refresh --type directory --urls "https://example.com/images/"
```

### 预热 CDN 缓存

```bash
# 预热单个文件
hwcctl cdn preload --urls "https://example.com/popular-video.mp4"

# 预热多个文件
hwcctl cdn preload --urls "https://example.com/file1.jpg,https://example.com/file2.css"
```

### 查询任务状态

```bash
# 查询任务状态（替换为实际的任务ID）
hwcctl cdn task task-123456789
```

## 第四步：配置输出格式

```bash
# 表格格式（默认）
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output table

# JSON 格式
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output json

# YAML 格式
hwcctl cdn refresh --urls "https://example.com/test.jpg" --output yaml
```

## 常用选项

- `--debug` - 启用调试模式，显示详细日志
- `--verbose` - 详细输出
- `--output` - 输出格式 (table/json/yaml)
- `--region` - 指定区域
- `--help` - 查看帮助信息

## 获取帮助

```bash
# 查看所有命令
hwcctl --help

# 查看 CDN 命令帮助
hwcctl cdn --help

# 查看特定子命令帮助
hwcctl cdn refresh --help
```

## 下一步

- 了解详细的 [CDN 管理](./04-cdn.md)
- 查看 [配置指南](./03-configuration.md) 进行高级配置
- 遇到问题请查看 [故障排查](./08-troubleshooting.md)
