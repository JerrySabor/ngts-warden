package cli

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/JerrySabor/ngts-warden/internal/httpclient"
	"github.com/JerrySabor/ngts-warden/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) requestCommand() *cobra.Command {
	root := &cobra.Command{Use: "request", Short: "use the authenticated raw NGTS HTTP escape hatch"}
	root.AddCommand(a.rawRequest("get", "GET", false))
	root.AddCommand(a.rawRequest("send", "", true))
	return root
}

func (a *App) rawRequest(name, fixedMethod string, methodArg bool) *cobra.Command {
	var queries, headers []string
	var body, bodyFile, out, accept string
	var execute bool
	use := name + " [METHOD] PATH"
	if !methodArg {
		use = name + " PATH"
	}
	cmd := &cobra.Command{Use: use, Short: "send a raw authenticated request", Args: cobra.RangeArgs(1, 2), RunE: func(cmd *cobra.Command, args []string) error {
		method := fixedMethod
		path := args[0]
		if methodArg {
			if len(args) != 2 {
				return a.fail("request.send", "usage", 2, fmt.Errorf("usage: request send METHOD PATH"), 0, "")
			}
			method = strings.ToUpper(args[0])
			path = args[1]
		}
		q, err := parseKeyValues(queries)
		if err != nil {
			return a.fail("request."+name, "usage", 2, err, 0, "")
		}
		query := url.Values{}
		for k, v := range q {
			query.Set(k, v)
		}
		h, err := parseKeyValues(headers)
		if err != nil {
			return a.fail("request."+name, "usage", 2, err, 0, "")
		}
		var bodyBytes []byte
		if body != "" && bodyFile != "" {
			return a.fail("request."+name, "usage", 2, fmt.Errorf("use only one of --body and --body-file"), 0, "")
		}
		if bodyFile != "" {
			if bodyFile == "-" {
				bodyBytes, err = readAll(cmd.InOrStdin())
			} else {
				bodyBytes, err = os.ReadFile(bodyFile)
			}
			if err != nil {
				return a.fail("request."+name, "file", 6, err, 0, "")
			}
		} else if body != "" {
			bodyBytes = []byte(body)
		}
		if method != "GET" && method != "HEAD" && !execute {
			preview := map[string]interface{}{"method": method, "path": path, "query": query, "headers": redactHeaders(h), "body": redactBody(bodyBytes), "execute_required": true}
			if a.opts.json {
				return output.WriteJSON(a.stdout, output.Success("request."+name, 0, preview, nil))
			}
			return output.WriteJSON(a.stdout, preview)
		}
		r, err := a.resolved()
		if err != nil {
			return a.fail("request."+name, "config", 2, err, 0, "")
		}
		client := httpclient.New(r)
		client.HTTP.Timeout = a.opts.timeout
		resp, err := client.Do(cmd.Context(), httpclient.Request{Method: method, Path: path, Query: query, Headers: h, Body: bodyBytes, Accept: accept})
		if err != nil {
			return a.fail("request."+name, "network", 4, err, 0, "")
		}
		if resp.Status < 200 || resp.Status >= 300 {
			return a.fail("request."+name, "api", 5, fmt.Errorf("API returned HTTP %d", resp.Status), resp.Status, resp.RequestID)
		}
		if out != "" {
			if err := os.WriteFile(out, resp.Body, 0600); err != nil {
				return a.fail("request."+name, "file", 6, err, resp.Status, resp.RequestID)
			}
			data := map[string]interface{}{"path": out, "bytes": len(resp.Body), "content_type": resp.ContentType}
			if a.opts.json {
				return output.WriteJSON(a.stdout, output.Success("request."+name, resp.Status, data, nil))
			}
			return output.WriteText(a.stdout, "wrote %s (%d bytes)\n", out, len(resp.Body))
		}
		data := httpclient.Decode(resp.Body, resp.ContentType)
		if a.opts.json {
			return output.WriteJSON(a.stdout, output.Success("request."+name, resp.Status, data, map[string]interface{}{"request_id": resp.RequestID, "content_type": resp.ContentType}))
		}
		if strings.Contains(strings.ToLower(resp.ContentType), "json") {
			return output.WriteJSON(a.stdout, data)
		}
		_, err = a.stdout.Write(resp.Body)
		return err
	}}
	cmd.Flags().StringArrayVarP(&queries, "query", "q", nil, "query parameter key=value (repeatable)")
	cmd.Flags().StringArrayVar(&headers, "header", nil, "request header key=value (repeatable)")
	cmd.Flags().StringVar(&body, "body", "", "request body")
	cmd.Flags().StringVar(&bodyFile, "body-file", "", "request body file, or - for stdin")
	cmd.Flags().StringVar(&out, "out", "", "response output file")
	cmd.Flags().StringVar(&accept, "accept", "", "Accept header")
	cmd.Flags().BoolVar(&execute, "execute", false, "send non-GET requests")
	return cmd
}

func readAll(r io.Reader) ([]byte, error) { return io.ReadAll(r) }
