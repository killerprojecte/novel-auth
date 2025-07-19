package infra

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EmailClient interface {
	SendVerifyEmail(email string, code string) error
}

type emailClient struct {
	domain string
	apiKey string
	client http.Client
}

func NewEmailClient(domain string, apiKey string) EmailClient {
	return &emailClient{
		domain: domain,
		apiKey: apiKey,
		client: http.Client{Timeout: 5 * time.Second},
	}
}

func (c *emailClient) SendVerifyEmail(email string, code string) error {
	return c.sendEmail(
		email,
		fmt.Sprintf(
			"%s 轻小说机翻机器人 注册激活码",
			code,
		),
		fmt.Sprintf(
			"您的注册激活码为 %s\n"+
				"激活码将会在15分钟后失效,请尽快完成注册\n"+
				"这是系统邮件，请勿回复",
			code,
		),
	)
}

func (c *emailClient) sendEmail(to string, subject string, text string) error {
	formData := url.Values{
		"from":    {"轻小说机翻机器人 <postmaster@" + c.domain + ">"},
		"to":      {to},
		"subject": {subject},
		"text":    {text},
	}
	body := strings.NewReader(formData.Encode())

	apiUrl := "https://api.en.mailgun.net/v3/" + c.domain + "/messages"
	req, err := http.NewRequest(http.MethodPost, apiUrl, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("api", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return &url.Error{
			Op:  "POST",
			URL: apiUrl,
			Err: http.ErrServerClosed,
		}
	}

	return nil
}
