package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/JerrySabor/ngts-warden/internal/catalog"
	"github.com/JerrySabor/ngts-warden/internal/config"
	"github.com/JerrySabor/ngts-warden/internal/output"
	"github.com/spf13/cobra"
)

type App struct {
	version        string
	stdout, stderr io.Writer
	root           *cobra.Command
	opts           *rootOptions
}
type rootOptions struct {
	json                                                                  bool
	profile, accessToken, clientID, clientSecret, tsgID, baseURL, authURL string
	timeout                                                               time.Duration
}
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string { return e.Err.Error() }
func (e *ExitError) Unwrap() error { return e.Err }

func New(version string, stdout, stderr io.Writer) *App {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	a := &App{version: version, stdout: stdout, stderr: stderr, opts: &rootOptions{timeout: 60 * time.Second}}
	spec, err := catalog.Load()
	root := &cobra.Command{Use: "ngts-warden", Short: "Composable CLI and MCP client for Palo Alto NGTS", SilenceUsage: true, SilenceErrors: true}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().BoolVar(&a.opts.json, "json", false, "emit a stable JSON envelope")
	root.PersistentFlags().StringVar(&a.opts.profile, "profile", "", "named credential profile")
	root.PersistentFlags().StringVar(&a.opts.accessToken, "access-token", "", "one-off bearer token (prefer environment/config)")
	root.PersistentFlags().StringVar(&a.opts.clientID, "client-id", "", "one-off client ID")
	root.PersistentFlags().StringVar(&a.opts.clientSecret, "client-secret", "", "one-off client secret (may appear in shell history)")
	root.PersistentFlags().StringVar(&a.opts.tsgID, "tsg-id", "", "one-off tenant service group ID")
	root.PersistentFlags().StringVar(&a.opts.baseURL, "base-url", "", "override NGTS API base URL")
	root.PersistentFlags().StringVar(&a.opts.authURL, "auth-url", "", "override authentication URL")
	root.PersistentFlags().DurationVar(&a.opts.timeout, "timeout", 60*time.Second, "HTTP request timeout")
	root.AddCommand(a.versionCommand(), a.doctorCommand(), a.configCommand(), a.authCommand(), a.requestCommand(), a.operationsCommand(), a.mcpCommand())
	if err == nil {
		addGeneratedCommands(root, spec, a)
	}
	a.root = root
	return a
}

func (a *App) ExecuteContext(ctx context.Context) int {
	err := a.root.ExecuteContext(ctx)
	if err == nil {
		return 0
	}
	var exit *ExitError
	if errors.As(err, &exit) {
		return exit.Code
	}
	if !a.opts.json {
		_, _ = fmt.Fprintln(a.stderr, err)
	}
	return 2
}

// ExecuteArgs is useful for embedding and tests that need to invoke the CLI
// without replacing the process argument vector.
func (a *App) ExecuteArgs(ctx context.Context, args ...string) int {
	a.root.SetArgs(args)
	return a.ExecuteContext(ctx)
}

func (a *App) resolved() (config.Resolved, error) {
	return config.Resolve(config.Flags{Profile: a.opts.profile, AccessToken: a.opts.accessToken, ClientID: a.opts.clientID, ClientSecret: a.opts.clientSecret, TSGID: a.opts.tsgID, BaseURL: a.opts.baseURL, AuthURL: a.opts.authURL})
}
func (a *App) fail(op, kind string, code int, err error, status int, requestID string) error {
	if a.opts.json {
		_ = output.WriteJSON(a.stdout, output.Failure(op, kind, err.Error(), status, requestID))
	}
	return &ExitError{Code: code, Err: err}
}

func parseKeyValues(items []string) (map[string]string, error) {
	out := map[string]string{}
	for _, item := range items {
		k, v, ok := strings.Cut(item, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("expected key=value, got %q", item)
		}
		out[k] = v
	}
	return out, nil
}
