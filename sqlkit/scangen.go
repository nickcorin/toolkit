package sqlkit

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/nickcorin/toolkit/sqlkit/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/imports"
)

// GenerateConfig contains configuration options for generation.
type GenerateConfig struct {
	Dialect   Dialect
	TableName string

	// Path to which to write the generated code.
	OutputFile string

	// Name of the struct to generate.
	OutputStruct string

	// A comma separated list of import paths to pass to goimports.
	// See: https://pkg.go.dev/golang.org/x/tools/imports/forward.go#LocalPrefix
	LocalPaths string
}

func Generate(config *GenerateConfig, pf *parsedFile) error {
	if !config.Dialect.Valid() {
		return fmt.Errorf("invalid dialect: '%s'", config.Dialect)
	}

	if config.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	if config.OutputFile == "" {
		return fmt.Errorf("output file cannot be empty")
	}

	td := templateData{
		Dialect:      config.Dialect,
		File:         pf,
		OutputStruct: config.OutputStruct,
		TableName:    config.TableName,
	}

	if td.OutputStruct == "" {
		td.OutputStruct = cases.Title(language.English).String(config.Dialect.String()) + "Repository"
	}

	var buffer bytes.Buffer
	t, err := template.New("scangen").Funcs(template.FuncMap{
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
	}).Parse(templates.Scangen)
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	if err := t.Execute(&buffer, td); err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	src, err := format.Source(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("could not format generated code: %w", err)
	}

	imports.LocalPrefix = config.LocalPaths
	src, err = imports.Process(config.OutputFile, src, nil)
	if err != nil {
		return fmt.Errorf("could not format imports: %w", err)
	}

	return os.WriteFile(config.OutputFile, src, 0o755)
}

type templateData struct {
	Dialect      Dialect
	OutputStruct string
	TableName    string
	File         *parsedFile
}

type parsedFile struct {
	Pkg        string
	Imports    []string
	ScanType   string
	SourceType string
	Fields     []*field
}

type field struct {
	Var        string
	Type       types.Type
	Conversion *conversion
	Tag        string
}

func (f *field) Col() string {
	if f.Tag != "" {
		return f.Tag
	}

	return toSnakeCase(f.Var)
}

type conversion struct {
	Field  string
	Method string
	Type   types.Type
}

type Parser struct {
	builtIns   map[string]bool
	wellKnowns map[string]*conversion
}

func NewParser() *Parser {
	p := Parser{
		builtIns:   make(map[string]bool),
		wellKnowns: make(map[string]*conversion),
	}

	p.builtIns = map[string]bool{
		"time.Time": true,
	}

	p.wellKnowns = map[string]*conversion{
		"sql.NullBool":    {Field: "Bool"},
		"sql.NullBytes":   {Field: "Bytes"},
		"sql.NullInt16":   {Field: "Int16"},
		"sql.NullInt32":   {Field: "Int32"},
		"sql.NullInt64":   {Field: "Int64"},
		"sql.NullString":  {Field: "String"},
		"sql.NullTime":    {Field: "Time"},
		"sql.Null[T any]": {Field: "V"},
	}

	return &p
}

func (p *Parser) isBuiltIn(s string) bool {
	return p.builtIns[s]
}

func (p *Parser) isWellKnown(s string) bool {
	_, ok := p.wellKnowns[s]
	return ok
}

func (p *Parser) Parse(sourceFile, scangenType string) (*parsedFile, error) {
	if filepath.Ext(sourceFile) != ".go" {
		return nil, fmt.Errorf("source file '%s' is not a Go file", sourceFile)
	}

	var file parsedFile
	fset := token.NewFileSet()
	asts := make([]*ast.File, 0)
	dir := filepath.Dir(sourceFile)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return fmt.Errorf("could not parse source file '%s': %w", path, err)
		}

		// Skip generated files.
		//
		// This reduces the search space later, and also prevents some errors when multiple generated test files
		// exist within the same package.
		if ast.IsGenerated(f) {
			return nil
		}

		if strings.Contains(f.Name.Name, "_test") {
			return nil
		}

		file.Pkg = f.Name.Name
		asts = append(asts, f)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk dir '%s': %w", dir, err)
	}

	conf := types.Config{Importer: importer.ForCompiler(fset, "source", nil)}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}

	pkg, err := conf.Check(file.Pkg, fset, asts, &info)
	if err != nil {
		return nil, fmt.Errorf("could not type check source file: %w", err)
	}

	o := pkg.Scope().Lookup(scangenType)
	if o == nil {
		return nil, fmt.Errorf("scangen struct '%s' not found in source file", scangenType)
	}

	scangenStruct, ok := o.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("target struct '%s' is not a struct", scangenType)
	}
	file.ScanType = o.Name()

	var baseStruct *types.Struct
	// Find the embedded source type, and ensure there is only one.
	for i := 0; i < scangenStruct.NumFields(); i++ {
		if !scangenStruct.Field(i).Embedded() {
			continue
		}
		if baseStruct != nil {
			return nil, fmt.Errorf("multiple embedded structs found in target struct")
		}

		named, ok := scangenStruct.Field(i).Type().(*types.Named)
		if !ok {
			return nil, fmt.Errorf("embedded struct is not a named type")
		}

		baseStruct = scangenStruct.Field(i).Type().Underlying().(*types.Struct)

		if named.Obj().Pkg().Name() != file.Pkg {
			file.SourceType = filepath.Base(scangenStruct.Field(i).Type().String())
		} else {
			file.SourceType = named.Obj().Name()
		}
	}

	if file.SourceType == "" {
		return nil, fmt.Errorf("no embedded struct found in target struct")
	}

	baseFields, err := p.processStruct(baseStruct)
	if err != nil {
		return nil, err
	}

	overrideFields, err := p.processStruct(scangenStruct)
	if err != nil {
		return nil, err
	}

	fs, err := p.processOverrides(baseFields, overrideFields)
	if err != nil {
		return nil, err
	}
	file.Fields = fs

	imports := make([]string, 0)
	for _, i := range pkg.Imports() {
		imports = append(imports, fmt.Sprintf("%s \"%s\"", i.Name(), i.Path()))
	}
	file.Imports = imports

	return &file, nil
}

func (p *Parser) processOverrides(fields, overrides []*field) ([]*field, error) {
	fm := make(map[string]*field)

	// Index the overrides by field name.
	for _, f := range overrides {
		fm[f.Var] = f
	}

	for i := 0; i < len(fields); i++ {
		override, ok := fm[fields[i].Var]
		if !ok {
			continue
		}

		// Remove the field if the tag is set to "-".
		if override.Tag == "-" {
			fields = append(fields[:i], fields[i+1:]...)
			i--

			continue
		}

		// Directly override the original field if the types are the same.
		if fields[i].Type == override.Type {
			fields[i] = override
			continue
		}

		// Add the override if the types are convertable.
		if fields[i].Type.Underlying() == override.Type.Underlying() {
			override.Conversion = &conversion{Type: fields[i].Type}
			fields[i] = override
			continue
		}

		// There is some magic happening here to check generic types.
		// For example sql.NullTime vs sql.Null[time.Time] (sql.Null[T any]).
		if s, ok := override.Type.(*types.Named); ok {
			t := filepath.Base(s.Origin().Obj().Type().String())
			c, ok := p.wellKnowns[t]
			if !ok {
				return nil, fmt.Errorf("type '%s' (field '%s') is not well-known", override.Type, override.Var)
			}

			if s.TypeArgs().Len() > 0 && s.TypeArgs().At(0).String() != fields[i].Type.String() {
				return nil, fmt.Errorf("cannot convert type '%s' to '%s' (field '%s')", override.Type, fields[i].Type, override.Var)
			}

			override.Conversion = c
			fields[i] = override
		}
	}

	return fields, nil
}

func (p *Parser) processStruct(s *types.Struct) ([]*field, error) {
	var fs []*field
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		// We skip embedded fields.
		if f.Embedded() {
			continue
		}

		switch t := f.Type().(type) {
		case *types.Basic:
			break
		case *types.Named:
			if _, ok := t.Underlying().(*types.Basic); ok {
				break
			}

			name := filepath.Base(t.Origin().String())

			if !p.isBuiltIn(name) && !p.isWellKnown(name) {
				return nil, fmt.Errorf("unsupported field type '%s' (field '%s')", f.Type(), f.Name())
			}

		default:
			return nil, fmt.Errorf("unsupported field type '%s' (field '%s')", f.Type(), f.Name())
		}

		fs = append(fs, &field{
			Var:  f.Name(),
			Type: f.Type(),
			Tag:  reflect.StructTag(s.Tag(i)).Get("sqlkit"),
		})
	}

	return fs, nil
}

func toSnakeCase(s string) string {
	var result string
	var prevCharIsUpper bool

	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 && !prevCharIsUpper {
				result += "_"
			}
			result += string(c + 32)
			prevCharIsUpper = true
		} else {
			result += string(c)
			prevCharIsUpper = false
		}
	}

	return result
}
