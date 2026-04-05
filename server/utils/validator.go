package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// PasswordStrength 密码强度等级
type PasswordStrength int

const (
	PasswordStrengthWeak       PasswordStrength = 1
	PasswordStrengthMedium     PasswordStrength = 2
	PasswordStrengthStrong     PasswordStrength = 3
	PasswordStrengthVeryStrong PasswordStrength = 4
)

// ValidatePassword 验证密码强度
// 返回密码强度等级和错误信息
func ValidatePassword(password string) (PasswordStrength, error) {
	if len(password) == 0 {
		return 0, errors.New("密码不能为空")
	}

	// 最小长度检查
	if len(password) < 8 {
		return 0, errors.New("密码长度至少8位")
	}

	// 最大长度检查
	if len(password) > 128 {
		return 0, errors.New("密码长度不能超过128位")
	}

	// 检查密码复杂度
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 计算复杂度得分
	score := 0
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	// 长度加分
	if len(password) >= 12 {
		score++
	}
	if len(password) >= 16 {
		score++
	}

	// 至少包含3种字符类型
	if score < 3 {
		return 0, errors.New("密码必须包含大写字母、小写字母、数字和特殊字符中的至少3种")
	}

	// 检查常见弱密码
	if isWeakPassword(password) {
		return 0, errors.New("密码过于简单，请使用更复杂的密码")
	}

	// 计算强度等级
	var strength PasswordStrength
	switch {
	case score >= 5:
		strength = PasswordStrengthVeryStrong
	case score >= 4:
		strength = PasswordStrengthStrong
	case score >= 3:
		strength = PasswordStrengthMedium
	default:
		strength = PasswordStrengthWeak
	}

	return strength, nil
}

// isWeakPassword 检查是否为常见弱密码
func isWeakPassword(password string) bool {
	// 常见弱密码列表
	weakPasswords := []string{
		"password", "12345678", "qwerty", "abc123", "11111111",
		"password123", "admin123", "letmein", "welcome", "monkey",
		"dragon", "master", "login", "admin", "qwertyuiop",
		"1234567890", "password1", "qwerty123", "password!",
	}

	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if lowerPassword == weak {
			return true
		}
	}

	// 检查连续字符（只检查字母键盘序列，不检查纯数字序列）
	if hasConsecutiveKeyboardChars(password, 4) {
		return true
	}

	// 检查重复字符
	if hasRepeatingChars(password, 4) {
		return true
	}

	return false
}

// hasConsecutiveChars 检查是否有连续字符（已废弃，改用 hasConsecutiveKeyboardChars）
func hasConsecutiveChars(s string, count int) bool {
	return hasConsecutiveKeyboardChars(s, count)
}

// hasConsecutiveKeyboardChars 检查是否有键盘连续字符（如 qwerty, asdfgh, zxcvbn 等）
// 只检查字母键盘序列，不检查纯数字序列
func hasConsecutiveKeyboardChars(s string, count int) bool {
	if len(s) < count {
		return false
	}

	// 键盘行序列（按键盘顺序）
	keyboardRows := []string{
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
		"1234567890",
	}

	lowerS := strings.ToLower(s)

	for _, row := range keyboardRows {
		// 检查 s 中是否包含该键盘序列的连续子串
		for i := 0; i <= len(lowerS)-count; i++ {
			substr := lowerS[i : i+count]
			// 检查这个子串是否是键盘序列的子串
			for j := 0; j <= len(row)-count; j++ {
				if row[j:j+count] == substr {
					return true
				}
			}
		}
	}

	return false
}

// hasRepeatingChars 检查是否有重复字符
func hasRepeatingChars(s string, count int) bool {
	if len(s) < count {
		return false
	}

	repeating := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			repeating++
			if repeating >= count {
				return true
			}
		} else {
			repeating = 1
		}
	}

	return false
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("用户名不能为空")
	}

	if len(username) < 3 {
		return errors.New("用户名长度至少3位")
	}

	if len(username) > 20 {
		return errors.New("用户名长度不能超过20位")
	}

	// 只允许字母、数字、下划线
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	// 不能以数字开头
	if unicode.IsNumber(rune(username[0])) {
		return errors.New("用户名不能以数字开头")
	}

	return nil
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("邮箱不能为空")
	}

	// 简单的邮箱格式验证
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("邮箱格式不正确")
	}

	return nil
}
