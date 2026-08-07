package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	DefaultBaseURL = "https://api.strata.paloaltonetworks.com/ngts"
	DefaultAuthURL = "https://auth.apps.paloaltonetworks.com/oauth2/access_token"
)

type Profile struct {
	ClientID     string `toml:"client_id,omitempty"`
	ClientSecret string `toml:"client_secret,omitempty"`
	TSGID        string `toml:"tsg_id,omitempty"`
	BaseURL      string `toml:"base_url,omitempty"`
	AuthURL      string `toml:"auth_url,omitempty"`
}

type File struct {
	CurrentProfile string             `toml:"current_profile,omitempty"`
	Profiles       map[string]Profile `toml:"profiles,omitempty"`
}

type Flags struct {
	Profile      string
	AccessToken  string
	ClientID     string
	ClientSecret string
	TSGID        string
	BaseURL      string
	AuthURL      string
}

type Resolved struct {
	Profile      string
	AccessToken  string
	ClientID     string
	ClientSecret string
	TSGID        string
	BaseURL      string
	AuthURL      string
	Source       string
}

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ngts-warden", "config.toml"), nil
}

func CachePath(profile string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	if profile == "" {
		profile = "default"
	}
	return filepath.Join(dir, "ngts-warden", "tokens", profile+".json"), nil
}

func Load() (File, string, error) {
	path, err := Path()
	if err != nil {
		return File{}, "", err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return File{Profiles: map[string]Profile{}}, path, nil
	}
	if err != nil {
		return File{}, path, fmt.Errorf("read config: %w", err)
	}
	var f File
	if err := toml.Unmarshal(b, &f); err != nil {
		return File{}, path, fmt.Errorf("parse config: %w", err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	return f, path, nil
}

func Save(f File) (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	b, err := toml.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("encode config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return "", fmt.Errorf("create config directory: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".config-*.tmp")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if err := tmp.Chmod(0600); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return "", fmt.Errorf("write config: %w", err)
	}
	return path, nil
}

func Resolve(flags Flags) (Resolved, error) {
	f, _, err := Load()
	if err != nil {
		return Resolved{}, err
	}
	profileName := first(flags.Profile, os.Getenv("NGTS_WARDEN_PROFILE"), f.CurrentProfile)
	profile := f.Profiles[profileName]
	r := Resolved{Profile: profileName, BaseURL: DefaultBaseURL, AuthURL: DefaultAuthURL, Source: "missing"}
	r.AccessToken, r.Source = choose(flags.AccessToken, os.Getenv("NGTS_WARDEN_ACCESS_TOKEN"), "flag", "env", profileName != "")
	r.ClientID, _ = choose(flags.ClientID, os.Getenv("NGTS_WARDEN_CLIENT_ID"), "flag", "env", false)
	r.ClientSecret, _ = choose(flags.ClientSecret, os.Getenv("NGTS_WARDEN_CLIENT_SECRET"), "flag", "env", false)
	r.TSGID, _ = choose(flags.TSGID, os.Getenv("NGTS_WARDEN_TSG_ID"), "flag", "env", false)
	r.BaseURL = first(flags.BaseURL, os.Getenv("NGTS_WARDEN_BASE_URL"), profile.BaseURL, DefaultBaseURL)
	r.AuthURL = first(flags.AuthURL, os.Getenv("NGTS_WARDEN_AUTH_URL"), profile.AuthURL, DefaultAuthURL)
	if r.AccessToken == "" {
		if profileName != "" {
			r.ClientID = first(flags.ClientID, os.Getenv("NGTS_WARDEN_CLIENT_ID"), profile.ClientID)
			r.ClientSecret = first(flags.ClientSecret, os.Getenv("NGTS_WARDEN_CLIENT_SECRET"), profile.ClientSecret)
			r.TSGID = first(flags.TSGID, os.Getenv("NGTS_WARDEN_TSG_ID"), profile.TSGID)
		}
		if r.ClientID != "" || r.ClientSecret != "" || r.TSGID != "" {
			r.Source = "credentials"
		}
	}
	if r.AccessToken == "" && r.ClientID == "" && r.ClientSecret == "" && r.TSGID == "" {
		r.Source = "missing"
	}
	return r, nil
}

func first(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func choose(flagValue, envValue, flagSource, envSource string, _ bool) (string, string) {
	if flagValue != "" {
		return flagValue, flagSource
	}
	if envValue != "" {
		return envValue, envSource
	}
	return "", "missing"
}
