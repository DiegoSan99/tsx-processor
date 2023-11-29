package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"

	"github.com/DiegoSan99/transaction-processor/src/config"
)

type SmtpClient struct {
	Host     string
	Port     string
	Username string
	Password string
	Auth     smtp.Auth
}

func NewSmtpClient(cfg *config.AppConfig) *SmtpClient {
	auth := smtp.PlainAuth("", cfg.EmailUsername, cfg.EmailPassword, cfg.EmailHost)
	return &SmtpClient{
		Host:     cfg.EmailHost,
		Port:     cfg.EmailPort,
		Username: cfg.EmailUsername,
		Password: cfg.EmailPassword,
		Auth:     auth,
	}
}
func (c *SmtpClient) SendEmail(recipient, subject, body, headerImagePath string) error {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// Headers
	headers := make(map[string]string)
	headers["From"] = c.Username
	headers["To"] = recipient
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary())
	for k, v := range headers {
		fmt.Fprintf(buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(buf, "\r\n")

	// Body
	htmlPartHeader := make(textproto.MIMEHeader)
	htmlPartHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlPartHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	htmlPart, err := writer.CreatePart(htmlPartHeader)
	if err != nil {
		return err
	}
	_, err = htmlPart.Write([]byte(body))
	if err != nil {
		return err
	}

	// Header
	imagePartHeader := make(textproto.MIMEHeader)
	imagePartHeader.Set("Content-Type", "image/png")
	imagePartHeader.Set("Content-ID", "<headerImage>")
	imagePartHeader.Set("Content-Disposition", "inline")
	imagePartHeader.Set("Content-Transfer-Encoding", "base64")
	imagePart, err := writer.CreatePart(imagePartHeader)
	if err != nil {
		return err
	}
	imageData, err := os.ReadFile(headerImagePath)
	if err != nil {
		return err
	}
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	_, err = imagePart.Write([]byte(base64Image))
	if err != nil {
		return err
	}

	writer.Close()

	serverAddr := c.Host + ":" + c.Port
	err = smtp.SendMail(serverAddr, c.Auth, c.Username, []string{recipient}, buf.Bytes())
	return err
}
