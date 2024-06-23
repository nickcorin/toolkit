package sqlkit

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
)

var scangenTemplate = template.Must(template.New("repository").Funcs(template.FuncMap{
	"cleanPath": func(pkg string) string {
		return filepath.Base(pkg)
	},
	"cleanPkg": func(pkg string) string {
		if strings.Index(pkg, ".") > 0 {
			return strings.Split(pkg, ".")[1]
		}
		return pkg
	},
	"cols": func(fields []*field) []string {
		cols := make([]string, 0)
		for _, f := range fields {
			cols = append(cols, fmt.Sprintf("\"%s\"", f.Col()))
		}
		return cols
	},
	"export": func(s string) string {
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"fields": func(prefix string, fields []*field) []string {
		fs := make([]string, 0)
		for _, f := range fields {
			fs = append(fs, fmt.Sprintf("&%s.%s", prefix, f.Var))
		}
		return fs
	},
	"join": func(s []string) string {
		return strings.Join(s, ", ")
	},
	"unexport": func(s string) string {
		return strings.ToLower(s[:1]) + s[1:]
	},
}).Parse(repoTpl))

const repoTpl = `// Code generated by scangen. DO NOT EDIT.
package {{ .File.Pkg }}

{{ if .File.Imports -}}
import (
{{- range $path := .File.Imports }}
    {{ $path }}
{{- end }}
)
{{ end }}
// Err{{ .File.SourceType | cleanPkg }}NotFound is returned when a query for a {{ .File.SourceType }} returns no results.
var Err{{ .File.SourceType | cleanPkg }}NotFound = errors.New("{{ .File.SourceType | cleanPkg | unexport }} not found")

type {{ .Dialect.String | export }}Repository struct {
    conn *sql.DB
    tableName string
    cols     []string
}

func New{{ .Dialect.String | export }}Repository(conn *sql.DB) *{{ .Dialect.String | export }}Repository {
    return &{{ .Dialect.String | export }}Repository{
        conn: conn,
        tableName: "{{ .TableName }}",
        cols: []string{ {{- .File.Fields | cols | join -}} },
    }
}

func (r *{{ .Dialect.String | export }}Repository) selectPrefix() string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.cols, ", "), r.tableName)
}

func (r *{{ .Dialect.String | export }}Repository) lookupWhere(ctx context.Context, where string, args ...any) (*{{ .File.SourceType }}, error) {
    row := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix() + " WHERE %s", where), args...)
    return r.scan(row)
}

func (r *{{ .Dialect.String | export }}Repository) listWhere(ctx context.Context, where string, args ...any) ([]*{{ .File.SourceType }}, error) {
    rows := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix() + " WHERE %s", where), args...)
    return r.list(rows)
}

func (r *{{ .Dialect.String | export }}Repository) list(rows *sql.Rows) ([]*{{ .File.SourceType }}, error) {
   ret := make([]*{{ .File.SourceType }}, 0)
    for rows.Next() {
        item, err := r.scan(rows)
        if err != nil {
            return nil, err
        }

       ret = append(ret, item)
    }

    return ret, nil
}

func (r *{{ .Dialect.String | export }}Repository) scan(row sqlkit.Scannable) (*{{ .File.SourceType }}, error) {
    var scan {{ .File.ScanType }}

    err := row.Scan({{ fields "scan" .File.Fields | join }})
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, Err{{ .File.SourceType | cleanPkg }}NotFound
        }

        return nil, fmt.Errorf("scan {{ .File.SourceType | cleanPkg | unexport }}: %w", err)
    }

    var ret {{ .File.SourceType }}
    {{ range $f := .File.Fields }}
	{{- if eq $f.Conversion nil }}
	ret.{{ $f.Var }} = scan.{{ $f.Var }}
	{{- else if $f.Conversion.Field }}
	ret.{{ $f.Var }} = scan.{{ $f.Var }}.{{ $f.Conversion.Field }}
	{{- else if $f.Conversion.Method }}
	ret.{{ $f.Var }} = scan.{{ $f.Var }}.{{ $f.Conversion.Method }}()
	{{- else if $f.Conversion.Type }}
	ret.{{ $f.Var }} = {{ $f.Conversion.Type.String | cleanPath }}(scan.{{ $f.Var }})
	{{- end }}
    {{- end }}

    return &ret, nil
}
`
