package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	APIBase = "http://localhost:8080/api/v1"
	TestDir = "test"
)

// API 响应类型
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ScreenshotData struct {
	ImageBase64 string `json:"image_base64"`
	URL         string `json:"url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

func main() {
	fmt.Println("=== 小红书 MCP 截图集成测试 ===")
	fmt.Println()

	// 检查服务
	if err := checkService(); err != nil {
		fmt.Printf("服务检查失败: %v\n", err)
		fmt.Println("请先启动 MCP 服务: go run .")
		os.Exit(1)
	}
	fmt.Println("✓ MCP 服务运行正常")
	fmt.Println()

	// 确保 test 目录存在
	testDir := getTestDir()
	os.MkdirAll(testDir, 0755)

	// 步骤 1: 关闭之前的浏览器（如果有）
	fmt.Println("步骤 1: 清理之前的浏览器实例...")
	postJSON("/persistent_browser", map[string]interface{}{"action": "close"})
	time.Sleep(1 * time.Second)
	fmt.Println()

	// 步骤 2: 启动持久浏览器并打开小红书
	fmt.Println("步骤 2: 启动持久浏览器并打开小红书...")
	_, err := postJSON("/persistent_browser", map[string]interface{}{
		"action":   "start",
		"headless": false,
		"url":      "https://www.xiaohongshu.com/explore",
	})
	if err != nil {
		fmt.Printf("启动浏览器失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 浏览器已启动，正在加载小红书...")
	fmt.Println()

	// 等待页面加载
	fmt.Println("步骤 3: 等待页面加载 (5秒)...")
	time.Sleep(5 * time.Second)
	fmt.Println()

	// 步骤 4: 截图整页
	fmt.Println("步骤 4: 截取整页截图...")
	resp, err := postJSON("/screenshot_current", map[string]interface{}{
		"full_page": true,
	})
	if err != nil {
		fmt.Printf("截图失败: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	if resp.Success && resp.Data != nil {
		dataJSON, _ := json.Marshal(resp.Data)
		var screenshot ScreenshotData
		json.Unmarshal(dataJSON, &screenshot)

		filename := "xiaohongshu_mcp_fullpage.png"
		if err := saveBase64Image(screenshot.ImageBase64, filepath.Join(testDir, filename)); err != nil {
			fmt.Printf("保存截图失败: %v\n", err)
		} else {
			fmt.Printf("✓ 整页截图已保存: %s (%dx%d)\n", filename, screenshot.Width, screenshot.Height)
		}
	} else {
		fmt.Printf("截图响应异常: success=%v, data=%v, message=%s\n", resp.Success, resp.Data, resp.Message)
	}
	fmt.Println()

	// 步骤 5: 截图视口
	fmt.Println("步骤 5: 截取视口截图...")
	resp, err = postJSON("/screenshot_current", map[string]interface{}{
		"full_page": false,
	})
	if err != nil {
		fmt.Printf("截图失败: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	if resp.Success && resp.Data != nil {
		dataJSON, _ := json.Marshal(resp.Data)
		var screenshot ScreenshotData
		json.Unmarshal(dataJSON, &screenshot)

		filename := "xiaohongshu_mcp_viewport.png"
		if err := saveBase64Image(screenshot.ImageBase64, filepath.Join(testDir, filename)); err != nil {
			fmt.Printf("保存截图失败: %v\n", err)
		} else {
			fmt.Printf("✓ 视口截图已保存: %s (%dx%d)\n", filename, screenshot.Width, screenshot.Height)
		}
	}
	fmt.Println()

	// 步骤 6: 关闭浏览器
	fmt.Println("步骤 6: 关闭浏览器...")
	postJSON("/persistent_browser", map[string]interface{}{"action": "close"})
	fmt.Println("✓ 浏览器已关闭")
	fmt.Println()

	// 显示结果
	fmt.Println("=== 测试完成 ===")
	fmt.Printf("截图文件保存在: %s\n", testDir)
	listScreenshots(testDir)
}

func checkService() error {
	resp, err := http.Get(APIBase + "/login/status")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func postJSON(path string, data interface{}) (*Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(APIBase+path, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func saveBase64Image(base64Data, filepath string) error {
	imgData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, imgData, 0644)
}

func getTestDir() string {
	wd, _ := os.Getwd()
	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return filepath.Join(dir, TestDir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return filepath.Join(wd, TestDir)
}

func cleanup() {
	fmt.Println("清理: 关闭浏览器...")
	postJSON("/persistent_browser", map[string]interface{}{"action": "close"})
}

func listScreenshots(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".png" {
			info, _ := file.Info()
			fmt.Printf("  - %s (%.2f KB)\n", file.Name(), float64(info.Size())/1024)
		}
	}
}
