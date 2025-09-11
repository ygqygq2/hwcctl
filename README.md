# hwcctl - 华为云命令行工具

[![Test](https://github.com/ygqygq2/hwcctl/actions/workflows/test.yml/badge.svg)](https://github.com/ygqygq2/hwcctl/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ygqygq2/hwcctl)](https://goreportcard.com/report/github.com/ygqygq2/hwcctl)
[![License](https://img.shields.io/github/license/ygqygq2/hwcctl)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ygqygq2/hwcctl)](https://github.com/ygqygq2/hwcctl/releases)

一个高性能、功能完整的华为云命令行工具，用 Go 语言编写，类似于 AWS CLI，专门用于调用华为云 API 进行各种运维操作。

## ✨ 核心特性

### 🔧 多服务支持

- **ECS**: 弹性云服务器管理（创建、删除、启动、停止等）
- **VPC**: 虚拟私有云管理（VPC、子网、安全组等）
- **更多服务**: 持续添加中...

### 🚀 高性能设计

- **命令行参数**: 基于参数而非配置文件，使用更灵活
- **并发处理**: 支持并发操作，提升执行效率
- **认证管理**: 支持多种认证方式（环境变量、命令行参数）
- **输出格式**: 支持多种输出格式（table、json、yaml）

### ⚙️ 灵活配置

- **多区域支持**: 支持华为云所有区域
- **认证方式**: 支持 Access Key/Secret Key 认证
- **调试模式**: 详细的调试和日志输出
- **自定义输出**: 可选择不同的输出格式

### 🔒 可靠性保证

- **错误处理**: 完善的错误处理和重试机制
- **详细日志**: 分级日志输出，便于问题排查
- **参数验证**: 严格的参数验证和提示
- **安全认证**: 安全的认证信息管理

## 📦 快速开始

### 下载安装

从 [Releases](https://github.com/ygqygq2/hwcctl/releases) 页面下载适合你系统的预编译二进制文件：

```bash
# Linux
wget https://github.com/ygqygq2/hwcctl/releases/latest/download/hwcctl_Linux_x86_64.zip
unzip hwcctl_Linux_x86_64.zip

# Windows
# 下载 hwcctl_Windows_x86_64.zip 并解压

# macOS
wget https://github.com/ygqygq2/hwcctl/releases/latest/download/hwcctl_Darwin_x86_64.zip
unzip hwcctl_Darwin_x86_64.zip
```

### 配置认证

支持两种认证方式：

**方式一：环境变量（推荐）**

```bash
export HUAWEICLOUD_ACCESS_KEY="your-access-key"
export HUAWEICLOUD_SECRET_KEY="your-secret-key"
export HUAWEICLOUD_REGION="cn-north-1"
```

**方式二：命令行参数**

```bash
hwcctl --access-key your-access-key --secret-key your-secret-key --region cn-north-1 <command>
```

### 基础使用

```bash
# 查看帮助信息
hwcctl --help

# 查看版本信息
hwcctl version

# 列出 ECS 实例
hwcctl ecs list

# 创建 ECS 实例
hwcctl ecs create --name my-server --image-id ubuntu-20.04 --flavor-id s3.large.2

# 列出 VPC
hwcctl vpc list

# 创建 VPC
hwcctl vpc create --name my-vpc --cidr 192.168.0.0/16

# 启用调试模式
hwcctl --debug ecs list

# 指定输出格式
hwcctl --output json ecs list
hwcctl --output yaml vpc list
```

## 🛠️ 开发和构建

### 本地开发

```bash
# 克隆代码
git clone https://github.com/ygqygq2/hwcctl.git
cd hwcctl

# 安装依赖 (需要先安装 Task)
task deps

# 运行测试
task test

# 本地构建
task build
```

### 自动化任务

项目使用 [Task](https://taskfile.dev/) 进行自动化：

```bash
task test          # 运行测试
task test-coverage # 测试覆盖率
task build         # 构建二进制
task release       # 发布构建 (多平台)
task clean         # 清理构建产物

# Git hooks 管理
task install-hooks # 安装 pre-commit hooks
task test-hooks    # 测试 hooks 状态
```

### Git Hooks 自动化

项目配置了 pre-commit hooks，每次提交时自动：

- 🔧 格式化 Go 代码
- 🔍 运行静态分析
- 📦 检查依赖状态

无需手动记住运行 `task fmt`！详见 [Git Hooks 说明](docs/GIT_HOOKS.md)。

## 📋 支持的服务

### 弹性云服务器 (ECS)

| 命令                | 说明              | 示例                                                         |
| ------------------- | ----------------- | ------------------------------------------------------------ |
| `hwcctl ecs list`   | 列出所有 ECS 实例 | `hwcctl ecs list`                                            |
| `hwcctl ecs create` | 创建 ECS 实例     | `hwcctl ecs create --name my-server --image-id ubuntu-20.04` |
| `hwcctl ecs delete` | 删除 ECS 实例     | `hwcctl ecs delete instance-id`                              |

### 虚拟私有云 (VPC)

| 命令                | 说明         | 示例                                                    |
| ------------------- | ------------ | ------------------------------------------------------- |
| `hwcctl vpc list`   | 列出所有 VPC | `hwcctl vpc list`                                       |
| `hwcctl vpc create` | 创建 VPC     | `hwcctl vpc create --name my-vpc --cidr 192.168.0.0/16` |
| `hwcctl vpc delete` | 删除 VPC     | `hwcctl vpc delete vpc-id`                              |

## 🌍 支持的区域

- `cn-north-1` - 华北-北京一
- `cn-north-4` - 华北-北京四
- `cn-east-2` - 华东-上海二
- `cn-east-3` - 华东-上海一
- `cn-south-1` - 华南-广州
- `cn-southwest-2` - 西南-贵阳一
- 更多区域持续支持中...

## 📈 性能特点

- **轻量级**: 单一二进制文件，无依赖
- **快速响应**: 高效的 API 调用和数据处理
- **并发支持**: 支持并发操作，提升批量操作效率
- **内存优化**: 优化的内存使用，适合大规模操作

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！请确保：

1. 代码通过所有测试：`task test`
2. 代码格式符合规范：`task fmt-check`
3. 测试覆盖率不低于 50%：`task test-coverage`

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。
