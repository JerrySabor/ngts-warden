package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadResolveProfile(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(root, "cache"))
	f := File{CurrentProfile: "prod", Profiles: map[string]Profile{"prod": {ClientID: "id", ClientSecret: "secret", TSGID: "123"}}}
	path, err := Save(f)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("config mode = %o, want 600", info.Mode().Perm())
	}
	loaded, _, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.CurrentProfile != "prod" {
		t.Fatalf("current profile = %q", loaded.CurrentProfile)
	}
	r, err := Resolve(Flags{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Profile != "prod" || r.ClientID != "id" || r.TSGID != "123" {
		t.Fatalf("unexpected resolved profile: %+v", r)
	}
}

func TestEnvironmentOverridesProfile(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	if _, err := Save(File{CurrentProfile: "prod", Profiles: map[string]Profile{"prod": {ClientID: "profile"}}}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("NGTS_WARDEN_CLIENT_ID", "environment")
	r, err := Resolve(Flags{})
	if err != nil {
		t.Fatal(err)
	}
	if r.ClientID != "environment" {
		t.Fatalf("client ID = %q", r.ClientID)
	}
}
