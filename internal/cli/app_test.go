package cli

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoctorWithoutCredentialsIsMachineReadable(t *testing.T) {
	var out, errOut bytes.Buffer
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	app := New("test", &out, &errOut)
	if code := app.ExecuteArgs(t.Context(), "--json", "doctor"); code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), `"ready": false`) || !strings.Contains(out.String(), `"auth_source": "missing"`) {
		t.Fatalf("doctor output = %s", out.String())
	}
}

func TestOperationListIncludesCompleteCatalog(t *testing.T) {
	var out, errOut bytes.Buffer
	app := New("test", &out, &errOut)
	if code := app.ExecuteArgs(t.Context(), "operations", "list"); code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, errOut.String())
	}
	if lines := strings.Count(out.String(), "\n"); lines != 147 {
		t.Fatalf("operation lines = %d", lines)
	}
}

func TestGeneratedReadUsesAuthAndOmitsUnchangedOptionalFlags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Query().Get("excludeSupersededInstances") != "" {
			t.Errorf("unexpected default boolean query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()
	var out, errOut bytes.Buffer
	app := New("test", &out, &errOut)
	if code := app.ExecuteArgs(t.Context(), "--json", "--base-url", server.URL, "--access-token", "test-token", "certificates", "get-all"); code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), `"ok": true`) {
		t.Fatalf("output = %s", out.String())
	}
}

func TestGeneratedWriteDefaultsToPreview(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true }))
	defer server.Close()
	var out, errOut bytes.Buffer
	app := New("test", &out, &errOut)
	if code := app.ExecuteArgs(t.Context(), "--json", "--base-url", server.URL, "--access-token", "test-token", "certificate-import", "certificateimports-create", "--body", `{"certificates":[]}`); code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, errOut.String())
	}
	if called {
		t.Fatal("preview contacted the API")
	}
	if !strings.Contains(out.String(), `"execute_required": true`) {
		t.Fatalf("preview = %s", out.String())
	}
}
