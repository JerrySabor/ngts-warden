package httpclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/JerrySabor/ngts-warden/internal/config"
)

func TestDoBuildsAuthenticatedRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/ngts/v1/machines" || r.URL.Query().Get("limit") != "3" {
			t.Errorf("request = %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"one"}]}`))
	}))
	defer server.Close()
	c := New(config.Resolved{BaseURL: server.URL + "/ngts", AccessToken: "token"})
	resp, err := c.Do(t.Context(), Request{Method: http.MethodGet, Path: "/v1/machines", Query: mapValues("limit", "3")})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != http.StatusOK || resp.ContentType != "application/json" {
		t.Fatalf("response = %+v", resp)
	}
}

func mapValues(k, v string) url.Values { return url.Values{k: []string{v}} }
