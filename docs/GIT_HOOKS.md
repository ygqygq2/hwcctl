# Git Hooks 说明

## 概述

本项目配置了 Git hooks，用于在代码提交前自动执行代码质量检查，确保代码库的质量和一致性。

## 已配置的 Hooks

### pre-commit

在每次 `git commit` 之前自动执行以下检查：

1. **代码格式化** (`task fmt`)

   - 自动格式化 Go 代码
   - 确保代码风格一致

2. **格式检查** (`task fmt-check`)

   - 验证代码格式是否正确
   - 如有未格式化的文件会阻止提交

3. **静态分析** (`task vet`)

   - 运行 `go vet` 静态分析
   - 检查潜在的代码问题

4. **依赖检查** (`task deps-check`)
   - 检查 `go.mod` 和 `go.sum` 是否是最新状态
   - 确保依赖管理正确

## 管理 Hooks

### 安装 Hooks

```bash
task install-hooks
```

这个命令会：

- 检查 pre-commit hook 文件是否存在
- 为 hook 文件添加执行权限
- 显示安装成功的确认信息

### 卸载 Hooks

```bash
task uninstall-hooks
```

### 测试 Hooks

```bash
task test-hooks
```

这个命令会验证 hooks 是否正确安装和可执行。

## Hooks 工作流程

1. **执行 `git commit`**
2. **触发 pre-commit hook**
3. **依次执行检查项**：
   - 如果所有检查通过 → 允许提交
   - 如果任何检查失败 → 阻止提交并显示错误信息

## 绕过 Hooks（不推荐）

在紧急情况下，可以使用以下命令绕过 hooks：

```bash
git commit --no-verify -m "紧急提交消息"
```

**注意**：不建议经常绕过 hooks，这会降低代码质量。

## 故障排除

### Hook 不执行

1. 检查 hook 文件是否存在：

   ```bash
   ls -la .git/hooks/pre-commit
   ```

2. 检查文件权限：

   ```bash
   chmod +x .git/hooks/pre-commit
   ```

3. 重新安装 hooks：
   ```bash
   task install-hooks
   ```

### Hook 执行失败

1. 手动运行相关任务来识别问题：

   ```bash
   task fmt
   task fmt-check
   task vet
   task deps-check
   ```

2. 修复报告的问题后重新尝试提交

## 最佳实践

1. **定期运行开发测试**：

   ```bash
   task dev-test
   ```

2. **在提交前手动运行检查**：

   ```bash
   task ci-quality
   ```

3. **保持依赖更新**：

   ```bash
   task tidy
   ```

4. **使用开发环境设置**：
   ```bash
   task dev-setup
   ```

这样可以确保在提交前就发现并解决问题，避免 hook 执行失败。
