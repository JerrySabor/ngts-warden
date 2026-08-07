package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/JerrySabor/ngts-warden/internal/auth"
	"github.com/JerrySabor/ngts-warden/internal/catalog"
	"github.com/JerrySabor/ngts-warden/internal/config"
	"github.com/JerrySabor/ngts-warden/internal/mcpserver"
	"github.com/JerrySabor/ngts-warden/internal/output"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func (a *App) versionCommand() *cobra.Command {
	return &cobra.Command{Use: "version", Short: "print the installed version", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("version", 200, map[string]string{"version": a.version}, nil))
		}
		return output.WriteText(a.stdout, "ngts-warden %s\n", a.version)
	}}
}

func (a *App) doctorCommand() *cobra.Command {
	var online, strict bool
	cmd := &cobra.Command{Use: "doctor", Short: "check local configuration and optional authentication", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		r, path, err := resolveWithPath(a)
		if err != nil {
			return a.fail("doctor", "config", 2, err, 0, "")
		}
		report := map[string]interface{}{"version": a.version, "config_path": path, "profile": r.Profile, "base_url": r.BaseURL, "auth_url": r.AuthURL, "auth_source": r.Source, "access_token_available": r.AccessToken != "", "client_credentials_available": r.ClientID != "" && r.ClientSecret != "" && r.TSGID != "", "ready": r.AccessToken != "" || (r.ClientID != "" && r.ClientSecret != "" && r.TSGID != ""), "online_checked": online}
		if online && report["ready"] == true {
			_, err = auth.Provider{Config: r}.Token(cmd.Context())
			if err != nil {
				report["online_error"] = err.Error()
				report["ready"] = false
			} else {
				report["online_ok"] = true
			}
		}
		if a.opts.json {
			if err := output.WriteJSON(a.stdout, output.Success("doctor", 200, report, nil)); err != nil {
				return err
			}
		} else {
			for _, key := range []string{"version", "config_path", "profile", "base_url", "auth_source", "ready", "online_checked"} {
				_, _ = fmt.Fprintf(a.stdout, "%-26s %v\n", key, report[key])
			}
		}
		if strict && report["ready"] != true {
			return &ExitError{Code: 3, Err: fmt.Errorf("NGTS authentication is not configured or failed")}
		}
		return nil
	}}
	cmd.Flags().BoolVar(&online, "online", false, "perform an authentication-service check")
	cmd.Flags().BoolVar(&strict, "strict", false, "return an error when authentication is not ready")
	return cmd
}

func (a *App) configCommand() *cobra.Command {
	root := &cobra.Command{Use: "config", Short: "manage named NGTS credential profiles"}
	var name, clientID, secret, tsgID, baseURL, authURL string
	initCmd := &cobra.Command{Use: "init NAME", Short: "create or replace a profile", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		name = args[0]
		f, _, err := config.Load()
		if err != nil {
			return a.fail("config.init", "config", 2, err, 0, "")
		}
		if f.Profiles == nil {
			f.Profiles = map[string]config.Profile{}
		}
		if secret == "" && os.Getenv("NGTS_WARDEN_CLIENT_SECRET") == "" && term.IsTerminal(int(os.Stdin.Fd())) {
			b, e := term.ReadPassword(int(os.Stdin.Fd()))
			_, _ = fmt.Fprintln(a.stderr)
			if e != nil {
				return a.fail("config.init", "config", 2, e, 0, "")
			}
			secret = string(b)
		}
		f.Profiles[name] = config.Profile{ClientID: clientID, ClientSecret: secret, TSGID: tsgID, BaseURL: baseURL, AuthURL: authURL}
		if f.CurrentProfile == "" {
			f.CurrentProfile = name
		}
		saved, err := config.Save(f)
		if err != nil {
			return a.fail("config.init", "config", 2, err, 0, "")
		}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("config.init", 200, map[string]string{"profile": name, "path": saved}, nil))
		}
		return output.WriteText(a.stdout, "saved profile %s in %s\n", name, saved)
	}}
	initCmd.Flags().StringVar(&clientID, "client-id", "", "service account client ID")
	initCmd.Flags().StringVar(&secret, "client-secret", "", "service account client secret (prefer the hidden prompt)")
	initCmd.Flags().StringVar(&tsgID, "tsg-id", "", "tenant service group ID")
	initCmd.Flags().StringVar(&baseURL, "base-url", config.DefaultBaseURL, "NGTS API base URL")
	initCmd.Flags().StringVar(&authURL, "auth-url", config.DefaultAuthURL, "OAuth token URL")
	listCmd := &cobra.Command{Use: "list", Short: "list configured profile names", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		f, path, err := config.Load()
		if err != nil {
			return a.fail("config.list", "config", 2, err, 0, "")
		}
		names := make([]string, 0, len(f.Profiles))
		for n := range f.Profiles {
			names = append(names, n)
		}
		sort.Strings(names)
		data := map[string]interface{}{"current_profile": f.CurrentProfile, "profiles": names, "path": path}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("config.list", 200, data, nil))
		}
		for _, n := range names {
			marker := " "
			if n == f.CurrentProfile {
				marker = "*"
			}
			_, _ = fmt.Fprintf(a.stdout, "%s %s\n", marker, n)
		}
		return nil
	}}
	showCmd := &cobra.Command{Use: "show [NAME]", Short: "show a profile with secrets redacted", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		f, path, err := config.Load()
		if err != nil {
			return a.fail("config.show", "config", 2, err, 0, "")
		}
		n := f.CurrentProfile
		if len(args) == 1 {
			n = args[0]
		}
		p, ok := f.Profiles[n]
		if !ok {
			return a.fail("config.show", "config", 2, fmt.Errorf("profile %q not found", n), 0, "")
		}
		data := map[string]interface{}{"name": n, "path": path, "client_id": p.ClientID, "client_secret": redact(p.ClientSecret), "tsg_id": p.TSGID, "base_url": firstNonEmpty(p.BaseURL, config.DefaultBaseURL), "auth_url": firstNonEmpty(p.AuthURL, config.DefaultAuthURL)}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("config.show", 200, data, nil))
		}
		return output.WriteJSON(a.stdout, data)
	}}
	useCmd := &cobra.Command{Use: "use NAME", Short: "select the default profile", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		f, _, err := config.Load()
		if err != nil {
			return a.fail("config.use", "config", 2, err, 0, "")
		}
		if _, ok := f.Profiles[args[0]]; !ok {
			return a.fail("config.use", "config", 2, fmt.Errorf("profile %q not found", args[0]), 0, "")
		}
		f.CurrentProfile = args[0]
		path, err := config.Save(f)
		if err != nil {
			return a.fail("config.use", "config", 2, err, 0, "")
		}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("config.use", 200, map[string]string{"current_profile": args[0], "path": path}, nil))
		}
		return output.WriteText(a.stdout, "current profile: %s\n", args[0])
	}}
	deleteCmd := &cobra.Command{Use: "delete NAME", Short: "delete a local profile", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		f, _, err := config.Load()
		if err != nil {
			return a.fail("config.delete", "config", 2, err, 0, "")
		}
		if _, ok := f.Profiles[args[0]]; !ok {
			return a.fail("config.delete", "config", 2, fmt.Errorf("profile %q not found", args[0]), 0, "")
		}
		delete(f.Profiles, args[0])
		if f.CurrentProfile == args[0] {
			f.CurrentProfile = ""
		}
		_, err = config.Save(f)
		if err != nil {
			return a.fail("config.delete", "config", 2, err, 0, "")
		}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("config.delete", 200, map[string]string{"deleted": args[0]}, nil))
		}
		return output.WriteText(a.stdout, "deleted profile: %s\n", args[0])
	}}
	root.AddCommand(initCmd, listCmd, showCmd, useCmd, deleteCmd)
	return root
}

func (a *App) authCommand() *cobra.Command {
	root := &cobra.Command{Use: "auth", Short: "inspect or clear authentication state"}
	status := &cobra.Command{Use: "status", Short: "show token and credential availability without secrets", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		r, _, err := resolveWithPath(a)
		if err != nil {
			return a.fail("auth.status", "config", 2, err, 0, "")
		}
		data := map[string]interface{}{"profile": r.Profile, "source": r.Source, "access_token_available": r.AccessToken != "", "client_credentials_available": r.ClientID != "" && r.ClientSecret != "" && r.TSGID != ""}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("auth.status", 200, data, nil))
		}
		return output.WriteJSON(a.stdout, data)
	}}
	clear := &cobra.Command{Use: "clear-cache", Short: "remove the cached token for the selected profile", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		r, _, err := resolveWithPath(a)
		if err != nil {
			return a.fail("auth.clear-cache", "config", 2, err, 0, "")
		}
		if err := auth.ClearCache(r.Profile); err != nil {
			return a.fail("auth.clear-cache", "file", 6, err, 0, "")
		}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("auth.clear-cache", 200, map[string]string{"profile": r.Profile}, nil))
		}
		return output.WriteText(a.stdout, "cleared token cache for %s\n", r.Profile)
	}}
	root.AddCommand(status, clear)
	return root
}

func (a *App) operationsCommand() *cobra.Command {
	root := &cobra.Command{Use: "operations", Short: "discover the pinned NGTS OpenAPI operation catalog"}
	list := &cobra.Command{Use: "list", Short: "list all generated operation commands", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		s, err := catalog.Load()
		if err != nil {
			return err
		}
		ops := s.Operations()
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("operations.list", 200, ops, map[string]interface{}{"count": len(ops)}))
		}
		for _, op := range ops {
			_, _ = fmt.Fprintf(a.stdout, "%-28s %-7s %-34s %s\n", op.Group+"/"+op.Action, op.Method, op.Path, op.ID)
		}
		return nil
	}}
	describe := &cobra.Command{Use: "describe OPERATION_ID", Short: "show parameters and response types for one operation", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		s, err := catalog.Load()
		if err != nil {
			return err
		}
		for _, op := range s.Operations() {
			if op.ID == args[0] {
				if a.opts.json {
					return output.WriteJSON(a.stdout, output.Success("operations.describe", 200, op, nil))
				}
				return output.WriteJSON(a.stdout, op)
			}
		}
		return a.fail("operations.describe", "usage", 2, fmt.Errorf("operation %q not found", args[0]), 0, "")
	}}
	root.AddCommand(list, describe)
	return root
}

func (a *App) mcpCommand() *cobra.Command {
	var allowWrite bool
	cmd := &cobra.Command{Use: "mcp", Short: "run the MCP stdio server", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		r, _, err := resolveWithPath(a)
		if err != nil {
			return err
		}
		return mcpserver.Run(cmd.Context(), mcpserver.Config{Version: a.version, Resolved: r, AllowWrite: allowWrite, Stderr: a.stderr})
	}}
	cmd.Flags().BoolVar(&allowWrite, "allow-write", false, "allow MCP execute_operation to send non-GET requests")
	return cmd
}

func resolveWithPath(a *App) (config.Resolved, string, error) {
	r, err := a.resolved()
	path, _ := config.Path()
	return r, path, err
}
func redact(s string) string {
	if s == "" {
		return ""
	}
	return "[REDACTED]"
}
