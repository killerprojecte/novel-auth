package infra

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EmailClient interface {
	SendEmail(to string, title string, content string) error
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

func (c *emailClient) SendEmail(to string, title string, content string) error {
	formData := url.Values{
		"from":    {"轻小说机翻机器人 <postmaster@" + c.domain + ">"},
		"to":      {to},
		"subject": {title},
		"text":    {content},
	}
	body := strings.NewReader(formData.Encode())

	apiUrl := "https://api.eu.mailgun.net/v3/" + c.domain + "/messages"
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
