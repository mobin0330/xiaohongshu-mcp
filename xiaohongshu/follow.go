package xiaohongshu

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

// 关注功能相关常量
const (
	// 关注按钮状态文本
	statusFollowed     = "已关注"
	statusMutualFollow = "互相关注"
	statusFollowing    = "Following"

	// 关注按钮文本
	btnFollow        = "关注"
	btnFollowPlus    = "+ 关注"
	btnFollowNoSpace = "+关注"
	btnFollowEn      = "Follow"

	// 默认超时时间
	followTimeout = 60 * time.Second
	// 页面加载等待时间
	pageLoadWait = 2 * time.Second
	// 点击后等待时间
	clickWait = 2 * time.Second
)

// FollowAction 负责处理关注用户相关交互
type FollowAction struct {
	page *rod.Page
}

func NewFollowAction(page *rod.Page) *FollowAction {
	pp := page.Timeout(followTimeout)
	return &FollowAction{page: pp}
}

// FollowUser 关注指定用户
// userID: 用户ID
// xsecToken: 安全token（可从帖子详情页获取）
func (f *FollowAction) FollowUser(ctx context.Context, userID, xsecToken string) error {
	page := f.page.Context(ctx).Timeout(followTimeout)

	// 构建用户主页URL
	url := makeUserProfileURL(userID, xsecToken)
	logrus.Infof("打开用户主页进行关注: %s", url)

	page.MustNavigate(url)
	page.MustWaitDOMStable()
	time.Sleep(pageLoadWait)

	// 用户主页的关注按钮选择器
	selectors := []string{
		".user-info .follow-btn",
		".user-info .follow-button",
		".user-profile-header .follow-btn",
		"button[class*='follow']",
	}

	return f.findAndClickFollowButton(page, selectors, userID)
}

// FollowUserFromFeed 从帖子详情页关注用户
// 直接打开帖子详情页，点击作者区域的关注按钮
func (f *FollowAction) FollowUserFromFeed(ctx context.Context, feedID, xsecToken string) error {
	page := f.page.Context(ctx).Timeout(followTimeout)

	// 打开帖子详情页
	url := makeFeedDetailURL(feedID, xsecToken)
	logrus.Infof("打开帖子详情页进行关注: %s", url)

	page.MustNavigate(url)
	page.MustWaitDOMStable()
	time.Sleep(pageLoadWait)

	// 帖子详情页的关注按钮选择器（更全面）
	selectors := []string{
		".author-container .follow-btn",
		".author .follow-btn",
		".note-author .follow-button",
		".author-info .follow-btn",
		"button.follow-btn",
		"[class*='follow']",
		"div[class*='follow']",
		"span[class*='follow']",
	}

	if err := f.findAndClickFollowButton(page, selectors, feedID); err == nil {
		return nil
	}

	// 如果CSS选择器都失败，尝试JavaScript方式
	return f.findAndClickFollowButtonByJS(page, feedID)
}

// findAndClickFollowButton 查找并点击关注按钮
func (f *FollowAction) findAndClickFollowButton(page *rod.Page, selectors []string, targetID string) error {
	for _, selector := range selectors {
		elements, err := page.Elements(selector)
		if err != nil {
			continue
		}

		for _, btn := range elements {
			text, err := btn.Text()
			if err != nil {
				continue
			}

			text = strings.TrimSpace(text)

			// 检查是否已经关注
			if f.isAlreadyFollowing(text) {
				logrus.Infof("目标 %s 已经关注过了", targetID)
				return nil
			}

			// 点击关注按钮
			if f.isFollowButton(text) {
				btn.MustClick()
				time.Sleep(clickWait)
				logrus.Infof("已关注目标 %s", targetID)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到关注按钮")
}

// findAndClickFollowButtonByJS 通过JavaScript查找并点击关注按钮
func (f *FollowAction) findAndClickFollowButtonByJS(page *rod.Page, targetID string) error {
	logrus.Infof("尝试通过 JavaScript 查找关注按钮...")

	result := page.MustEval(`() => {
		const allElements = document.querySelectorAll('button, div, span, a');
		for (const el of allElements) {
			const text = el.textContent.trim();
			if (text === '关注' || text === '+ 关注' || text === '+关注') {
				el.click();
				return 'clicked';
			}
			if (text === '已关注' || text === '互相关注') {
				return 'already_followed';
			}
		}
		return 'not_found';
	}`).String()

	switch result {
	case "clicked":
		time.Sleep(clickWait)
		logrus.Infof("已通过 JavaScript 点击关注按钮")
		return nil
	case "already_followed":
		logrus.Infof("目标已经关注过了")
		return nil
	default:
		return fmt.Errorf("未找到关注按钮或无法关注目标 %s", targetID)
	}
}

// isAlreadyFollowing 检查文本是否表示已关注状态
func (f *FollowAction) isAlreadyFollowing(text string) bool {
	return text == statusFollowed || text == statusMutualFollow || text == statusFollowing
}

// isFollowButton 检查文本是否表示关注按钮
func (f *FollowAction) isFollowButton(text string) bool {
	return text == btnFollow || text == btnFollowPlus || text == btnFollowEn || text == btnFollowNoSpace
}
