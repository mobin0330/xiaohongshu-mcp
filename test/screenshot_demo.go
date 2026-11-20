package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
)

func main() {
	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("=== 小红书网页截图测试 ===")
	fmt.Println()

	// 启动浏览器（非无头模式，方便观察）
	fmt.Println("正在启动浏览器...")
	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	// 设置视口大小
	page.MustSetViewport(1280, 720, 1.0, false)
	fmt.Println("浏览器已启动，视口设置为 1280x720")

	// 导航到小红书
	fmt.Println("\n正在打开小红书...")
	page.MustNavigate("https://www.xiaohongshu.com/explore")
	page.MustWaitLoad()
	fmt.Println("小红书页面已加载")

	// 等待页面稳定
	fmt.Println("等待页面稳定...")
	page.MustWaitStable()
	time.Sleep(2 * time.Second) // 额外等待 2 秒确保内容加载完成

	// 获取 test 目录路径
	testDir := getTestDir()

	// 截图 1: 当前视口
	screenshotPath1 := filepath.Join(testDir, "xiaohongshu_viewport.png")
	fmt.Printf("\n正在截取视口截图: %s\n", screenshotPath1)
	page.MustScreenshot(screenshotPath1)
	fmt.Println("视口截图已保存 ✓")

	// 截图 2: 整页截图
	screenshotPath2 := filepath.Join(testDir, "xiaohongshu_fullpage.png")
	fmt.Printf("\n正在截取整页截图: %s\n", screenshotPath2)
	page.MustScreenshotFullPage(screenshotPath2)
	fmt.Println("整页截图已保存 ✓")

	// 截图 3: 指定元素（如果有）
	// 尝试截取第一个笔记卡片
	fmt.Println("\n尝试截取第一个笔记卡片...")
	elements, err := page.Elements("section.note-item")
	elementScreenshotSaved := false
	if err == nil && len(elements) > 0 {
		screenshotPath3 := filepath.Join(testDir, "xiaohongshu_first_note.png")
		fmt.Printf("正在截取第一个笔记卡片: %s\n", screenshotPath3)
		elements[0].MustScreenshot(screenshotPath3)
		elementScreenshotSaved = true
		fmt.Println("笔记卡片截图已保存 ✓")
	} else {
		if err != nil {
			fmt.Printf("查找笔记卡片失败: %v\n", err)
		}
		fmt.Println("未找到笔记卡片元素，跳过元素截图")
	}

	fmt.Println("\n=== 截图测试完成 ===")
	fmt.Printf("\n截图文件保存在: %s\n", testDir)
	fmt.Println("- xiaohongshu_viewport.png: 当前视口截图")
	fmt.Println("- xiaohongshu_fullpage.png: 整页截图")
	if elementScreenshotSaved {
		fmt.Println("- xiaohongshu_first_note.png: 第一个笔记卡片截图")
	}

	// 保持浏览器打开 5 秒，方便用户查看
	fmt.Println("\n5 秒后关闭浏览器...")
	time.Sleep(5 * time.Second)
}

// getTestDir 获取 test 目录路径
func getTestDir() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		// 如果失败，使用当前目录
		return "test"
	}

	// 尝试找到项目根目录
	// 从当前目录向上查找，直到找到 go.mod 文件
	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// 找到了项目根目录
			testDir := filepath.Join(dir, "test")
			// 确保目录存在
			os.MkdirAll(testDir, 0755)
			return testDir
		}

		// 向上查找父目录
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达根目录，停止查找
			break
		}
		dir = parent
	}

	// 如果没找到，使用当前目录下的 test
	testDir := filepath.Join(wd, "test")
	os.MkdirAll(testDir, 0755)
	return testDir
}
