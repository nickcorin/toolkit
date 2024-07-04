package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickcorin/toolkit/sqlkit"
)

/*
 *
 * go:generate --in=. --out=. --dialect=postgres --type=Foo --local=github.com/nickcorin/toolkit/sqlkit/testdata/foo
 *
 */

var (
	inFile  = flag.String("in", os.Getenv("GOFILE"), "file containing the scangen type")
	outFile = flag.String("out", "", "output file to write")
	table   = flag.String("table", "", "name of the table")
	inType  = flag.String("type", "", "input struct")
	outType = flag.String("outType", "", "output struct")
	dialect = flag.String("dialect", "", "sql dialect")
	locals  = flag.String("local", "", "comma separated list of local package paths")
)

func main() {
	flag.Parse()

	if *inFile == "" {
		errorOut(1, "input file is required")
	}

	if filepath.Ext(*inFile) != ".go" {
		errorOut(1, "input file must be a .go file")
	}

	if *outFile == "" {
		*outFile = fmt.Sprintf("%s_gen.go", strings.TrimSuffix(*inFile, ".go"))
	}

	if *inType == "" {
		errorOut(1, "target type cannot be empty")
	}

	if *dialect == "" {
		errorOut(1, "dialect is required")
	}

	p := sqlkit.NewParser()

	parsed, err := p.Parse(*inFile, *inType)
	if err != nil {
		errorOut(1, "could not parse file: %w", err)
	}

	config := sqlkit.GenerateConfig{
		Dialect:      sqlkit.GetDialectFromString(*dialect),
		TableName:    *table,
		OutputFile:   *outFile,
		OutputStruct: *outType,
		LocalPaths:   *locals,
	}

	err = sqlkit.Generate(&config, parsed)
	if err != nil {
		errorOut(1, "could not generate code: %w", err)
	}
}

func errorOut(exitCode int, msg string, args ...any) {
	slog.Error(msg, args...)
	flag.PrintDefaults()
	os.Exit(exitCode)
}
