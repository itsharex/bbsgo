package middleware

import (
	"html"
	"regexp"
	"strings"
)

// XSS 攻击的常见模式
var xssPatterns = []string{
	`<script[^>]*>.*?</script>`,
	`<script[^>]*/>`,
	`javascript:`,
	`vbscript:`,
	`onabort\s*=`,
	`onblur\s*=`,
	`onchange\s*=`,
	`onclick\s*=`,
	`ondblclick\s*=`,
	`onerror\s*=`,
	`onfocus\s*=`,
	`onkeydown\s*=`,
	`onkeypress\s*=`,
	`onkeyup\s*=`,
	`onload\s*=`,
	`onmousedown\s*=`,
	`onmousemove\s*=`,
	`onmouseout\s*=`,
	`onmouseover\s*=`,
	`onmouseup\s*=`,
	`onreset\s*=`,
	`onresize\s*=`,
	`onscroll\s*=`,
	`onselect\s*=`,
	`onsubmit\s*=`,
	`onunload\s*=`,
	`<iframe[^>]*>.*?</iframe>`,
	`<iframe[^>]*/>`,
	`<object[^>]*>.*?</object>`,
	`<object[^>]*/>`,
	`<embed[^>]*>`,
	`<embed[^>]*/>`,
	`<link[^>]*>`,
	`<meta[^>]*>`,
	`<style[^>]*>.*?</style>`,
	`expression\s*\(`,
	`url\s*\(`,
	`behavior\s*:`,
	`-moz-binding\s*:`,
	`@import`,
	`document\.cookie`,
	`document\.write`,
	`window\.location`,
	`window\.open`,
	`eval\s*\(`,
	`setTimeout\s*\(`,
	`setInterval\s*\(`,
	`Function\s*\(`,
	`<base[^>]*>`,
	`<form[^>]*>`,
	`<input[^>]*>`,
	`<button[^>]*>`,
	`<textarea[^>]*>`,
	`<select[^>]*>`,
	`<applet[^>]*>`,
	`<body[^>]*>`,
	`<html[^>]*>`,
	`<svg[^>]*>.*?</svg>`,
	`<math[^>]*>.*?</math>`,
}

// compiledPatterns 编译后的正则表达式
var compiledPatterns = make([]*regexp.Regexp, 0, len(xssPatterns))

func init() {
	for _, pattern := range xssPatterns {
		compiledPatterns = append(compiledPatterns, regexp.MustCompile(`(?i)`+pattern))
	}
}

// ContainsXSS 检测内容是否包含潜在的 XSS 攻击模式
// 这个函数只做基础检测，不能完全替代专业的 HTML sanitization 库
func ContainsXSS(content string) bool {
	if content == "" {
		return false
	}

	for _, pattern := range compiledPatterns {
		if pattern.MatchString(content) {
			return true
		}
	}

	return false
}

// SanitizeHTML 对 HTML 内容进行安全清理
// 移除危险的标签和属性
func SanitizeHTML(content string) string {
	if content == "" {
		return ""
	}

	// 移除 script 标签及其内容
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	content = scriptPattern.ReplaceAllString(content, "")

	// 移除 style 标签及其内容
	stylePattern := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	content = stylePattern.ReplaceAllString(content, "")

	// 移除 iframe 标签
	iframePattern := regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`)
	content = iframePattern.ReplaceAllString(content, "")

	// 移除 object 标签
	objectPattern := regexp.MustCompile(`(?i)<object[^>]*>.*?</object>`)
	content = objectPattern.ReplaceAllString(content, "")

	// 移除 embed 标签
	embedPattern := regexp.MustCompile(`(?i)<embed[^>]*/?>`)
	content = embedPattern.ReplaceAllString(content, "")

	// 移除所有事件处理器属性
	eventPattern := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	content = eventPattern.ReplaceAllString(content, "")

	// 移除 javascript: 伪协议
	jsProtocol := regexp.MustCompile(`(?i)javascript\s*:`)
	content = jsProtocol.ReplaceAllString(content, "")

	// 移除 vbscript: 伪协议
	vbsProtocol := regexp.MustCompile(`(?i)vbscript\s*:`)
	content = vbsProtocol.ReplaceAllString(content, "")

	// 移除 data: 伪协议（可能用于 XSS）
	dataProtocol := regexp.MustCompile(`(?i)data\s*:`)
	content = dataProtocol.ReplaceAllString(content, "")

	// 移除 expression()
	exprPattern := regexp.MustCompile(`(?i)expression\s*\([^)]*\)`)
	content = exprPattern.ReplaceAllString(content, "")

	// 移除 behavior:
	behaviorPattern := regexp.MustCompile(`(?i)behavior\s*:`)
	content = behaviorPattern.ReplaceAllString(content, "")

	// 移除 -moz-binding:
	mozBindingPattern := regexp.MustCompile(`(?i)-moz-binding\s*:`)
	content = mozBindingPattern.ReplaceAllString(content, "")

	return content
}

// EscapeHTML 转义 HTML 特殊字符
// 用于纯文本内容的转义
func EscapeHTML(content string) string {
	if content == "" {
		return ""
	}
	return html.EscapeString(content)
}

// UnescapeHTML 反转义 HTML 特殊字符
func UnescapeHTML(content string) string {
	if content == "" {
		return ""
	}
	return html.UnescapeString(content)
}

// SanitizeContent 对用户发布的内容进行基础的安全处理
// 注意：这个函数不能完全替代专业的 HTML sanitization 库
func SanitizeContent(content string) string {
	if content == "" {
		return ""
	}

	// 移除换行和空白符后检查是否为空
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}

	// 移除潜在的 XSS 模式
	result := SanitizeHTML(content)

	return result
}

// SanitizeURL 清理 URL，移除危险的伪协议
func SanitizeURL(url string) string {
	if url == "" {
		return ""
	}

	// 转换为小写进行检测
	lowerURL := strings.ToLower(url)

	// 检查危险的伪协议
	dangerousProtocols := []string{
		"javascript:",
		"vbscript:",
		"data:",
	}

	for _, protocol := range dangerousProtocols {
		if strings.HasPrefix(lowerURL, protocol) {
			return ""
		}
	}

	return url
}

// ValidateContentType 验证内容类型是否安全
func ValidateContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	// 允许的内容类型
	allowedTypes := map[string]bool{
		"text/plain":                true,
		"text/html":                 true,
		"text/markdown":             true,
		"image/jpeg":                true,
		"image/jpg":                 true,
		"image/png":                 true,
		"image/gif":                 true,
		"image/webp":                true,
		"image/bmp":                 true,
		"video/mp4":                 true,
		"video/webm":                true,
		"video/ogg":                 true,
		"application/json":          true,
		"application/pdf":           true,
	}

	return allowedTypes[contentType]
}
