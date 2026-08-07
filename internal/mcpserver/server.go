package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/JerrySabor/ngts-warden/internal/catalog"
	"github.com/JerrySabor/ngts-warden/internal/config"
	"github.com/JerrySabor/ngts-warden/internal/httpclient"
	"github.com/JerrySabor/ngts-warden/internal/output"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config struct {
	Version    string
	Resolved   config.Resolved
	AllowWrite bool
	Stderr     interface{ Write([]byte) (int, error) }
}

type operationInput struct {
	OperationID string                 `json:"operation_id" jsonschema:"OpenAPI operationId"`
	Path        map[string]string      `json:"path,omitempty" jsonschema:"Path parameter values"`
	Query       map[string]string      `json:"query,omitempty" jsonschema:"Query parameter values"`
	Headers     map[string]string      `json:"headers,omitempty" jsonschema:"Additional request headers"`
	Body        map[string]interface{} `json:"body,omitempty" jsonschema:"JSON request body"`
	Confirm     bool                   `json:"confirm,omitempty" jsonschema:"Required for non-GET execution"`
}

type describeInput struct {
	OperationID string `json:"operation_id"`
}

func Run(ctx context.Context, c Config) error {
	spec, err := catalog.Load()
	if err != nil {
		return err
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "ngts-warden", Version: c.Version}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "doctor", Description: "Read-only credential and configuration status. Never returns secrets."}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		r := c.Resolved
		value := map[string]interface{}{"version": c.Version, "profile": r.Profile, "base_url": r.BaseURL, "auth_url": r.AuthURL, "auth_source": r.Source, "access_token_available": r.AccessToken != "", "client_credentials_available": r.ClientID != "" && r.ClientSecret != "" && r.TSGID != ""}
		return jsonResult(output.Success("doctor", 200, value, nil))
	})
	mcp.AddTool(server, &mcp.Tool{Name: "list_operations", Description: "List all 147 generated NGTS operations and their CLI commands."}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		return jsonResult(output.Success("operations.list", 200, spec.Operations(), map[string]interface{}{"count": len(spec.Operations())}))
	})
	mcp.AddTool(server, &mcp.Tool{Name: "describe_operation", Description: "Describe one OpenAPI operation, its parameters, body, and response content types."}, func(_ context.Context, _ *mcp.CallToolRequest, in describeInput) (*mcp.CallToolResult, any, error) {
		for _, op := range spec.Operations() {
			if op.ID == in.OperationID {
				return jsonResult(output.Success("operations.describe", 200, op, nil))
			}
		}
		return nil, nil, fmt.Errorf("operation %q not found", in.OperationID)
	})
	mcp.AddTool(server, &mcp.Tool{Name: "preview_operation", Description: "Build a redacted NGTS request without contacting the API."}, func(_ context.Context, _ *mcp.CallToolRequest, in operationInput) (*mcp.CallToolResult, any, error) {
		op, ok := findOperation(spec, in.OperationID)
		if !ok {
			return nil, nil, fmt.Errorf("operation %q not found", in.OperationID)
		}
		request, err := buildRequest(op, in)
		if err != nil {
			return nil, nil, err
		}
		preview := map[string]interface{}{"method": request.Method, "path": request.Path, "query": request.Query, "headers": request.Headers, "body": redactBody(request.Body), "execute_required": op.Method != "GET" && op.Method != "HEAD"}
		return jsonResult(output.Success(op.ID, 0, preview, nil))
	})
	mcp.AddTool(server, &mcp.Tool{Name: "execute_operation", Description: "Execute one NGTS operation. Non-GET operations require server --allow-write and confirm=true."}, func(ctx context.Context, _ *mcp.CallToolRequest, in operationInput) (*mcp.CallToolResult, any, error) {
		op, ok := findOperation(spec, in.OperationID)
		if !ok {
			return nil, nil, fmt.Errorf("operation %q not found", in.OperationID)
		}
		if op.Method != "GET" && op.Method != "HEAD" && (!c.AllowWrite || !in.Confirm) {
			return nil, nil, fmt.Errorf("non-GET operation requires mcp --allow-write and confirm=true")
		}
		req, err := buildRequest(op, in)
		if err != nil {
			return nil, nil, err
		}
		client := httpclient.New(c.Resolved)
		resp, err := client.Do(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		if resp.Status < 200 || resp.Status >= 300 {
			return jsonResult(output.Failure(op.ID, "api", fmt.Sprintf("API returned HTTP %d", resp.Status), resp.Status, resp.RequestID))
		}
		return jsonResult(output.Success(op.ID, resp.Status, httpclient.Decode(resp.Body, resp.ContentType), map[string]interface{}{"request_id": resp.RequestID, "content_type": resp.ContentType}))
	})
	return server.Run(ctx, &mcp.StdioTransport{})
}

func findOperation(spec *catalog.Spec, id string) (catalog.OperationInfo, bool) {
	for _, op := range spec.Operations() {
		if op.ID == id {
			return op, true
		}
	}
	return catalog.OperationInfo{}, false
}

func buildRequest(op catalog.OperationInfo, in operationInput) (httpclient.Request, error) {
	path := op.Path
	for key, value := range in.Path {
		path = strings.ReplaceAll(path, "{"+key+"}", url.PathEscape(value))
	}
	for _, p := range op.Parameters {
		if p.Required && p.In == "path" && in.Path[p.Name] == "" {
			return httpclient.Request{}, fmt.Errorf("missing path parameter %q", p.Name)
		}
	}
	query := url.Values{}
	for k, v := range in.Query {
		query.Set(k, v)
	}
	headers := map[string]string{}
	for k, v := range in.Headers {
		headers[k] = v
	}
	var body []byte
	if in.Body != nil {
		var err error
		body, err = json.Marshal(in.Body)
		if err != nil {
			return httpclient.Request{}, err
		}
	}
	return httpclient.Request{Method: op.Method, Path: path, Query: query, Headers: headers, Body: body}, nil
}

func jsonResult(v interface{}) (*mcp.CallToolResult, any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, v, nil
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

func redactValue(value interface{}) {
	switch current := value.(type) {
	case map[string]interface{}:
		for key, child := range current {
			lower := strings.ToLower(key)
			if strings.Contains(lower, "secret") || strings.Contains(lower, "password") || strings.Contains(lower, "token") || strings.Contains(lower, "privatekey") {
				current[key] = "[REDACTED]"
			} else {
				redactValue(child)
			}
		}
	case []interface{}:
		for _, child := range current {
			redactValue(child)
		}
	}
}
