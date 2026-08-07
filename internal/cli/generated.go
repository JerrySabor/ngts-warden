package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JerrySabor/ngts-warden/internal/catalog"
	"github.com/JerrySabor/ngts-warden/internal/httpclient"
	"github.com/JerrySabor/ngts-warden/internal/output"
	"github.com/spf13/cobra"
)

func addGeneratedCommands(root *cobra.Command, spec *catalog.Spec, app *App) {
	groups := map[string]*cobra.Command{}
	for _, op := range spec.Operations() {
		group := groups[op.Group]
		if group == nil {
			group = &cobra.Command{Use: op.Group, Short: op.GroupName + " operations"}
			groups[op.Group] = group
			root.AddCommand(group)
		}
		use := op.Action
		if use == "" {
			use = strings.ToLower(op.Method)
		}
		if group.Commands() != nil {
			for _, existing := range group.Commands() {
				if existing.Name() == use {
					use += "-" + strings.ToLower(op.Method)
				}
			}
		}
		cmd := &cobra.Command{Use: use, Short: firstNonEmpty(op.Summary, op.Description, op.ID), Args: cobra.NoArgs}
		for _, p := range op.Parameters {
			addParameterFlag(cmd, p)
		}
		cmd.Flags().String("body", "", "inline JSON request body")
		cmd.Flags().String("body-file", "", "JSON request body file, or - for stdin")
		cmd.Flags().String("accept", "", "override the Accept header")
		cmd.Flags().String("out", "", "write non-JSON or binary response to this file")
		cmd.Flags().Bool("execute", false, "send non-GET requests; without this flag, print a redacted preview")
		captured := op
		cmd.RunE = func(cmd *cobra.Command, _ []string) error { return app.runOperation(cmd, captured) }
		group.AddCommand(cmd)
	}
}

func addParameterFlag(cmd *cobra.Command, p catalog.Parameter) {
	name := p.Name
	help := p.Description
	if p.Required {
		help = firstNonEmpty(help, "required") + " (required)"
	}
	switch p.Schema.Type {
	case "boolean":
		cmd.Flags().Bool(name, false, help)
	case "integer":
		cmd.Flags().Int64(name, 0, help)
	case "number":
		cmd.Flags().Float64(name, 0, help)
	case "array":
		cmd.Flags().StringArray(name, nil, help)
	default:
		cmd.Flags().String(name, "", help)
	}
}

func (a *App) runOperation(cmd *cobra.Command, op catalog.OperationInfo) error {
	resolved, err := a.resolved()
	if err != nil {
		return a.fail(op.ID, "config", 2, err, 0, "")
	}
	values := map[string]string{}
	path := op.Path
	query := url.Values{}
	headers := map[string]string{}
	for _, p := range op.Parameters {
		value, changed, err := parameterValue(cmd, p)
		if err != nil {
			return a.fail(op.ID, "usage", 2, err, 0, "")
		}
		if p.Required && !changed && value == "" {
			return a.fail(op.ID, "usage", 2, fmt.Errorf("missing required --%s", p.Name), 0, "")
		}
		if value == "" && !changed {
			continue
		}
		switch p.In {
		case "path":
			values[p.Name] = url.PathEscape(value)
		case "query":
			query.Set(p.Name, value)
		case "header":
			headers[p.Name] = value
		}
	}
	for key, value := range values {
		path = strings.ReplaceAll(path, "{"+key+"}", value)
	}
	body, err := requestBody(cmd, op.RequestBody)
	if err != nil {
		return a.fail(op.ID, "usage", 2, err, 0, "")
	}
	execute, _ := cmd.Flags().GetBool("execute")
	if op.Method != "GET" && op.Method != "HEAD" && !execute {
		preview := map[string]interface{}{"method": op.Method, "path": path, "query": query, "headers": redactHeaders(headers), "body": redactBody(body), "execute_required": true}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success(op.ID, 0, preview, nil))
		}
		return output.WriteJSON(a.stdout, preview)
	}
	accept, _ := cmd.Flags().GetString("accept")
	outPath, _ := cmd.Flags().GetString("out")
	client := httpclient.New(resolved)
	client.HTTP.Timeout = a.opts.timeout
	resp, err := client.Do(cmd.Context(), httpclient.Request{Method: op.Method, Path: path, Query: query, Headers: headers, Body: body, Accept: accept})
	if err != nil {
		return a.fail(op.ID, "network", 4, err, 0, "")
	}
	if resp.Status < 200 || resp.Status >= 300 {
		return a.fail(op.ID, "api", 5, fmt.Errorf("API returned HTTP %d", resp.Status), resp.Status, resp.RequestID)
	}
	if outPath != "" {
		if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
			return a.fail(op.ID, "file", 6, err, resp.Status, resp.RequestID)
		}
		if err := os.WriteFile(outPath, resp.Body, 0600); err != nil {
			return a.fail(op.ID, "file", 6, err, resp.Status, resp.RequestID)
		}
		data := map[string]interface{}{"path": outPath, "bytes": len(resp.Body), "content_type": resp.ContentType}
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success(op.ID, resp.Status, data, map[string]interface{}{"request_id": resp.RequestID}))
		}
		return output.WriteText(a.stdout, "wrote %s (%d bytes)\n", outPath, len(resp.Body))
	}
	data := httpclient.Decode(resp.Body, resp.ContentType)
	if a.opts.json {
		return output.WriteJSON(a.stdout, output.Success(op.ID, resp.Status, data, map[string]interface{}{"request_id": resp.RequestID, "content_type": resp.ContentType}))
	}
	if strings.Contains(strings.ToLower(resp.ContentType), "json") {
		return output.WriteJSON(a.stdout, data)
	}
	_, err = a.stdout.Write(resp.Body)
	if err == nil && len(resp.Body) > 0 && resp.Body[len(resp.Body)-1] != '\n' {
		_, err = a.stdout.Write([]byte("\n"))
	}
	return err
}

func parameterValue(cmd *cobra.Command, p catalog.Parameter) (string, bool, error) {
	changed := cmd.Flags().Changed(p.Name)
	switch p.Schema.Type {
	case "boolean":
		v, e := cmd.Flags().GetBool(p.Name)
		if !changed {
			return "", false, e
		}
		return strconv.FormatBool(v), changed, e
	case "integer":
		v, e := cmd.Flags().GetInt64(p.Name)
		if !changed {
			return "", false, e
		}
		return strconv.FormatInt(v, 10), changed, e
	case "number":
		v, e := cmd.Flags().GetFloat64(p.Name)
		if !changed {
			return "", false, e
		}
		return strconv.FormatFloat(v, 'f', -1, 64), changed, e
	case "array":
		v, e := cmd.Flags().GetStringArray(p.Name)
		return strings.Join(v, ","), changed, e
	default:
		v, e := cmd.Flags().GetString(p.Name)
		return v, changed, e
	}
}

func requestBody(cmd *cobra.Command, bodySpec *catalog.RequestBody) ([]byte, error) {
	body, _ := cmd.Flags().GetString("body")
	file, _ := cmd.Flags().GetString("body-file")
	if body != "" && file != "" {
		return nil, fmt.Errorf("use only one of --body and --body-file")
	}
	if file != "" {
		if file == "-" {
			return ioReadAll(cmd.InOrStdin())
		}
		return os.ReadFile(file)
	}
	if body != "" {
		if !json.Valid([]byte(body)) {
			return nil, fmt.Errorf("--body must contain valid JSON")
		}
		return []byte(body), nil
	}
	if bodySpec != nil && bodySpec.Required {
		return nil, fmt.Errorf("request body required; provide --body or --body-file")
	}
	return nil, nil
}

func ioReadAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }

func redactHeaders(headers map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range headers {
		if strings.Contains(strings.ToLower(k), "authorization") || strings.Contains(strings.ToLower(k), "secret") {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func redactBody(body []byte) interface{} {
	if len(body) == 0 {
		return nil
	}
	var value interface{}
	if json.Unmarshal(body, &value) != nil {
		return "[JSON BODY REDACTED]"
	}
	redactValue(value)
	return value
}
func redactValue(v interface{}) {
	switch x := v.(type) {
	case map[string]interface{}:
		for k, val := range x {
			low := strings.ToLower(k)
			if strings.Contains(low, "secret") || strings.Contains(low, "password") || strings.Contains(low, "token") || strings.Contains(low, "privatekey") {
				x[k] = "[REDACTED]"
			} else {
				redactValue(val)
			}
		}
	case []interface{}:
		for _, item := range x {
			redactValue(item)
		}
	}
}
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
