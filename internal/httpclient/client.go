package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JerrySabor/ngts-warden/internal/auth"
	"github.com/JerrySabor/ngts-warden/internal/config"
)

type Client struct {
	HTTP   *http.Client
	Config config.Resolved
}
type Request struct {
	Method  string
	Path    string
	Query   url.Values
	Headers map[string]string
	Body    []byte
	Accept  string
}
type Response struct {
	Status      int
	Header      http.Header
	Body        []byte
	ContentType string
	RequestID   string
}

func New(c config.Resolved) *Client {
	return &Client{HTTP: &http.Client{Timeout: 60 * time.Second}, Config: c}
}

func (c *Client) Do(ctx context.Context, in Request) (*Response, error) {
	base := strings.TrimRight(c.Config.BaseURL, "/")
	target := base + "/" + strings.TrimLeft(in.Path, "/")
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("parse API URL: %w", err)
	}
	if in.Query != nil {
		u.RawQuery = in.Query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(in.Method), u.String(), strings.NewReader(string(in.Body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", first(in.Accept, "application/json"))
	if len(in.Body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range in.Headers {
		req.Header.Set(k, v)
	}
	token, err := (auth.Provider{HTTP: c.HTTP, Config: c.Config}).Token(ctx)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Status: resp.StatusCode, Header: resp.Header, Body: body, ContentType: resp.Header.Get("Content-Type"), RequestID: first(resp.Header.Get("x-request-id"), resp.Header.Get("x-correlation-id"))}, nil
}

func Decode(body []byte, contentType string) any {
	if strings.Contains(strings.ToLower(contentType), "json") {
		var value any
		if json.Unmarshal(body, &value) == nil {
			return value
		}
	}
	return string(body)
}

func first(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
