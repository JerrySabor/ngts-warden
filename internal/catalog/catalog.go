package catalog

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

//go:embed spec.json
var specFS embed.FS

type Spec struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Tags       []Tag               `json:"tags"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type Server struct {
	URL string `json:"url"`
}
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Components struct {
	Parameters map[string]Parameter `json:"parameters"`
}

type PathItem struct {
	Summary    string      `json:"summary"`
	Parameters []Parameter `json:"parameters"`
	Get        *Operation  `json:"get"`
	Post       *Operation  `json:"post"`
	Put        *Operation  `json:"put"`
	Patch      *Operation  `json:"patch"`
	Delete     *Operation  `json:"delete"`
}

type Operation struct {
	OperationID string              `json:"operationId"`
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	Parameters  []Parameter         `json:"parameters"`
	RequestBody *RequestBody        `json:"requestBody"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Ref         string `json:"$ref"`
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Schema      Schema `json:"schema"`
}

type Schema struct {
	Ref         string        `json:"$ref"`
	Type        string        `json:"type"`
	Format      string        `json:"format"`
	Description string        `json:"description"`
	Enum        []interface{} `json:"enum"`
	Items       *Schema       `json:"items"`
}

type RequestBody struct {
	Description string               `json:"description"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type OperationInfo struct {
	ID            string
	Method        string
	Path          string
	Group         string
	GroupName     string
	Action        string
	Summary       string
	Description   string
	Parameters    []Parameter
	RequestBody   *RequestBody
	ResponseTypes []string
}

var groupSlugs = map[string]string{
	"Certificates":                        "certificates",
	"Certificate Import":                  "certificate-import",
	"TLS Server Endpoints":                "tls-server-endpoints",
	"Certificate Discovery":               "certificate-discovery",
	"Private Key Import":                  "private-key-import",
	"Certificate Request":                 "certificate-requests",
	"Certificate Policy":                  "certificate-policies",
	"Credential Management":               "credentials",
	"Machine Installations":               "machine-installations",
	"Machine Types":                       "machine-types",
	"Machines":                            "machines",
	"Event Logs":                          "event-logs",
	"VSatellite":                          "vsatellite",
	"Certificate Inventory Monitoring":    "inventory-monitoring",
	"Certificate Auto-renewal Monitoring": "renewal-monitoring",
	"Certificate Tags":                    "certificate-tags",
	"Issuer Configurations":               "issuer-configurations",
	"Issuer Sub CA Providers":             "sub-ca-providers",
	"Workload Issuance Policies":          "workload-policies",
	"Issuer Certificates":                 "issuer-certificates",
	"Certificate Approvals":               "certificate-approvals",
	"Certificate Revocation Approvals":    "revocation-approvals",
	"Plugins (Connectors)":                "plugins",
	"Built-In Accounts":                   "built-in-accounts",
}

func Load() (*Spec, error) {
	b, err := specFS.ReadFile("spec.json")
	if err != nil {
		return nil, fmt.Errorf("read embedded OpenAPI spec: %w", err)
	}
	var s Spec
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("parse embedded OpenAPI spec: %w", err)
	}
	return &s, nil
}

func (s *Spec) Operations() []OperationInfo {
	var out []OperationInfo
	for path, item := range s.Paths {
		for method, op := range map[string]*Operation{"GET": item.Get, "POST": item.Post, "PUT": item.Put, "PATCH": item.Patch, "DELETE": item.Delete} {
			if op == nil {
				continue
			}
			groupName := "Other"
			if len(op.Tags) > 0 {
				groupName = op.Tags[0]
			}
			group := groupSlugs[groupName]
			if group == "" {
				group = kebab(groupName)
			}
			params := make([]Parameter, 0, len(item.Parameters)+len(op.Parameters))
			for _, p := range item.Parameters {
				params = append(params, s.resolveParameter(p))
			}
			for _, p := range op.Parameters {
				params = append(params, s.resolveParameter(p))
			}
			seen := map[string]bool{}
			unique := params[:0]
			for _, p := range params {
				key := p.In + ":" + p.Name
				if !seen[key] {
					seen[key] = true
					unique = append(unique, p)
				}
			}
			params = unique
			respTypes := []string{}
			for _, resp := range op.Responses {
				for ct := range resp.Content {
					respTypes = append(respTypes, ct)
				}
			}
			sort.Strings(respTypes)
			out = append(out, OperationInfo{ID: op.OperationID, Method: method, Path: path, Group: group, GroupName: groupName, Action: actionName(op.OperationID, group), Summary: op.Summary, Description: op.Description, Parameters: params, RequestBody: op.RequestBody, ResponseTypes: respTypes})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Group != out[j].Group {
			return out[i].Group < out[j].Group
		}
		if out[i].Action != out[j].Action {
			return out[i].Action < out[j].Action
		}
		return out[i].Method < out[j].Method
	})
	return out
}

func (s *Spec) resolveParameter(p Parameter) Parameter {
	if p.Ref == "" {
		return p
	}
	name := strings.TrimPrefix(p.Ref, "#/components/parameters/")
	if resolved, ok := s.Components.Parameters[name]; ok {
		return resolved
	}
	return p
}

func actionName(id, group string) string {
	a := kebab(id)
	for _, prefix := range []string{group + "-", strings.TrimSuffix(group, "s") + "-"} {
		a = strings.TrimPrefix(a, prefix)
	}
	return a
}

func kebab(s string) string {
	var b bytes.Buffer
	for i, r := range s {
		if r == '_' || r == ' ' || r == '/' {
			if b.Len() > 0 && b.Bytes()[b.Len()-1] != '-' {
				b.WriteByte('-')
			}
			continue
		}
		if r >= 'A' && r <= 'Z' {
			if i > 0 && b.Len() > 0 && b.Bytes()[b.Len()-1] != '-' {
				b.WriteByte('-')
			}
			b.WriteByte(byte(r - 'A' + 'a'))
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return strings.Trim(b.String(), "-")
}
