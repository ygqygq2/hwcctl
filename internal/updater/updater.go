package updater

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ygqygq2/hwcctl/internal/logx"
)

// Config 更新器配置
type Config struct {
	Owner      string // GitHub 仓库所有者
	Repo       string // GitHub 仓库名称
	CurrentVer string // 当前版本
	OS         string // 操作系统
	Arch       string // 架构
	Verbose    bool   // 详细输出
	Debug      bool   // 调试模式
}

// Updater 更新器
type Updater struct {
	config     *Config
	httpClient *http.Client
}

// GitHubRelease GitHub Release API 响应结构
type GitHubRelease struct {
	TagName    string        `json:"tag_name"`
	Name       string        `json:"name"`
	Body       string        `json:"body"`
	Draft      bool          `json:"draft"`
	Prerelease bool          `json:"prerelease"`
	Assets     []GitHubAsset `json:"assets"`
}

// GitHubAsset GitHub Release Asset
type GitHubAsset struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"browser_download_url"`
}

// New 创建新的更新器
func New(config *Config) *Updater {
	return &Updater{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckUpdate 检查是否有可用更新
func (u *Updater) CheckUpdate() error {
	logx.Infof("检查 %s/%s 的最新版本...", u.config.Owner, u.config.Repo)

	latest, err := u.getLatestRelease()
	if err != nil {
		return fmt.Errorf("获取最新版本信息失败: %w", err)
	}

	currentVer := strings.TrimPrefix(u.config.CurrentVer, "v")
	latestVer := strings.TrimPrefix(latest.TagName, "v")

	fmt.Printf("当前版本: %s\n", currentVer)
	fmt.Printf("最新版本: %s\n", latestVer)

	if currentVer == latestVer {
		fmt.Println("✅ 已经是最新版本！")
		return nil
	}

	if currentVer == "dev" {
		fmt.Println("⚠️  当前运行的是开发版本")
	}

	fmt.Printf("🎉 发现新版本: %s\n", latest.TagName)
	if latest.Body != "" {
		fmt.Printf("\n更新说明:\n%s\n", latest.Body)
	}

	return nil
}

// Update 执行更新
func (u *Updater) Update(force bool, targetVersion string) error {
	var release *GitHubRelease
	var err error

	if targetVersion != "" {
		release, err = u.getReleaseByTag(targetVersion)
		if err != nil {
			return fmt.Errorf("获取指定版本 %s 失败: %w", targetVersion, err)
		}
	} else {
		release, err = u.getLatestRelease()
		if err != nil {
			return fmt.Errorf("获取最新版本信息失败: %w", err)
		}
	}

	// 检查是否需要更新
	currentVer := strings.TrimPrefix(u.config.CurrentVer, "v")
	targetVer := strings.TrimPrefix(release.TagName, "v")

	if !force && currentVer == targetVer && currentVer != "dev" {
		fmt.Println("✅ 已经是最新版本！")
		return nil
	}

	fmt.Printf("准备更新到版本: %s\n", release.TagName)

	// 查找合适的资源
	asset, err := u.findAsset(release)
	if err != nil {
		return fmt.Errorf("未找到适合的安装包: %w", err)
	}

	fmt.Printf("找到安装包: %s (%.2f MB)\n", asset.Name, float64(asset.Size)/1024/1024)

	// 下载文件
	tempFile, err := u.downloadAsset(asset)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer os.Remove(tempFile)

	// 提取并安装
	if err := u.installUpdate(tempFile, asset.Name); err != nil {
		return fmt.Errorf("安装更新失败: %w", err)
	}

	fmt.Printf("✅ 成功更新到版本 %s\n", release.TagName)
	return nil
}

// getLatestRelease 获取最新 release
func (u *Updater) getLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.config.Owner, u.config.Repo)

	resp, err := u.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回错误: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &release, nil
}

// getReleaseByTag 根据标签获取 release
func (u *Updater) getReleaseByTag(tag string) (*GitHubRelease, error) {
	// 确保标签以 v 开头
	if !strings.HasPrefix(tag, "v") {
		tag = "v" + tag
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", u.config.Owner, u.config.Repo, tag)

	resp, err := u.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("版本 %s 不存在", tag)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API 返回错误: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &release, nil
}

// findAsset 查找适合当前平台的资源
func (u *Updater) findAsset(release *GitHubRelease) (*GitHubAsset, error) {
	osName := u.config.OS
	archName := u.config.Arch

	// 标准化操作系统名称
	switch osName {
	case "windows":
		osName = "Windows"
	case "darwin":
		osName = "Darwin"
	case "linux":
		osName = "Linux"
	}

	// 标准化架构名称
	switch archName {
	case "amd64":
		archName = "x86_64"
	case "386":
		archName = "i386"
	case "arm64":
		archName = "arm64"
	}

	// 查找匹配的资源
	for _, asset := range release.Assets {
		name := asset.Name

		// 检查是否包含操作系统和架构信息
		if strings.Contains(name, osName) && strings.Contains(name, archName) {
			// 优先选择 .zip 文件
			if strings.HasSuffix(name, ".zip") {
				return &asset, nil
			}
		}
	}

	// 如果没找到 .zip，再找其他格式
	for _, asset := range release.Assets {
		name := asset.Name
		if strings.Contains(name, osName) && strings.Contains(name, archName) {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("未找到适合 %s/%s 的安装包", osName, archName)
}

// downloadAsset 下载资源文件
func (u *Updater) downloadAsset(asset *GitHubAsset) (string, error) {
	logx.Infof("开始下载: %s", asset.DownloadURL)

	resp, err := u.httpClient.Get(asset.DownloadURL)
	if err != nil {
		return "", fmt.Errorf("下载请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败: %s", resp.Status)
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "hwcctl-update-*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer tempFile.Close()

	// 下载并显示进度
	written, err := u.copyWithProgress(tempFile, resp.Body, asset.Size)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("下载失败: %w", err)
	}

	if written != asset.Size {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("下载不完整: 期望 %d 字节，实际 %d 字节", asset.Size, written)
	}

	fmt.Println("\n✅ 下载完成")
	return tempFile.Name(), nil
}

// copyWithProgress 带进度显示的复制
func (u *Updater) copyWithProgress(dst io.Writer, src io.Reader, total int64) (int64, error) {
	var written int64
	buf := make([]byte, 32*1024) // 32KB 缓冲区

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}

			// 显示进度
			if total > 0 {
				progress := float64(written) / float64(total) * 100
				fmt.Printf("\r下载进度: %.1f%% (%.2f/%.2f MB)",
					progress,
					float64(written)/1024/1024,
					float64(total)/1024/1024)
			}
		}
		if er != nil {
			if er != io.EOF {
				return written, er
			}
			break
		}
	}
	return written, nil
}

// installUpdate 安装更新
func (u *Updater) installUpdate(tempFile, assetName string) error {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前可执行文件路径失败: %w", err)
	}

	var binaryData []byte

	// 根据文件类型处理
	if strings.HasSuffix(assetName, ".zip") {
		binaryData, err = u.extractFromZip(tempFile)
		if err != nil {
			return fmt.Errorf("解压文件失败: %w", err)
		}
	} else {
		// 直接读取二进制文件
		binaryData, err = os.ReadFile(tempFile)
		if err != nil {
			return fmt.Errorf("读取下载文件失败: %w", err)
		}
	}

	// 验证文件是否为有效的可执行文件（简单检查）
	if len(binaryData) < 100 {
		return fmt.Errorf("下载的文件太小，可能不是有效的可执行文件")
	}

	// 使用延迟替换策略来解决 "text file busy" 问题
	return u.performDelayedReplacement(execPath, binaryData)
}

// performDelayedReplacement 执行延迟替换
func (u *Updater) performDelayedReplacement(execPath string, binaryData []byte) error {
	// 创建新的可执行文件路径
	newExecPath := execPath + ".new"

	// 写入新的可执行文件
	if err := os.WriteFile(newExecPath, binaryData, 0755); err != nil {
		return fmt.Errorf("写入新文件失败: %w", err)
	}

	// 备份当前文件
	backupPath := execPath + ".backup"
	if err := u.copyFile(execPath, backupPath); err != nil {
		os.Remove(newExecPath)
		return fmt.Errorf("备份当前文件失败: %w", err)
	}

	logx.Infof("已备份当前文件到: %s", backupPath)

	// 根据操作系统选择不同的替换策略
	if runtime.GOOS == "windows" {
		return u.performWindowsReplacement(execPath, newExecPath, backupPath)
	} else {
		return u.performUnixReplacement(execPath, newExecPath, backupPath)
	}
}

// performUnixReplacement Unix系统的文件替换
func (u *Updater) performUnixReplacement(execPath, newExecPath, backupPath string) error {
	// 在 Unix 系统上，我们可以创建一个脚本来延迟替换
	scriptContent := fmt.Sprintf(`#!/bin/bash
# 等待当前进程退出
sleep 2

# 替换可执行文件
if mv "%s" "%s"; then
    echo "✅ 更新完成！"
    # 删除备份文件
    rm -f "%s"
    # 删除脚本自身
    rm -f "$0"
else
    echo "❌ 更新失败，正在恢复备份..."
    # 恢复备份
    mv "%s" "%s" || true
    # 删除新文件
    rm -f "%s"
    # 删除脚本自身
    rm -f "$0"
    exit 1
fi
`, newExecPath, execPath, backupPath, backupPath, execPath, newExecPath)

	scriptPath := filepath.Dir(execPath) + "/.hwcctl_update.sh"
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		os.Remove(newExecPath)
		os.Remove(backupPath)
		return fmt.Errorf("创建更新脚本失败: %w", err)
	}

	// 启动后台脚本
	cmd := exec.Command("/bin/bash", scriptPath)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		os.Remove(newExecPath)
		os.Remove(backupPath)
		return fmt.Errorf("启动更新脚本失败: %w", err)
	}

	// 提示用户
	fmt.Println("🚀 更新将在程序退出后完成...")
	fmt.Println("💡 请稍等片刻，然后重新运行程序以验证更新。")

	return nil
}

// performWindowsReplacement Windows系统的文件替换
func (u *Updater) performWindowsReplacement(execPath, newExecPath, backupPath string) error {
	// Windows 批处理脚本
	scriptContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
move "%s" "%s" >nul 2>&1
if %%errorlevel%% equ 0 (
    echo 更新完成！
    del "%s" >nul 2>&1
) else (
    echo 更新失败，正在恢复备份...
    move "%s" "%s" >nul 2>&1
    del "%s" >nul 2>&1
)
del "%%~f0"
`, newExecPath, execPath, backupPath, backupPath, execPath, newExecPath)

	scriptPath := filepath.Dir(execPath) + "\\hwcctl_update.bat"
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		os.Remove(newExecPath)
		os.Remove(backupPath)
		return fmt.Errorf("创建更新脚本失败: %w", err)
	}

	// 启动批处理脚本
	cmd := exec.Command("cmd", "/C", "start", "/B", scriptPath)
	if err := cmd.Start(); err != nil {
		os.Remove(scriptPath)
		os.Remove(newExecPath)
		os.Remove(backupPath)
		return fmt.Errorf("启动更新脚本失败: %w", err)
	}

	// 提示用户
	fmt.Println("🚀 更新将在程序退出后完成...")
	fmt.Println("💡 请稍等片刻，然后重新运行程序以验证更新。")

	return nil
}

// extractFromZip 从 ZIP 文件中提取二进制文件
func (u *Updater) extractFromZip(zipPath string) ([]byte, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// 查找可执行文件
	for _, f := range r.File {
		name := strings.ToLower(f.Name)

		// 跳过目录和其他文件
		if f.FileInfo().IsDir() {
			continue
		}

		// 查找 hwcctl 可执行文件
		if strings.Contains(name, "hwcctl") && !strings.Contains(name, ".") {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("在 ZIP 文件中未找到 hwcctl 可执行文件")
}

// copyFile 复制文件
func (u *Updater) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
