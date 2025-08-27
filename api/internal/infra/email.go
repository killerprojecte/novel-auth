package infra

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/mail"
	"net/smtp"
)

type EmailClient interface {
	SendEmail(to string, title string, content string) error
}

type emailClient struct {
	email    string
	server   string
	password string
}

func NewEmailClient(email string, server string, password string) EmailClient {
	return &emailClient{
		email:    email,
		server:   server,
		password: password,
	}
}

func (c *emailClient) SendEmail(to string, title string, content string) error {
	from := mail.Address{Name: "轻小说机翻机器人", Address: c.email}
	target := mail.Address{Address: to}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = target.String()
	headers["Subject"] = title

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + content

	// Connect to the SMTP Server
	servername := c.server
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", c.email, c.password, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	// Auth
	if err = client.Auth(auth); err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	// To && From
	if err = client.Mail(from.Address); err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	if err = client.Rcpt(target.Address); err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	// Data
	w, err := client.Data()
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	err = w.Close()
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}

	err = client.Quit()
	if err != nil {
		slog.Error("Failed to send email", "error", err)
		return err
	}
	return nil
}
