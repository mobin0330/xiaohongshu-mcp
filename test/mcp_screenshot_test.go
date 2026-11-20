package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// 注意：这个测试需要直接导入 main 包中的类型
// 由于 Go 的限制，我们需要在同一个包中运行
// 或者使用 go run 运行整个项目

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("=== 小红书 MCP 截图集成测试 ===")
	fmt.Println()
	fmt.Println("这个测试演示了 MCP 服务的工作流程:")
	fmt.Println("1. 启动持久浏览器")
	fmt.Println("2. 导航到小红书")
	fmt.Println("3. 截图当前页面")
	fmt.Println("4. 保存截图到 test 文件夹")
	fmt.Println("5. 关闭浏览器")
	fmt.Println()

	// 获取 test 目录
	testDir := getTestDir()
	fmt.Printf("截图将保存到: %s\n", testDir)
	fmt.Println()

	// 由于不能跨包导入 main 包中的类型，
	// 这里我们使用 HTTP API 方式演示
	fmt.Println("请使用以下步骤测试:")
	fmt.Println()
	fmt.Println("=== 步骤 1: 启动 MCP 服务 ===")
	fmt.Println("go run .")
	fmt.Println()
	fmt.Println("等待服务启动后，执行以下命令:")
	fmt.Println()
	fmt.Println("=== 步骤 2: 启动持久浏览器并打开小红书 ===")
	fmt.Println(`curl -X POST http://localhost:8080/api/v1/persistent_browser \`)
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"action": "start", "headless": false, "url": "https://www.xiaohongshu.com/explore"}'`)
	fmt.Println()
	fmt.Println("=== 步骤 3: 等待页面加载（约 3-5 秒）===")
	fmt.Println("sleep 5")
	fmt.Println()
	fmt.Println("=== 步骤 4: 截图当前页面（整页）===")
	fmt.Printf(`curl -s -X POST http://localhost:8080/api/v1/screenshot_current \`)
	fmt.Println()
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"full_page": true}' > /tmp/screenshot_response.json`)
	fmt.Println()
	fmt.Println("# 提取并保存 base64 图片")
	fmt.Println(`cat /tmp/screenshot_response.json | jq -r '.data.image_base64' | base64 -d > test/xiaohongshu_mcp_fullpage.png`)
	fmt.Println()
	fmt.Println("=== 步骤 5: 截图当前页面（视口）===")
	fmt.Printf(`curl -s -X POST http://localhost:8080/api/v1/screenshot_current \`)
	fmt.Println()
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"full_page": false}' > /tmp/screenshot_response2.json`)
	fmt.Println()
	fmt.Println(`cat /tmp/screenshot_response2.json | jq -r '.data.image_base64' | base64 -d > test/xiaohongshu_mcp_viewport.png`)
	fmt.Println()
	fmt.Println("=== 步骤 6: 关闭浏览器 ===")
	fmt.Printf(`curl -X POST http://localhost:8080/api/v1/persistent_browser \`)
	fmt.Println()
	fmt.Println(`  -H "Content-Type: application/json" \`)
	fmt.Println(`  -d '{"action": "close"}'`)
	fmt.Println()
	fmt.Println("=== 检查截图文件 ===")
	fmt.Println("ls -lh test/xiaohongshu_mcp_*.png")
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
