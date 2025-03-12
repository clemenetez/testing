package services

import (
	"fmt"
	"net/smtp"

	"github.com/your_username/testing/config" // Замініть на ваш шлях
)

func SendVerificationEmail(email, verificationCode string) error {
	auth := smtp.PlainAuth("", config.SMTP_USERNAME, config.SMTP_PASSWORD, config.SMTP_HOST)

	to := []string{email}
	subject := "Підтвердження пошти"
	body := fmt.Sprintf(`
        <h1>Код підтвердження</h1>
        <p>Ваш код: <strong>%s</strong></p>
    `, verificationCode)

	msg := []byte(
		"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" + body,
	)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", config.SMTP_HOST, config.SMTP_PORT),
		auth,
		config.SMTP_USERNAME,
		to,
		msg,
	)
}
