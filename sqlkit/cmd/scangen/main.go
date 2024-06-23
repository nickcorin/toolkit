package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/nickcorin/toolkit/sqlkit"
)

/*
 *
 * go:generate --in=. --out=. --dialect=postgres --type=Foo --local=github.com/nickcorin/toolkit/sqlkit/testdata/foo
 *
 */

var (
	inFile  = flag.String("in", os.Getenv("GOFILE"), "file containing the scangen type")
	outFile = flag.String("out", "", "file to ")

	tableName    = flag.String("tableName", "", "name of the table")
	targetStruct = flag.String("type", "", "target struct")
	dialect      = flag.String("dialect", "", "sql dialect")
	locals       = flag.String("local", "", "comma separated list of local package paths")
)

func main() {
	if *inFile == "" {
		errorOut(1, "input file is required")
	}

	if *outFile == "" {
		*outFile = fmt.Sprintf("%s_gen.go", *inFile)
	}

	if *targetStruct == "" {
		errorOut(1, "target struct is required")
	}

	if *dialect == "" {
		errorOut(1, "dialect is required")
	}

	p := sqlkit.NewParser()

	parsed, err := p.Parse(*inFile, *targetStruct)
	if err != nil {
		errorOut(1, "could not parse file: %w", err)
	}

	config := sqlkit.GenerateConfig{
		Dialect:    sqlkit.GetDialectFromString(*dialect),
		TableName:  *tableName,
		OutputFile: *outFile,
		LocalPaths: *locals,
	}

	err = sqlkit.Generate(&config, parsed)
	if err != nil {
		errorOut(1, "could not generate code: %w", err)
	}
}

func errorOut(exitCode int, msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(exitCode)
}
