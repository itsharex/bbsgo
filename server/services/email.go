package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"bbsgo/config"
)

type EmailService struct {
	enabled  bool
	host     string
	port     int
	user     string
	password string
	from     string
	fromName string
}

func NewEmailService() *EmailService {
	return &EmailService{
		enabled:  config.GetConfigBool("email_enabled", false),
		host:     config.GetConfig("email_host"),
		port:     config.GetConfigInt("email_port", 465),
		user:     config.GetConfig("email_user"),
		password: config.GetConfig("email_password"),
		from:     config.GetConfig("email_from"),
		fromName: config.GetConfig("email_from_name"),
	}
}

func (s *EmailService) Send(to, subject, body string) error {
	if !s.enabled {
		return fmt.Errorf("email service is disabled")
	}

	// 打印邮件配置，方便定位问题
	log.Printf("[EMAIL] 配置信息: host=%s, port=%d, user=%s, from=%s, fromName=%s",
		s.host, s.port, s.user, s.from, s.fromName)
	log.Printf("[EMAIL] 目标邮箱: %s", to)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	if s.port == 465 {
		return s.sendMailWithSSL(addr, s.user, s.password, s.from, to, subject, body)
	}
	return s.sendMailWithTLS(addr, s.user, s.password, s.from, to, subject, body)
}

func (s *EmailService) sendMailWithTLS(addr, user, password, from, to, subject, body string) error {
	auth := smtp.PlainAuth("", user, password, strings.Split(addr, ":")[0])

	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: strings.Split(addr, ":")[0]})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, strings.Split(addr, ":")[0])
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Quit()

	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	msg := s.buildMessage(from, to, subject, body)
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return w.Close()
}

func (s *EmailService) sendMailWithSSL(addr, user, password, from, to, subject, body string) error {
	auth := smtp.PlainAuth("", user, password, strings.Split(addr, ":")[0])

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         strings.Split(addr, ":")[0],
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, strings.Split(addr, ":")[0])
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Quit()

	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %v", err)
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	msg := s.buildMessage(from, to, subject, body)
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return w.Close()
}

func (s *EmailService) buildMessage(from, to, subject, body string) string {
	msg := fmt.Sprintf("From: %s <%s>\r\n", s.fromName, from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body
	return msg
}

func SendVerificationCode(to, code string) error {
	service := NewEmailService()
	if !service.enabled {
		return fmt.Errorf("email service is disabled")
	}

	subject := "验证码 - " + service.fromName
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>您的验证码</h2>
			<p>您的验证码是：<strong style="font-size: 24px; color: #4F46E5;">%s</strong></p>
			<p>验证码有效期为5分钟，请尽快使用。</p>
			<p>如果这不是您的操作，请忽略此邮件。</p>
			<br>
			<p>%s</p>
		</body>
		</html>
	`, code, service.fromName)

	return service.Send(to, subject, body)
}
