package utils

import (
	"regexp"
	"strings"
)

// 敏感信息模式
var sensitivePatterns = []struct {
	pattern     *regexp.Regexp
	replacement string
}{
	// JWT Token
	{
		regexp.MustCompile(`(?i)(bearer\s+)[a-zA-Z0-9\-_]+\.?[a-zA-Z0-9\-_]*\.?[a-zA-Z0-9\-_]*`),
		"${1}***REDACTED***",
	},
	// Email
	{
		regexp.MustCompile(`(?i)([a-zA-Z0-9._%+-]+)@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
		"${1}***@***.***",
	},
	// Password in logs
	{
		regexp.MustCompile(`(?i)(password["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// Password hash
	{
		regexp.MustCompile(`(?i)(password_hash["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// API Key
	{
		regexp.MustCompile(`(?i)(api[_-]?key["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// Secret
	{
		regexp.MustCompile(`(?i)(secret["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// Token
	{
		regexp.MustCompile(`(?i)(token["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// Authorization header
	{
		regexp.MustCompile(`(?i)(authorization["\s:=]+)["']?[^"'\s,}]+["']?`),
		"${1}***REDACTED***",
	},
	// Phone number (Chinese)
	{
		regexp.MustCompile(`(?i)(1[3-9]\d)\d{4}(\d{4})`),
		"${1}****${2}",
	},
	// ID Card (Chinese)
	{
		regexp.MustCompile(`(?i)(\d{6})\d{8}(\d{4})`),
		"${1}********${2}",
	},
	// Bank card
	{
		regexp.MustCompile(`(?i)(\d{4})\d{8,12}(\d{4})`),
		"${1}********${2}",
	},
	// Credit card
	{
		regexp.MustCompile(`(?i)(\d{4})\d{8,11}(\d{4})`),
		"${1}********${2}",
	},
}

// MaskSensitiveInfo 脱敏敏感信息
// 用于日志输出前对敏感信息进行脱敏处理
func MaskSensitiveInfo(content string) string {
	if content == "" {
		return ""
	}

	result := content
	for _, sp := range sensitivePatterns {
		result = sp.pattern.ReplaceAllString(result, sp.replacement)
	}

	return result
}

// MaskEmail 脱敏邮箱地址
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	name := parts[0]
	domain := parts[1]

	if len(name) <= 2 {
		return name[:1] + "***@" + domain
	}

	return name[:2] + "***@" + domain
}

// MaskPhone 脱敏手机号
func MaskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}

	return phone[:3] + "****" + phone[7:]
}

// MaskIDCard 脱敏身份证号
func MaskIDCard(idCard string) string {
	if len(idCard) != 18 && len(idCard) != 15 {
		return idCard
	}

	return idCard[:6] + "********" + idCard[len(idCard)-4:]
}

// MaskBankCard 脱敏银行卡号
func MaskBankCard(cardNumber string) string {
	if len(cardNumber) < 8 {
		return cardNumber
	}

	return cardNumber[:4] + "********" + cardNumber[len(cardNumber)-4:]
}

// MaskName 脱敏姓名
func MaskName(name string) string {
	if name == "" {
		return ""
	}

	runes := []rune(name)
	length := len(runes)

	if length == 1 {
		return "*"
	} else if length == 2 {
		return string(runes[0]) + "*"
	} else {
		return string(runes[0]) + strings.Repeat("*", length-2) + string(runes[length-1])
	}
}

// MaskToken 脱敏 Token
func MaskToken(token string) string {
	if token == "" {
		return ""
	}

	if len(token) <= 10 {
		return "***REDACTED***"
	}

	return token[:5] + "..." + token[len(token)-5:]
}

// MaskPassword 脱敏密码
func MaskPassword(password string) string {
	if password == "" {
		return ""
	}

	return "***REDACTED***"
}

// MaskURL 脱敏 URL 中的敏感参数
func MaskURL(url string) string {
	if url == "" {
		return ""
	}

	// 脱敏 URL 中的 token 参数
	tokenPattern := regexp.MustCompile(`(?i)(token=)[^&]+`)
	url = tokenPattern.ReplaceAllString(url, "${1}***REDACTED***")

	// 脱敏 URL 中的 password 参数
	passwordPattern := regexp.MustCompile(`(?i)(password=)[^&]+`)
	url = passwordPattern.ReplaceAllString(url, "${1}***REDACTED***")

	// 脱敏 URL 中的 key 参数
	keyPattern := regexp.MustCompile(`(?i)(key=)[^&]+`)
	url = keyPattern.ReplaceAllString(url, "${1}***REDACTED***")

	// 脱敏 URL 中的 secret 参数
	secretPattern := regexp.MustCompile(`(?i)(secret=)[^&]+`)
	url = secretPattern.ReplaceAllString(url, "${1}***REDACTED***")

	return url
}
