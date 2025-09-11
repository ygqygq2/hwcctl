# hwcctl 项目架构

## 概述

hwcctl 是一个现代化的华为云命令行工具，采用 Go 语言开发，基于 Cobra 框架构建。本文档描述了项目的整体架构和设计理念。

## 项目结构

```
hwcctl/
├── cmd/                    # 命令行接口层
│   ├── root.go            # 根命令和全局配置
│   ├── ecs.go             # ECS 服务命令
│   ├── ecs_operations.go  # ECS 具体操作
│   ├── vpc.go             # VPC 服务命令
│   └── vpc_operations.go  # VPC 具体操作
├── internal/              # 内部包（不对外暴露）
│   ├── auth/              # 认证管理
│   │   └── auth.go        # 华为云认证逻辑
│   ├── logx/              # 日志系统
│   │   └── logx.go        # 分级日志实现
│   └── utils/             # 工具函数
│       └── output.go      # 输出格式处理
├── docs/                  # 文档
│   ├── ARCHITECTURE.md    # 架构文档
│   └── GIT_HOOKS.md       # Git Hooks 说明
├── .github/               # GitHub 配置
│   └── workflows/         # CI/CD 工作流
├── main.go               # 程序入口
├── go.mod               # Go 模块定义
├── Taskfile.yml         # 任务自动化
└── .goreleaser.yaml     # 发布配置
```

## 架构设计

### 1. 分层架构

```
┌─────────────────────────┐
│      CLI Interface      │  <- cmd/ 包
├─────────────────────────┤
│    Business Logic       │  <- internal/ 包
├─────────────────────────┤
│   Huawei Cloud APIs     │  <- 华为云 SDK
└─────────────────────────┘
```

### 2. 核心组件

#### 命令层 (cmd/)

- **职责**: 处理用户输入，参数解析，命令路由
- **技术栈**: Cobra CLI 框架
- **特点**:
  - 支持子命令和嵌套命令
  - 自动生成帮助信息
  - 参数验证和类型转换

#### 业务逻辑层 (internal/)

- **auth/**: 华为云认证管理

  - 支持多种认证方式（环境变量、参数）
  - 认证信息安全处理
  - 区域配置管理

- **logx/**: 统一日志系统

  - 分级日志输出 (DEBUG/INFO/WARN/ERROR)
  - 可配置的日志级别
  - 结构化日志支持

- **utils/**: 通用工具函数
  - 多格式输出 (Table/JSON/YAML)
  - 字符串处理工具
  - 数据格式转换

### 3. 设计原则

#### 3.1 单一职责原则

每个包和模块都有明确的职责边界：

- `cmd/` 只负责 CLI 交互
- `internal/auth/` 只负责认证
- `internal/logx/` 只负责日志
- `internal/utils/` 提供通用工具

#### 3.2 依赖倒置原则

- 高层模块不依赖低层模块
- 通过接口定义契约
- 便于单元测试和模块替换

#### 3.3 开闭原则

- 对扩展开放：易于添加新的华为云服务
- 对修改封闭：现有功能稳定，不轻易修改

## 扩展指南

### 添加新服务

1. **创建服务命令文件**

   ```go
   // cmd/rds.go
   var rdsCmd = &cobra.Command{
       Use:   "rds",
       Short: "云数据库 RDS 相关操作",
   }
   ```

2. **创建操作子命令**

   ```go
   // cmd/rds_operations.go
   var rdsListCmd = &cobra.Command{
       Use:   "list",
       Short: "列出 RDS 实例",
       RunE:  listRDSInstances,
   }
   ```

3. **注册命令**
   ```go
   func init() {
       rootCmd.AddCommand(rdsCmd)
       rdsCmd.AddCommand(rdsListCmd)
   }
   ```

### 添加新的输出格式

在 `internal/utils/output.go` 中扩展 `FormatOutput` 函数。

### 添加新的认证方式

在 `internal/auth/auth.go` 中扩展认证逻辑。

## 配置管理

### 全局配置

- 区域选择 (`--region`)
- 认证信息 (`--access-key`, `--secret-key`)
- 输出格式 (`--output`)
- 调试模式 (`--debug`, `--verbose`)

### 环境变量支持

- `HUAWEICLOUD_ACCESS_KEY`
- `HUAWEICLOUD_SECRET_KEY`
- `HUAWEICLOUD_REGION`

## 错误处理

### 分级错误处理

1. **参数验证错误**: 在命令层捕获，提供用户友好的错误信息
2. **API 调用错误**: 在业务逻辑层处理，支持重试机制
3. **系统错误**: 记录详细日志，向用户提供简化的错误信息

### 错误信息本地化

所有用户面向的错误信息都使用中文，提供更好的用户体验。

## 性能考虑

### 并发处理

- 支持并发 API 调用
- 合理控制并发数量，避免 API 限流

### 内存优化

- 流式处理大数据集
- 及时释放不需要的资源

### 缓存策略

- 认证信息缓存
- API 响应缓存（适用场景）

## 安全考虑

### 认证信息安全

- 不在日志中输出敏感信息
- 支持从环境变量读取认证信息
- 内存中及时清理敏感数据

### 输入验证

- 严格的参数验证
- 防止注入攻击
- 合理的输入长度限制

## 测试策略

### 单元测试

- 每个包都有对应的测试文件
- 覆盖率要求不低于 50%
- Mock 外部依赖

### 集成测试

- 真实环境测试（可选）
- API 集成测试
- 端到端测试

### 性能测试

- 并发性能测试
- 内存使用测试
- 响应时间测试

## 部署和发布

### 构建系统

- 使用 GoReleaser 进行多平台构建
- 自动化版本管理
- 支持快照构建和正式发布

### CI/CD 流程

- 代码质量检查
- 自动化测试
- 自动化发布
- 文档更新

这个架构设计确保了 hwcctl 的可维护性、可扩展性和高性能，同时保持了良好的用户体验。
