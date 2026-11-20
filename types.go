package main

import "github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"

// HTTP API 响应类型

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"`
}

// MCP 相关类型（用于内部转换）

// MCPToolResult MCP 工具结果（内部使用）
type MCPToolResult struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP 内容（内部使用）
type MCPContent struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// CommentLoadConfig 评论加载配置
type CommentLoadConfig struct {
	// 是否点击"更多回复"按钮
	ClickMoreReplies bool `json:"click_more_replies,omitempty"`
	// 回复数量阈值，超过这个数量的"更多"按钮将被跳过（0表示不跳过任何）
	MaxRepliesThreshold int `json:"max_replies_threshold,omitempty"`
	// 最大加载评论数（.parent-comment数量），0表示加载所有
	MaxCommentItems int `json:"max_comment_items,omitempty"`
	// 滚动速度等级: slow(慢速), normal(正常), fast(快速)
	ScrollSpeed string `json:"scroll_speed,omitempty"`
}

// FeedDetailRequest Feed详情请求
type FeedDetailRequest struct {
	FeedID          string             `json:"feed_id" binding:"required"`
	XsecToken       string             `json:"xsec_token" binding:"required"`
	LoadAllComments bool               `json:"load_all_comments,omitempty"`
	CommentConfig   *CommentLoadConfig `json:"comment_config,omitempty"`
}

type SearchFeedsRequest struct {
	Keyword string                   `json:"keyword" binding:"required"`
	Filters xiaohongshu.FilterOption `json:"filters,omitempty"`
}

// FeedDetailResponse Feed详情响应
type FeedDetailResponse struct {
	FeedID string `json:"feed_id"`
	Data   any    `json:"data"`
}

// PostCommentRequest 发表评论请求
type PostCommentRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// PostCommentResponse 发表评论响应
type PostCommentResponse struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ReplyCommentRequest 回复评论请求
type ReplyCommentRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
	CommentID string `json:"comment_id" binding:"required_without=UserID"`
	UserID    string `json:"user_id" binding:"required_without=CommentID"`
	Content   string `json:"content" binding:"required"`
}

// ReplyCommentResponse 回复评论响应
type ReplyCommentResponse struct {
	FeedID          string `json:"feed_id"`
	TargetCommentID string `json:"target_comment_id,omitempty"`
	TargetUserID    string `json:"target_user_id,omitempty"`
	Success         bool   `json:"success"`
	Message         string `json:"message"`
}

// UserProfileRequest 用户主页请求
type UserProfileRequest struct {
	UserID    string `json:"user_id" binding:"required"`
	XsecToken string `json:"xsec_token" binding:"required"`
}

// ActionResult 通用动作响应（点赞/收藏等）
type ActionResult struct {
	FeedID  string `json:"feed_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ScreenshotRequest 浏览器截图请求
type ScreenshotRequest struct {
	URL          string `json:"url" binding:"required"`  // 目标网页 URL
	FullPage     bool   `json:"full_page,omitempty"`     // 是否截取整页
	Width        int    `json:"width,omitempty"`         // 视口宽度（默认 1280）
	Height       int    `json:"height,omitempty"`        // 视口高度（默认 720）
	Selector     string `json:"selector,omitempty"`      // CSS 选择器，截取指定元素
	WaitSelector string `json:"wait_selector,omitempty"` // 等待指定元素出现后再截图
	WaitTimeout  int    `json:"wait_timeout,omitempty"`  // 等待超时（秒，默认 30）
}

// ScreenshotResponse 浏览器截图响应
type ScreenshotResponse struct {
	ImageBase64 string `json:"image_base64"` // Base64 编码的图片
	URL         string `json:"url"`          // 实际截图的 URL
	Width       int    `json:"width"`        // 图片宽度
	Height      int    `json:"height"`       // 图片高度
}

// PersistentBrowserRequest 持久浏览器操作请求
type PersistentBrowserRequest struct {
	Action   string `json:"action" binding:"required"` // 操作类型: start|close|navigate|status
	Headless bool   `json:"headless,omitempty"`        // 是否无头模式（默认 false，显示界面）
	URL      string `json:"url,omitempty"`             // 启动后要导航的 URL（可选）
}

// PersistentBrowserResponse 持久浏览器操作响应
type PersistentBrowserResponse struct {
	Running bool   `json:"running"`       // 浏览器是否正在运行
	URL     string `json:"url,omitempty"` // 当前页面 URL
	Message string `json:"message"`       // 操作结果消息
}

// FollowUserRequest 关注用户请求
// FollowUserRequest 关注用户请求
type FollowUserRequest struct {
	UserID    string `json:"user_id" binding:"required"`    // 用户ID
	XsecToken string `json:"xsec_token" binding:"required"` // 访问令牌
}

// FollowUserFromFeedRequest 从帖子详情页关注用户请求
type FollowUserFromFeedRequest struct {
	FeedID    string `json:"feed_id" binding:"required"`    // 帖子ID
	XsecToken string `json:"xsec_token" binding:"required"` // 访问令牌
}
