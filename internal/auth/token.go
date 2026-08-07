package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JerrySabor/ngts-warden/internal/config"
)

type Provider struct {
	HTTP   *http.Client
	Config config.Resolved
}

type cachedToken struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (p Provider) Token(ctx context.Context) (string, error) {
	if p.Config.AccessToken != "" {
		return p.Config.AccessToken, nil
	}
	if p.Config.ClientID == "" || p.Config.ClientSecret == "" || p.Config.TSGID == "" {
		return "", fmt.Errorf("missing client_id, client_secret, or tsg_id")
	}
	cachePath, err := config.CachePath(p.Config.Profile)
	if err == nil {
		if b, readErr := os.ReadFile(cachePath); readErr == nil {
			var cached cachedToken
			if json.Unmarshal(b, &cached) == nil && cached.AccessToken != "" && time.Until(cached.ExpiresAt) > time.Minute {
				return cached.AccessToken, nil
			}
		}
	}
	httpClient := p.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {"tsg_id:" + p.Config.TSGID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Config.AuthURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.SetBasicAuth(p.Config.ClientID, p.Config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request access token: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	var parsed tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", fmt.Errorf("decode access token response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || parsed.AccessToken == "" {
		return "", fmt.Errorf("authentication service returned HTTP %d", resp.StatusCode)
	}
	if cachePath != "" && parsed.ExpiresIn > 0 {
		_ = writeCache(cachePath, cachedToken{AccessToken: parsed.AccessToken, ExpiresAt: time.Now().Add(time.Duration(parsed.ExpiresIn) * time.Second)})
	}
	return parsed.AccessToken, nil
}

func ClearCache(profile string) error {
	path, err := config.CachePath(profile)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func writeCache(path string, token cachedToken) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".token-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer func() { _ = os.Remove(name) }()
	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}
