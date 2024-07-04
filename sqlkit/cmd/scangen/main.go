package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickcorin/toolkit/sqlkit"
)

var (
	inFile  = flag.String("inFile", os.Getenv("GOFILE"), "file containing the scangen type")
	outFile = flag.String("outFile", "", "output file to write")

	table   = flag.String("table", "", "name of the table")
	inType  = flag.String("inType", "", "input struct")
	outType = flag.String("outType", "", "output struct")

	dialect = flag.String("dialect", "", "sql dialect")
	locals  = flag.String("local", "", "comma separated list of local package paths")
)

func main() {
	flag.Parse()

	if *inFile == "" {
		errorOut(1, errors.New("input file is required"))
	}

	if filepath.Ext(*inFile) != ".go" {
		errorOut(1, errors.New("input file must be a .go file"))
	}

	if *outFile == "" {
		*outFile = fmt.Sprintf("%s_gen.go", strings.TrimSuffix(*inFile, ".go"))
	}

	if *inType == "" {
		errorOut(1, errors.New("target type cannot be empty"))
	}

	if *dialect == "" {
		errorOut(1, errors.New("dialect is required"))
	}

	p := sqlkit.NewParser()

	parsed, err := p.Parse(*inFile, *inType)
	if err != nil {
		errorOut(1, fmt.Errorf("could not parse file: %w", err))
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
		errorOut(1, fmt.Errorf("could not generate code: %w", err))
	}
}

func errorOut(exitCode int, err error) {
	slog.Error(err.Error())
	flag.PrintDefaults()
	os.Exit(exitCode)
}
