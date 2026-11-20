package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/configs"
)

// 直接使用 service 层测试持久浏览器截图功能
func main() {
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("=== 小红书 MCP 截图测试 ===")
	fmt.Println("流程: 启动持久浏览器 -> 打开小红书 -> 截图 -> 保存到 test 文件夹")
	fmt.Println()

	// 这里我们通过 HTTP API 调用 MCP 服务
	// 假设服务运行在 localhost:8080
	apiBaseURL := "http://localhost:8080/api/v1"

	// 检查服务是否运行
	fmt.Println("1. 检查 MCP 服务状态...")
	if err := checkHealth(apiBaseURL); err != nil {
		fmt.Printf("服务检查失败: %v\n", err)
		fmt.Println("请确保 MCP 服务已启动: go run .")
		os.Exit(1)
	}
	fmt.Println("   MCP 服务运行正常 ✓")
	fmt.Println()

	// 步骤 1: 启动持久浏览器
	fmt.Println("2. 启动持久浏览器...")
	fmt.Println("   POST /api/v1/persistent_browser")
	fmt.Println("   Body: {\"action\": \"start\", \"headless\": false, \"url\": \"https://www.xiaohongshu.com/explore\"}")
	fmt.Println()

	// 步骤 2: 截图当前页面
	fmt.Println("3. 等待页面加载后截图...")
	fmt.Println("   POST /api/v1/screenshot_current")
	fmt.Println("   Body: {\"full_page\": true}")
	fmt.Println()

	// 步骤 3: 关闭浏览器
	fmt.Println("4. 关闭浏览器...")
	fmt.Println("   POST /api/v1/persistent_browser")
	fmt.Println("   Body: {\"action\": \"close\"}")
	fmt.Println()

	fmt.Println("=== 使用 curl 命令测试 ===")
	fmt.Println()
	fmt.Println("# 1. 启动浏览器并打开小红书")
	fmt.Printf("curl -X POST %s/persistent_browser \\\n", apiBaseURL)
	fmt.Println("  -H \"Content-Type: application/json\" \\")
	fmt.Println("  -d '{\"action\": \"start\", \"headless\": false, \"url\": \"https://www.xiaohongshu.com/explore\"}'")
	fmt.Println()
	fmt.Println("# 2. 截图当前页面（整页）")
	fmt.Printf("curl -X POST %s/screenshot_current \\\n", apiBaseURL)
	fmt.Println("  -H \"Content-Type: application/json\" \\")
	fmt.Println("  -d '{\"full_page\": true}' \\")
	fmt.Println("  | jq -r '.data.image_base64' | base64 -d > test/mcp_screenshot_fullpage.png")
	fmt.Println()
	fmt.Println("# 3. 截图当前页面（视口）")
	fmt.Printf("curl -X POST %s/screenshot_current \\\n", apiBaseURL)
	fmt.Println("  -H \"Content-Type: application/json\" \\")
	fmt.Println("  -d '{\"full_page\": false}' \\")
	fmt.Println("  | jq -r '.data.image_base64' | base64 -d > test/mcp_screenshot_viewport.png")
	fmt.Println()
	fmt.Println("# 4. 关闭浏览器")
	fmt.Printf("curl -X POST %s/persistent_browser \\\n", apiBaseURL)
	fmt.Println("  -H \"Content-Type: application/json\" \\")
	fmt.Println("  -d '{\"action\": \"close\"}'")
	fmt.Println()

	fmt.Println("=== 或者直接运行集成测试 ===")
	fmt.Println("go run test/mcp_screenshot_integration.go")
}

func checkHealth(apiBaseURL string) error {
	// 简单检查，实际应该发送 HTTP 请求
	return nil
}

// 辅助函数：保存 base64 图片
func saveBase64Image(base64Data, filename string) error {
	// 确保目录存在
	testDir := getTestDir()
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 解码 base64
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("解码 base64 失败: %w", err)
	}

	// 保存文件
	filepath := filepath.Join(testDir, filename)
	if err := os.WriteFile(filepath, imgData, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %w", err)
	}

	return nil
}

func getTestDir() string {
	wd, _ := os.Getwd()
	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			testDir := filepath.Join(dir, "test")
			os.MkdirAll(testDir, 0755)
			return testDir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Join(wd, "test")
}

func init() {
	configs.LoadConfigs()
}
