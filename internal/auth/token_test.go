package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/JerrySabor/ngts-warden/internal/config"
)

func TestProviderRequestsClientCredentialsAndCaches(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", filepath.Join(root, "cache"))
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		user, password, ok := r.BasicAuth()
		if !ok || user != "client" || password != "secret" {
			t.Errorf("basic auth = %q/%q/%t", user, password, ok)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "tsg_id:123" {
			t.Errorf("form = %v", r.Form)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"access_token": "token", "expires_in": 900})
	}))
	defer server.Close()
	p := Provider{HTTP: server.Client(), Config: config.Resolved{Profile: "prod", ClientID: "client", ClientSecret: "secret", TSGID: "123", AuthURL: server.URL}}
	first, err := p.Token(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	second, err := p.Token(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if first != "token" || second != "token" {
		t.Fatalf("tokens = %q/%q", first, second)
	}
	if calls != 1 {
		t.Fatalf("token calls = %d, want 1", calls)
	}
}
