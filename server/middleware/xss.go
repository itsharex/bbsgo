package middleware

import (
	"regexp"
	"strings"
)

// XSS 攻击的常见模式
var xssPatterns = []string{
	`<script[^>]*>.*?</script>`,
	`javascript:`,
	`onerror\s*=`,
	`onload\s*=`,
	`onclick\s*=`,
	`onmouseover\s*=`,
	`onfocus\s*=`,
	`onblur\s*=`,
	`onchange\s*=`,
	`onsubmit\s*=`,
	`<iframe[^>]*>.*?</iframe>`,
	`<object[^>]*>.*?</object>`,
	`<embed[^>]*>`,
	`<link[^>]*>`,
	`<meta[^>]*>`,
	`expression\s*\(`,
	`url\s*\(`,
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

// RemoveScriptTags 移除可能包含脚本的内容
// 用于用户输入的初步过滤
func RemoveScriptTags(content string) string {
	if content == "" {
		return ""
	}

	// 移除 <script> 标签及其内容
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	content = scriptPattern.ReplaceAllString(content, "")

	// 移除 <style> 标签及其内容
	stylePattern := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	content = stylePattern.ReplaceAllString(content, "")

	// 移除 on* 事件属性
	eventPattern := regexp.MustCompile(`(?i)\s+on\w+\s*=\s*["'][^"']*["']`)
	content = eventPattern.ReplaceAllString(content, "")

	// 移除 javascript: 伪协议
	jsProtocol := regexp.MustCompile(`(?i)javascript:`)
	content = jsProtocol.ReplaceAllString(content, "")

	// 移除 data: 伪协议（可能用于 XSS）
	dataProtocol := regexp.MustCompile(`(?i)data:`)
	content = dataProtocol.ReplaceAllString(content, "")

	return content
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
	result := RemoveScriptTags(content)

	return result
}
