# 01. 安装指南

## 系统要求

- **操作系统**: Linux, macOS, Windows
- **架构**: x86_64, ARM64
- **内存**: 最小 64MB

## 安装方式

### 1. 从 GitHub Releases 下载（推荐）

```bash
# 下载最新版本
wget https://github.com/ygqygq2/hwcctl/releases/latest/download/hwcctl_linux_amd64.tar.gz

# 解压
tar -xzf hwcctl_linux_amd64.tar.gz

# 移动到系统路径
sudo mv hwcctl /usr/local/bin/

# 验证安装
hwcctl version
```

### 2. 使用 Go 安装

```bash
go install github.com/ygqygq2/hwcctl@latest
```

### 3. 从源码编译

```bash
# 克隆仓库
git clone https://github.com/ygqygq2/hwcctl.git
cd hwcctl

# 编译
go build -o hwcctl .

# 安装到系统路径
sudo mv hwcctl /usr/local/bin/
```

## 验证安装

```bash
# 查看版本信息
hwcctl version

# 查看帮助信息
hwcctl --help
```

## 下一步

安装完成后，请查看 [配置指南](./03-configuration.md) 进行初始配置。
