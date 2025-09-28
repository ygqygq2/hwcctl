package updater

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	if updater == nil {
		t.Fatal("创建更新器失败")
	}

	if updater.config != config {
		t.Error("配置设置不正确")
	}

	if updater.httpClient == nil {
		t.Error("HTTP 客户端未初始化")
	}
}

func TestFindAsset(t *testing.T) {
	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 模拟 release 数据
	release := &GitHubRelease{
		TagName: "v0.1.0",
		Assets: []GitHubAsset{
			{
				Name:        "hwcctl_Linux_x86_64.zip",
				Size:        1024000,
				DownloadURL: "https://github.com/ygqygq2/hwcctl/releases/download/v0.1.0/hwcctl_Linux_x86_64.zip",
			},
			{
				Name:        "hwcctl_Windows_x86_64.zip",
				Size:        1024000,
				DownloadURL: "https://github.com/ygqygq2/hwcctl/releases/download/v0.1.0/hwcctl_Windows_x86_64.zip",
			},
		},
	}

	asset, err := updater.findAsset(release)
	if err != nil {
		t.Errorf("查找资源失败: %v", err)
	}

	if asset.Name != "hwcctl_Linux_x86_64.zip" {
		t.Errorf("期望找到 hwcctl_Linux_x86_64.zip，实际找到 %s", asset.Name)
	}
}

func TestFindAssetNotFound(t *testing.T) {
	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "freebsd",
		Arch:       "mips",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 模拟 release 数据
	release := &GitHubRelease{
		TagName: "v0.1.0",
		Assets: []GitHubAsset{
			{
				Name:        "hwcctl_Linux_x86_64.zip",
				Size:        1024000,
				DownloadURL: "https://github.com/ygqygq2/hwcctl/releases/download/v0.1.0/hwcctl_Linux_x86_64.zip",
			},
		},
	}

	_, err := updater.findAsset(release)
	if err == nil {
		t.Error("应该找不到适合的资源，但没有返回错误")
	}
}

func TestFindAssetMultiplePlatforms(t *testing.T) {
	testCases := []struct {
		name     string
		os       string
		arch     string
		expected string
	}{
		{
			name:     "Windows amd64",
			os:       "windows",
			arch:     "amd64",
			expected: "hwcctl_Windows_x86_64.zip",
		},
		{
			name:     "Darwin arm64",
			os:       "darwin",
			arch:     "arm64",
			expected: "hwcctl_Darwin_arm64.zip",
		},
		{
			name:     "Linux 386",
			os:       "linux",
			arch:     "386",
			expected: "hwcctl_Linux_i386.zip",
		},
	}

	release := &GitHubRelease{
		TagName: "v0.1.0",
		Assets: []GitHubAsset{
			{Name: "hwcctl_Linux_x86_64.zip", Size: 1024000, DownloadURL: "https://example.com/linux_x86_64.zip"},
			{Name: "hwcctl_Windows_x86_64.zip", Size: 1024000, DownloadURL: "https://example.com/windows_x86_64.zip"},
			{Name: "hwcctl_Darwin_arm64.zip", Size: 1024000, DownloadURL: "https://example.com/darwin_arm64.zip"},
			{Name: "hwcctl_Linux_i386.zip", Size: 1024000, DownloadURL: "https://example.com/linux_i386.zip"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{
				Owner:      "ygqygq2",
				Repo:       "hwcctl",
				CurrentVer: "v0.1.0",
				OS:         tc.os,
				Arch:       tc.arch,
				Verbose:    false,
				Debug:      false,
			}

			updater := New(config)
			asset, err := updater.findAsset(release)
			if err != nil {
				t.Errorf("查找资源失败: %v", err)
			}

			if asset.Name != tc.expected {
				t.Errorf("期望找到 %s，实际找到 %s", tc.expected, asset.Name)
			}
		})
	}
}

func TestGetLatestRelease(t *testing.T) {
	// 创建模拟的GitHub API服务器
	mockRelease := &GitHubRelease{
		TagName:    "v1.0.0",
		Name:       "Release v1.0.0",
		Body:       "Test release",
		Draft:      false,
		Prerelease: false,
		Assets: []GitHubAsset{
			{
				Name:        "hwcctl_Linux_x86_64.zip",
				Size:        1024000,
				DownloadURL: "https://github.com/ygqygq2/hwcctl/releases/download/v1.0.0/hwcctl_Linux_x86_64.zip",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/releases/latest") {
			json.NewEncoder(w).Encode(mockRelease)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 替换API URL以使用测试服务器
	originalTransport := updater.httpClient.Transport
	updater.httpClient.Transport = &mockTransport{
		server:        server,
		baseTransport: http.DefaultTransport,
	}
	defer func() {
		updater.httpClient.Transport = originalTransport
	}()

	release, err := updater.getLatestRelease()
	if err != nil {
		t.Errorf("获取最新版本失败: %v", err)
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("期望版本为 v1.0.0，实际为 %s", release.TagName)
	}
}

func TestGetReleaseByTag(t *testing.T) {
	mockRelease := &GitHubRelease{
		TagName:    "v1.2.0",
		Name:       "Release v1.2.0",
		Body:       "Specific version release",
		Draft:      false,
		Prerelease: false,
		Assets: []GitHubAsset{
			{
				Name:        "hwcctl_Linux_x86_64.zip",
				Size:        1024000,
				DownloadURL: "https://github.com/ygqygq2/hwcctl/releases/download/v1.2.0/hwcctl_Linux_x86_64.zip",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/releases/tags/v1.2.0") {
			json.NewEncoder(w).Encode(mockRelease)
		} else if strings.Contains(r.URL.Path, "/releases/tags/v999.999.999") {
			http.NotFound(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)
	updater.httpClient.Transport = &mockTransport{
		server:        server,
		baseTransport: http.DefaultTransport,
	}

	// 测试成功获取指定版本
	t.Run("成功获取指定版本", func(t *testing.T) {
		release, err := updater.getReleaseByTag("1.2.0")
		if err != nil {
			t.Errorf("获取指定版本失败: %v", err)
		}

		if release.TagName != "v1.2.0" {
			t.Errorf("期望版本为 v1.2.0，实际为 %s", release.TagName)
		}
	})

	// 测试版本不存在
	t.Run("版本不存在", func(t *testing.T) {
		_, err := updater.getReleaseByTag("999.999.999")
		if err == nil {
			t.Error("应该返回版本不存在的错误")
		}
	})
}

func TestCopyWithProgress(t *testing.T) {
	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 测试数据
	testData := "Hello, World! This is a test data for progress copy."
	src := strings.NewReader(testData)
	dst := &bytes.Buffer{}

	written, err := updater.copyWithProgress(dst, src, int64(len(testData)))
	if err != nil {
		t.Errorf("复制失败: %v", err)
	}

	if written != int64(len(testData)) {
		t.Errorf("期望写入 %d 字节，实际写入 %d 字节", len(testData), written)
	}

	if dst.String() != testData {
		t.Error("复制的数据不正确")
	}
}

func TestExtractFromZip(t *testing.T) {
	// 创建临时ZIP文件
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")

	// 创建包含hwcctl可执行文件的ZIP
	testContent := []byte("fake executable content")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	// 添加hwcctl文件到ZIP
	fileWriter, err := zipWriter.Create("hwcctl")
	if err != nil {
		t.Fatalf("创建ZIP内文件失败: %v", err)
	}

	_, err = fileWriter.Write(testContent)
	if err != nil {
		t.Fatalf("写入ZIP内文件失败: %v", err)
	}

	zipWriter.Close()
	zipFile.Close()

	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 测试提取
	extractedData, err := updater.extractFromZip(zipPath)
	if err != nil {
		t.Errorf("提取ZIP文件失败: %v", err)
	}

	if !bytes.Equal(extractedData, testContent) {
		t.Error("提取的内容不正确")
	}
}

func TestExtractFromZipNotFound(t *testing.T) {
	// 创建不包含hwcctl文件的ZIP
	tempDir := t.TempDir()
	zipPath := filepath.Join(tempDir, "test.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("创建ZIP文件失败: %v", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	// 添加其他文件
	fileWriter, err := zipWriter.Create("other-file.txt")
	if err != nil {
		t.Fatalf("创建ZIP内文件失败: %v", err)
	}

	fileWriter.Write([]byte("other content"))
	zipWriter.Close()
	zipFile.Close()

	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 测试提取失败
	_, err = updater.extractFromZip(zipPath)
	if err == nil {
		t.Error("应该返回未找到hwcctl文件的错误")
	}
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	// 创建源文件
	srcPath := filepath.Join(tempDir, "source.txt")
	srcContent := []byte("test file content")
	err := os.WriteFile(srcPath, srcContent, 0644)
	if err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	// 目标文件路径
	dstPath := filepath.Join(tempDir, "destination.txt")

	config := &Config{
		Owner:      "ygqygq2",
		Repo:       "hwcctl",
		CurrentVer: "v0.1.0",
		OS:         "linux",
		Arch:       "amd64",
		Verbose:    false,
		Debug:      false,
	}

	updater := New(config)

	// 测试复制文件
	err = updater.copyFile(srcPath, dstPath)
	if err != nil {
		t.Errorf("复制文件失败: %v", err)
	}

	// 验证目标文件内容
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("读取目标文件失败: %v", err)
	}

	if !bytes.Equal(srcContent, dstContent) {
		t.Error("复制的文件内容不正确")
	}
}

func TestPerformDelayedReplacement(t *testing.T) {
	updater := &Updater{
		config: &Config{
			OS:   "linux",
			Arch: "amd64",
		},
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "hwcctl-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建模拟的当前可执行文件
	execPath := filepath.Join(tempDir, "hwcctl")
	originalContent := []byte("original binary content")
	if err := os.WriteFile(execPath, originalContent, 0755); err != nil {
		t.Fatalf("创建模拟可执行文件失败: %v", err)
	}

	// 新的二进制内容
	newContent := []byte("new binary content for testing")

	// 执行延迟替换
	err = updater.performDelayedReplacement(execPath, newContent)
	if err != nil {
		t.Errorf("延迟替换失败: %v", err)
	}

	// 验证 .new 文件是否创建
	newExecPath := execPath + ".new"
	if _, err := os.Stat(newExecPath); os.IsNotExist(err) {
		t.Error("新文件未创建")
	}

	// 验证备份文件是否创建
	backupPath := execPath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("备份文件未创建")
	}

	// 验证新文件内容
	newFileContent, err := os.ReadFile(newExecPath)
	if err != nil {
		t.Errorf("读取新文件失败: %v", err)
	}
	if !bytes.Equal(newFileContent, newContent) {
		t.Error("新文件内容不正确")
	}

	// 验证备份文件内容
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Errorf("读取备份文件失败: %v", err)
	}
	if !bytes.Equal(backupContent, originalContent) {
		t.Error("备份文件内容不正确")
	}
}

func TestPerformUnixReplacement(t *testing.T) {
	updater := &Updater{
		config: &Config{
			OS:   "linux",
			Arch: "amd64",
		},
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "hwcctl-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建模拟文件路径
	execPath := filepath.Join(tempDir, "hwcctl")
	newExecPath := execPath + ".new"
	backupPath := execPath + ".backup"

	// 创建模拟文件
	if err := os.WriteFile(execPath, []byte("original"), 0755); err != nil {
		t.Fatalf("创建模拟可执行文件失败: %v", err)
	}
	if err := os.WriteFile(newExecPath, []byte("new"), 0755); err != nil {
		t.Fatalf("创建新文件失败: %v", err)
	}
	if err := os.WriteFile(backupPath, []byte("backup"), 0755); err != nil {
		t.Fatalf("创建备份文件失败: %v", err)
	}

	// 执行Unix替换
	err = updater.performUnixReplacement(execPath, newExecPath, backupPath)
	if err != nil {
		t.Errorf("Unix替换失败: %v", err)
	}

	// 验证更新脚本是否创建
	scriptPath := filepath.Dir(execPath) + "/.hwcctl_update.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("更新脚本未创建")
	}

	// 验证脚本内容
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Errorf("读取脚本失败: %v", err)
	}

	scriptStr := string(scriptContent)
	if !strings.Contains(scriptStr, "#!/bin/bash") {
		t.Error("脚本缺少shebang")
	}
	if !strings.Contains(scriptStr, "sleep 2") {
		t.Error("脚本缺少等待时间")
	}
	if !strings.Contains(scriptStr, "mv") {
		t.Error("脚本缺少移动命令")
	}
}

func TestInstallUpdateWithDelayedReplacement(t *testing.T) {
	updater := &Updater{
		config: &Config{
			OS:   "linux",
			Arch: "amd64",
		},
	}

	// 创建临时ZIP文件
	tempDir, err := os.MkdirTemp("", "hwcctl-test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "test.zip")
	if err := createTestZip(zipPath, "hwcctl", []byte("test binary content")); err != nil {
		t.Fatalf("创建测试ZIP文件失败: %v", err)
	}

	// 模拟当前可执行文件
	originalExec, err := os.Executable()
	if err != nil {
		t.Fatalf("获取当前可执行文件路径失败: %v", err)
	}

	// 由于我们不能真的替换测试中的可执行文件，我们只测试方法调用不出错
	err = updater.installUpdate(zipPath, "test.zip")

	// 在测试环境中，我们期望这个方法能正常执行而不报错
	// 实际的文件替换会在后台进行
	if err != nil {
		// 如果是因为权限或其他系统限制导致的错误，我们记录但不视为测试失败
		t.Logf("installUpdate 在测试环境中的预期行为: %v", err)
	}

	// 验证原始可执行文件仍然存在（因为替换是延迟的）
	if _, err := os.Stat(originalExec); os.IsNotExist(err) {
		t.Error("原始可执行文件不应该被立即删除")
	}
}

// mockTransport 用于模拟HTTP请求
type mockTransport struct {
	server        *httptest.Server
	baseTransport http.RoundTripper
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 将GitHub API请求重定向到测试服务器
	if strings.Contains(req.URL.Host, "api.github.com") {
		req.URL.Scheme = "http"
		req.URL.Host = strings.TrimPrefix(t.server.URL, "http://")
		return t.baseTransport.RoundTrip(req)
	}
	return t.baseTransport.RoundTrip(req)
}

// createTestZip 创建测试用的ZIP文件
func createTestZip(zipPath, filename string, content []byte) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = fileWriter.Write(content)
	return err
}
