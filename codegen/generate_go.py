#!/usr/bin/env python3
"""Emit models_gen.go from codegen/models.json.

Usage (from the repo root):
    python3 codegen/normalize_spec.py spec/openapi.json > codegen/models.json
    python3 codegen/generate_go.py codegen/models.json > models_gen.go
    gofmt -w models_gen.go

`go generate ./...` runs all three steps (see generate.go).
"""
import json
import re
import sys

INITIALISMS = {
    "id": "ID", "ids": "IDs", "url": "URL", "urls": "URLs", "uri": "URI", "api": "API",
    "http": "HTTP", "https": "HTTPS", "json": "JSON", "uuid": "UUID", "sql": "SQL",
    "ip": "IP", "ttl": "TTL", "html": "HTML", "mcp": "MCP", "jwt": "JWT", "sso": "SSO",
    "oauth": "OAuth", "otlp": "OTLP", "s3": "S3", "hmac": "HMAC", "ui": "UI",
}

# Names the hand-written part of the package already uses.
RESERVED = {"Client", "Error", "Record", "Params", "Page", "Iterator", "Option",
            "RequestOption", "RetryInfo", "RateLimitInfo"}


def words(name):
    name = name.replace("$", " ")
    parts = re.split(r"[^A-Za-z0-9]+", name)
    out = []
    for p in parts:
        # split camelCase / PascalCase, keep runs of capitals together
        out += re.findall(r"[A-Z]+(?=[A-Z][a-z])|[A-Z]?[a-z]+|[A-Z]+|[0-9]+", p)
    return [w for w in out if w]


def go_field_name(name):
    ws = words(name)
    if not ws:
        return "Field"
    out = "".join(INITIALISMS.get(w.lower(), w[0].upper() + w[1:]) for w in ws)
    if out[0].isdigit():
        out = "N" + out
    return out


def type_name(name):
    return name + "Model" if name in RESERVED else name


def comment(text, indent=""):
    text = (text or "").strip()
    if not text:
        return ""
    para = text.split("\n\n")[0]
    lines = []
    for line in para.splitlines():
        line = line.rstrip()
        lines.append(f"{indent}// {line}" if line else f"{indent}//")
    return "\n".join(lines) + "\n"


class Gen:
    def __init__(self, ir):
        self.ir = ir
        self.models = {m["name"]: m for m in ir["models"]}

    def go_type(self, t, optional):
        k = t["kind"]
        if k == "string":
            base = "string"
        elif k == "integer":
            base = "int64"
        elif k == "number":
            base = "float64"
        elif k == "boolean":
            base = "bool"
        elif k == "any":
            return "any"
        elif k == "array":
            return "[]" + self.go_type(t["items"], False)
        elif k == "map":
            return "map[string]" + self.go_type(t["values"], False)
        elif k == "ref":
            target = self.models.get(t["name"])
            base = type_name(t["name"])
            # Map-only models are Go maps: no pointer needed.
            if target and target["kind"] == "object" and not target["fields"]:
                return base
        else:
            return "any"
        return "*" + base if optional else base

    def emit_object(self, m):
        name = type_name(m["name"])
        out = comment(m["description"]) or f"// {name} is a Truto API object.\n"
        if not m["fields"]:
            values = self.go_type(m["additional"], False) if m.get("additional") else "any"
            out = out.rstrip("\n") + "\n" if out else ""
            return out + f"type {name} map[string]{values}\n"
        out += f"type {name} struct {{\n"
        used = set()
        for f in m["fields"]:
            fname = go_field_name(f["name"])
            while fname in used:
                fname += "_"
            used.add(fname)
            optional = (not f["required"]) or f["nullable"]
            gtype = self.go_type(f["type"], optional)
            tag = f["name"] + (",omitempty" if optional else "")
            out += comment(f["description"], "\t")
            out += f"\t{fname} {gtype} `json:\"{tag}\"`\n"
        out += "}\n"
        return out

    def emit_enum(self, m):
        name = type_name(m["name"])
        out = comment(m["description"]) or f"// {name} is an open string enum.\n"
        out = out.rstrip("\n") + "\n//\n// Values the SDK does not know about decode unchanged.\n"
        out += f"type {name} string\n\n"
        out += f"// Known {name} values.\nconst (\n"
        used = set()
        for v in m["values"]:
            cname = name + (go_field_name(v) if v else "Empty")
            while cname in used:
                cname += "_"
            used.add(cname)
            out += f"\t{cname} {name} = {json.dumps(v)}\n"
        out += ")\n"
        return out

    def emit_union(self, m):
        name = type_name(m["name"])
        out = comment(m["description"]) or f"// {name} is one of several object shapes.\n"
        out = out.rstrip("\n") + "\n//\n// It holds the raw JSON; use the As* methods to decode a specific shape.\n"
        out += f"type {name} struct {{\n\traw json.RawMessage\n}}\n\n"
        out += f"// Raw returns the JSON this {name} holds.\n"
        out += f"func (u {name}) Raw() json.RawMessage {{ return u.raw }}\n\n"
        out += f"// MarshalJSON implements json.Marshaler.\n"
        out += f"func (u {name}) MarshalJSON() ([]byte, error) {{\n\tif len(u.raw) == 0 {{\n\t\treturn []byte(\"null\"), nil\n\t}}\n\treturn u.raw, nil\n}}\n\n"
        out += f"// UnmarshalJSON implements json.Unmarshaler.\n"
        out += f"func (u *{name}) UnmarshalJSON(b []byte) error {{\n\tu.raw = append(json.RawMessage(nil), b...)\n\treturn nil\n}}\n\n"
        if m.get("discriminator"):
            d = m["discriminator"]
            out += f"// Discriminator returns the value of the {json.dumps(d)} property, which tells the shapes apart.\n"
            out += f"func (u {name}) Discriminator() string {{\n\tvar probe struct {{\n\t\tV string `json:{json.dumps(d)}`\n\t}}\n\t_ = json.Unmarshal(u.raw, &probe)\n\treturn probe.V\n}}\n\n"
        for v in m["variants"]:
            vt = type_name(v)
            out += f"// As{vt} decodes the value as a {vt}.\n"
            out += f"func (u {name}) As{vt}() (*{vt}, error) {{\n\tvar v {vt}\n\tif err := json.Unmarshal(u.raw, &v); err != nil {{\n\t\treturn nil, err\n\t}}\n\treturn &v, nil\n}}\n\n"
            out += f"// {name}From{vt} wraps a {vt}.\n"
            out += f"func {name}From{vt}(v {vt}) ({name}, error) {{\n\tb, err := json.Marshal(v)\n\treturn {name}{{raw: b}}, err\n}}\n\n"
        return out.rstrip("\n") + "\n"

    def run(self):
        parts = [
            "// Code generated by codegen/generate_go.py from spec/openapi.json. DO NOT EDIT.\n\n"
            "package truto\n\nimport \"encoding/json\"\n\n"
            "var _ = json.Marshal\n"
        ]
        for m in self.ir["models"]:
            if m["kind"] == "object":
                parts.append(self.emit_object(m))
            elif m["kind"] == "enum":
                parts.append(self.emit_enum(m))
            elif m["kind"] == "union":
                parts.append(self.emit_union(m))
        return "\n".join(parts)


def main():
    ir = json.load(open(sys.argv[1]))
    sys.stdout.write(Gen(ir).run())


if __name__ == "__main__":
    main()
