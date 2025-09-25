# hwcctl

[![License](https://img.shields.io/github/license/ygqygq2/hwcctl)](https://github.com/ygqygq2/hwcctl/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/ygqygq2/hwcctl)](https://github.com/ygqygq2/hwcctl/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygqygq2/hwcctl)](https://goreportcard.com/report/github.com/ygqygq2/hwcctl)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ygqygq2/hwcctl)](https://golang.org/)

> 华为云命令行工具 - 专为华为云服务设计的现代 CLI 工具

hwcctl 是一个强大、易用的华为云命令行工具，类似于 AWS CLI，但专门针对华为云服务设计。它提供了简洁的命令行界面来管理华为云资源，支持多种输出格式，具备完善的错误处理和重试机制。

## ✨ 主要特性

### 🔧 已实现功能

- **CDN 管理**: 内容分发网络缓存刷新、内容预热、任务查询
- **配置管理**: 交互式配置、多种认证方式
- **输出格式**: 支持 table、json、yaml 格式
- **错误处理**: 完善的错误处理和可选重试机制
- **调试支持**: 详细的调试和日志功能

### 🚀 设计特点

- **简洁易用**: 直观的命令行接口，类似 AWS CLI 体验
- **灵活配置**: 支持环境变量、配置文件、命令行参数多种配置方式
- **智能重试**: 可配置的重试机制，默认快速失败避免延迟
- **安全可靠**: 安全的认证信息管理，详细的错误提示

## � 快速开始

### 1. 安装

```bash
# 使用 Go 安装
go install github.com/ygqygq2/hwcctl@latest

# 或从 GitHub Releases 下载
wget https://github.com/ygqygq2/hwcctl/releases/latest/download/hwcctl_linux_amd64.tar.gz
tar -xzf hwcctl_linux_amd64.tar.gz
sudo mv hwcctl /usr/local/bin/
```

### 2. 配置认证

```bash
# 交互式配置（推荐）
hwcctl configure

# 或使用环境变量
export HUAWEICLOUD_ACCESS_KEY="your-access-key"
export HUAWEICLOUD_SECRET_KEY="your-secret-key"
export HUAWEICLOUD_REGION="cn-north-4"
export HUAWEICLOUD_DOMAIN_ID="your-domain-id"  # CDN 必需
```

### 3. 开始使用

```bash
# 查看版本
hwcctl version

# CDN 缓存刷新
hwcctl cdn refresh --urls "https://your-domain.com/file.jpg"

# CDN 内容预热
hwcctl cdn preload --urls "https://your-domain.com/popular-file.mp4"

# 查看帮助
hwcctl --help
hwcctl cdn --help
```

## � 详细文档

- **[📖 完整文档](./docs/README.md)** - 查看所有文档
- **[🚀 快速开始](./docs/02-quick-start.md)** - 详细入门指南
- **[⚙️ 配置指南](./docs/03-configuration.md)** - 认证和配置详解
- **[🔄 CDN 管理](./docs/04-cdn.md)** - CDN 功能使用指南
- **[🔧 故障排查](./docs/08-troubleshooting.md)** - 问题解决方案

## 🛠️ 开发

```bash
# 克隆项目
git clone https://github.com/ygqygq2/hwcctl.git
cd hwcctl

# 安装依赖
go mod tidy

# 本地构建
go build -o hwcctl .

# 使用 task 构建（推荐）
task build
```

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](./docs/11-development.md) 了解如何参与项目开发。

## 📄 许可证

本项目使用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢 [华为云](https://www.huaweicloud.com/) 提供云服务平台，感谢 [Cobra](https://github.com/spf13/cobra) 提供 CLI 框架支持。
