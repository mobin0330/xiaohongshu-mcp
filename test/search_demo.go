package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

func main() {
	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("=== 小红书搜索+互动测试脚本 ===")
	fmt.Println()

	// 检查 cookies 是否存在
	cookiePath := cookies.GetCookiesFilePath()
	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		fmt.Printf("未找到 cookies 文件: %s\n", cookiePath)
		fmt.Println("请先登录小红书，运行: go run cmd/login/main.go")
		fmt.Println()
		fmt.Println("登录步骤:")
		fmt.Println("1. 运行 go run cmd/login/main.go")
		fmt.Println("2. 浏览器会打开小红书网页并显示二维码")
		fmt.Println("3. 使用小红书 APP 扫描二维码登录")
		fmt.Println("4. 登录成功后 cookies 会自动保存")
		os.Exit(1)
	}
	fmt.Printf("找到 cookies 文件: %s\n", cookiePath)

	// 启动浏览器（非无头模式，方便观察）
	fmt.Println("正在启动浏览器...")
	b := browser.NewBrowser(false)
	defer b.Close()

	page := b.NewPage()
	defer page.Close()

	// 检查登录状态
	fmt.Println("检查登录状态...")
	loginAction := xiaohongshu.NewLogin(page)
	isLoggedIn, err := loginAction.CheckLoginStatus(context.Background())
	if err != nil {
		fmt.Printf("检查登录状态失败: %v\n", err)
		fmt.Println("可能是 cookies 已过期，请重新登录")
		os.Exit(1)
	}

	if !isLoggedIn {
		fmt.Println("登录状态已失效，请重新登录")
		fmt.Println("运行: go run cmd/login/main.go")
		os.Exit(1)
	}
	fmt.Println("登录状态有效 ✓")
	fmt.Println()

	// 执行搜索+互动测试
	fmt.Println("=== 开始搜索+互动测试 ===")
	testSearchAndInteract(page)

	fmt.Println()
	fmt.Println("=== 测试完成 ===")
}

// testSearchAndInteract 执行搜索并对每个关键词的第1个帖子进行点赞、关注、收藏
func testSearchAndInteract(page *rod.Page) {
	// 创建搜索动作
	searchAction := xiaohongshu.NewSearchAction(page)

	// 测试关键词
	keywords := []string{"旅游", "美食"}

	for _, keyword := range keywords {
		fmt.Printf("\n========================================\n")
		fmt.Printf("搜索关键词: %s\n", keyword)
		fmt.Printf("========================================\n")

		ctx := context.Background()
		feeds, err := searchAction.Search(ctx, keyword)
		if err != nil {
			fmt.Printf("搜索失败: %v\n", err)
			continue
		}

		fmt.Printf("成功获取到 %d 个 Feed\n", len(feeds))

		if len(feeds) == 0 {
			fmt.Println("没有找到任何帖子，跳过互动操作")
			continue
		}

		// 获取第1个帖子
		firstFeed := feeds[0]
		fmt.Printf("\n选择第1个帖子进行互动:\n")
		printFeedInfo(1, firstFeed)

		// 对该帖子执行点赞、关注、收藏
		performInteractions(page, firstFeed)
	}
}

// performInteractions 执行点赞、关注、收藏操作
func performInteractions(page *rod.Page, feed xiaohongshu.Feed) {
	ctx := context.Background()
	feedID := feed.ID
	xsecToken := feed.XsecToken
	userID := feed.NoteCard.User.UserID

	fmt.Println("\n--- 开始互动操作 ---")

	// 1. 点赞
	fmt.Println("\n[1/3] 执行点赞...")
	likeAction := xiaohongshu.NewLikeAction(page)
	if err := likeAction.Like(ctx, feedID, xsecToken); err != nil {
		fmt.Printf("点赞失败: %v\n", err)
	} else {
		fmt.Println("点赞成功 ✓")
	}

	// 2. 关注作者
	fmt.Println("\n[2/3] 执行关注...")
	followAction := xiaohongshu.NewFollowAction(page)
	// 使用帖子详情页方式关注
	if err := followAction.FollowUserFromFeed(ctx, feedID, xsecToken); err != nil {
		// 如果失败，尝试直接访问用户主页关注
		fmt.Printf("从帖子详情页关注失败，尝试访问用户主页: %v\n", err)
		if err := followAction.FollowUser(ctx, userID, xsecToken); err != nil {
			fmt.Printf("关注失败: %v\n", err)
		} else {
			fmt.Println("关注成功 ✓")
		}
	} else {
		fmt.Println("关注成功 ✓")
	}

	// 3. 收藏
	fmt.Println("\n[3/3] 执行收藏...")
	favoriteAction := xiaohongshu.NewFavoriteAction(page)
	if err := favoriteAction.Favorite(ctx, feedID, xsecToken); err != nil {
		fmt.Printf("收藏失败: %v\n", err)
	} else {
		fmt.Println("收藏成功 ✓")
	}

	fmt.Println("\n--- 互动操作完成 ---")
}

// printFeedInfo 打印 Feed 信息
func printFeedInfo(index int, feed xiaohongshu.Feed) {
	fmt.Printf("  ID: %s\n", feed.ID)
	fmt.Printf("  标题: %s\n", feed.NoteCard.DisplayTitle)
	fmt.Printf("  作者: %s (@%s)\n", feed.NoteCard.User.Nickname, feed.NoteCard.User.UserID)
	fmt.Printf("  类型: %s\n", feed.NoteCard.Type)
	fmt.Printf("  点赞数: %s\n", feed.NoteCard.InteractInfo.LikedCount)
}
