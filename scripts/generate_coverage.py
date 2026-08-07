#!/usr/bin/env python3
"""Generate the checked-in NGTS operation coverage document from OpenAPI."""
import argparse
import json
import re
from pathlib import Path

GROUPS = {
    "Certificates": "certificates", "Certificate Import": "certificate-import",
    "TLS Server Endpoints": "tls-server-endpoints", "Certificate Discovery": "certificate-discovery",
    "Private Key Import": "private-key-import", "Certificate Request": "certificate-requests",
    "Certificate Policy": "certificate-policies", "Credential Management": "credentials",
    "Machine Installations": "machine-installations", "Machine Types": "machine-types",
    "Machines": "machines", "Event Logs": "event-logs", "VSatellite": "vsatellite",
    "Certificate Inventory Monitoring": "inventory-monitoring", "Certificate Auto-renewal Monitoring": "renewal-monitoring",
    "Certificate Tags": "certificate-tags", "Issuer Configurations": "issuer-configurations",
    "Issuer Sub CA Providers": "sub-ca-providers", "Workload Issuance Policies": "workload-policies",
    "Issuer Certificates": "issuer-certificates", "Certificate Approvals": "certificate-approvals",
    "Certificate Revocation Approvals": "revocation-approvals", "Plugins (Connectors)": "plugins",
    "Built-In Accounts": "built-in-accounts",
}

def kebab(value: str) -> str:
    value = re.sub(r"([a-z0-9])([A-Z])", r"\1-\2", value).replace("_", "-")
    return re.sub(r"[^a-zA-Z0-9-]+", "-", value).strip("-").lower()

def action(operation_id: str, group: str) -> str:
    result = kebab(operation_id)
    for prefix in (group + "-", group.removesuffix("s") + "-"):
        result = result.removeprefix(prefix)
    return result

def deref(value, root):
    if not isinstance(value, dict) or "$ref" not in value:
        return value
    current = root
    for part in value["$ref"].removeprefix("#/").split("/"):
        current = current[part]
    return current

def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--spec", default="api/openapi.json")
    parser.add_argument("--output", default="docs/api-coverage.md")
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    root = json.loads(Path(args.spec).read_text())
    operations = []
    for path, item in root["paths"].items():
        for method in ("get", "post", "put", "patch", "delete"):
            if method not in item:
                continue
            op = item[method]
            tag = (op.get("tags") or ["Other"])[0]
            group = GROUPS.get(tag, kebab(tag))
            parameters = [deref(p, root) for p in item.get("parameters", []) + op.get("parameters", [])]
            required = ", ".join(f"{p.get('in')}:{p.get('name')}" for p in parameters if p.get("required")) or "—"
            body = "required" if op.get("requestBody", {}).get("required") else ("optional" if op.get("requestBody") else "—")
            responses = set()
            for response in op.get("responses", {}).values():
                response = deref(response, root)
                responses.update(response.get("content", {}).keys())
            operations.append((group, action(op["operationId"], group), op["operationId"], method.upper(), path, op.get("summary", ""), required, body, ", ".join(sorted(responses)) or "—"))
    operations.sort()
    lines = ["# NGTS API coverage", "", f"Generated from `{args.spec}`. This snapshot contains **{len(operations)} operations** across **{len({row[0] for row in operations})} resource groups**.", "", "| CLI command | Method | Path | Operation ID | Required inputs | Body | Response content types | Summary |", "|---|---|---|---|---|---|---|---|"]
    for group, leaf, opid, method, path, summary, required, body, responses in operations:
        summary = summary.replace("|", "\\|").replace("\n", " ")
        lines.append(f"| `{group} {leaf}` | `{method}` | `{path}` | `{opid}` | `{required}` | `{body}` | `{responses}` | {summary} |")
    content = "\n".join(lines) + "\n"
    output = Path(args.output)
    if args.check:
        if not output.exists() or output.read_text() != content:
            print(f"generated coverage is stale: {output}")
            return 1
    else:
        output.parent.mkdir(parents=True, exist_ok=True)
        output.write_text(content)
    print(f"{len(operations)} operations, {len({row[0] for row in operations})} groups")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
